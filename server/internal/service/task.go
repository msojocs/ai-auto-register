package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/msojocs/free2api/server/internal/core"
	"github.com/msojocs/free2api/server/internal/executor"
	"github.com/msojocs/free2api/server/internal/model"
	"github.com/msojocs/free2api/server/internal/repository"
	"github.com/msojocs/free2api/server/internal/resource"
	"gorm.io/gorm"
)

var jobCounter uint64

type TaskService struct {
	repo     repository.TaskRepository
	pool     *core.WorkerPool
	db       *gorm.DB
	proxyRes *resource.ProxyResource
}

func NewTaskService(repo repository.TaskRepository, pool *core.WorkerPool, db *gorm.DB, proxyRes *resource.ProxyResource) *TaskService {
	return &TaskService{repo: repo, pool: pool, db: db, proxyRes: proxyRes}
}

func (s *TaskService) List(page, limit int) ([]model.TaskBatch, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(offset, limit)
}

func (s *TaskService) Create(name, taskType string, total int, config map[string]interface{}) (*model.TaskBatch, error) {
	validTypes := map[string]bool{
		"chatgpt": true,
		"cursor":  true,
		"trae":    true,
		"grok":    true,
		"tavily":  true,
		"kiro":    true,
	}
	if !validTypes[taskType] {
		return nil, fmt.Errorf("invalid task type: must be one of: chatgpt, cursor, trae, grok, tavily, kiro")
	}
	task := &model.TaskBatch{
		Name:   name,
		Type:   taskType,
		Status: model.TaskStatusPending,
		Total:  total,
		Config: model.JSONMap(config),
	}
	if err := s.repo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Get(id uint) (*model.TaskBatch, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (s *TaskService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *TaskService) Start(id uint) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.Status == model.TaskStatusRunning {
		return errors.New("task is already running")
	}
	if err := s.repo.UpdateFields(id, map[string]interface{}{"status": model.TaskStatusRunning}); err != nil {
		return err
	}
	task.Status = model.TaskStatusRunning
	go s.dispatchJobs(*task)
	return nil
}

func (s *TaskService) dispatchJobs(task model.TaskBatch) {
	total := task.Total
	if total <= 0 {
		total = 1
	}

	var wg sync.WaitGroup
	for i := 0; i < total; i++ {
		current, err := s.repo.FindByID(task.ID)
		if err != nil || current == nil {
			break
		}
		if current.Status == model.TaskStatusPaused {
			return
		}

		jobID := atomic.AddUint64(&jobCounter, 1)
		taskID := task.ID
		taskType := task.Type
		cfg := map[string]interface{}(task.Config)
		cfg = s.resolveProxyConfig(cfg)
		// If the task config references a temp mail provider by ID, resolve it
		// and merge the provider's settings into the job config so executors
		// receive mail_provider + mail_* keys transparently.
		cfg = s.resolveMailProviderConfig(cfg)
		db := s.db

		wg.Add(1)
		s.pool.Submit(core.Job{
			ID: uint(jobID),
			Execute: func(ctx context.Context, publish func(core.ProgressUpdate)) {
				defer wg.Done()
				var exec executor.Executor
				switch taskType {
				case "chatgpt":
					exec = executor.NewChatGPTExecutor(db)
				case "cursor":
					exec = executor.NewCursorExecutor(db)
				default:
					return
				}
				err := exec.Execute(ctx, taskID, cfg, publish)
				if err != nil {
					s.repo.UpdateFields(taskID, map[string]interface{}{
						"failed": gorm.Expr("failed + ?", 1),
					})
				} else {
					s.repo.UpdateFields(taskID, map[string]interface{}{
						"completed": gorm.Expr("completed + ?", 1),
					})
				}
			},
		})
	}

	wg.Wait()
	finalTask, err := s.repo.FindByID(task.ID)
	if err == nil && finalTask != nil && finalTask.Status == model.TaskStatusRunning {
		s.repo.UpdateFields(task.ID, map[string]interface{}{"status": model.TaskStatusCompleted})
	}
}

func (s *TaskService) Pause(id uint) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.Status != model.TaskStatusRunning {
		return errors.New("task is not running")
	}
	return s.repo.UpdateFields(id, map[string]interface{}{"status": model.TaskStatusPaused})
}

func (s *TaskService) Retry(id uint) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.Status != model.TaskStatusFailed && task.Status != model.TaskStatusPaused {
		return fmt.Errorf("task status is %s, can only retry failed or paused tasks", task.Status)
	}
	if err := s.repo.UpdateFields(id, map[string]interface{}{
		"status": model.TaskStatusRunning,
		"failed": 0,
	}); err != nil {
		return err
	}
	task.Status = model.TaskStatusRunning
	go s.dispatchJobs(*task)
	return nil
}

func (s *TaskService) Subscribe(taskID uint) chan core.ProgressUpdate {
	return s.pool.Subscribe(taskID)
}

func (s *TaskService) Unsubscribe(taskID uint, ch chan core.ProgressUpdate) {
	s.pool.Unsubscribe(taskID, ch)
}

func (s *TaskService) GetLogs(id uint) ([]string, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	if task.Logs == "" {
		return []string{}, nil
	}
	return strings.Split(task.Logs, "\n"), nil
}

func (s *TaskService) resolveProxyConfig(cfg map[string]interface{}) map[string]interface{} {
	if strings.TrimSpace(taskConfigString(cfg, "proxy")) != "" {
		return cfg
	}
	if s.proxyRes == nil {
		return cfg
	}
	var proxy *model.Proxy
	if groupID := taskConfigUint(cfg, "proxy_group_id"); groupID > 0 {
		proxy = s.proxyRes.NextByGroupID(groupID)
	} else {
		group := strings.TrimSpace(taskConfigString(cfg, "proxy_group"))
		if group != "" {
			proxy = s.proxyRes.NextByGroupName(group)
		}
	}
	if proxy == nil {
		return cfg
	}

	merged := make(map[string]interface{}, len(cfg)+1)
	for k, v := range cfg {
		merged[k] = v
	}
	merged["proxy"] = buildProxyURL(proxy)
	return merged
}

func taskConfigString(cfg map[string]interface{}, key string) string {
	if value, ok := cfg[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func taskConfigUint(cfg map[string]interface{}, key string) uint {
	if value, ok := cfg[key]; ok {
		switch typed := value.(type) {
		case float64:
			return uint(typed)
		case int:
			return uint(typed)
		case uint:
			return typed
		}
	}
	return 0
}

func buildProxyURL(proxy *model.Proxy) string {
	protocol := strings.TrimSpace(proxy.Protocol)
	if protocol == "" {
		protocol = "http"
	}
	u := &url.URL{
		Scheme: protocol,
		Host:   net.JoinHostPort(proxy.Host, proxy.Port),
	}
	if proxy.Username != "" || proxy.Password != "" {
		u.User = url.UserPassword(proxy.Username, proxy.Password)
	}
	return u.String()
}

// resolveMailProviderConfig looks up temp_mail_provider_id in cfg, loads the
// corresponding TempMailProvider record from the database, and merges its
// settings as mail_provider / mail_* keys so executors can consume them
// without being aware of the TempMailProvider model.
func (s *TaskService) resolveMailProviderConfig(cfg map[string]interface{}) map[string]interface{} {
	raw, ok := cfg["temp_mail_provider_id"]
	if !ok {
		return cfg
	}
	var id uint
	switch v := raw.(type) {
	case float64:
		id = uint(v)
	case int:
		id = uint(v)
	case uint:
		id = v
	default:
		return cfg
	}
	if id == 0 {
		return cfg
	}

	var p model.TempMailProvider
	if err := s.db.First(&p, id).Error; err != nil {
		// Provider not found — fall through and let the executor handle the missing config.
		return cfg
	}

	// Build a merged copy so the original task config is not mutated.
	merged := make(map[string]interface{}, len(cfg)+len(p.Config)+1)
	for k, v := range cfg {
		merged[k] = v
	}
	// Provider type and per-provider config keys.
	merged["mail_provider"] = p.ProviderType
	for k, v := range p.Config {
		if _, already := merged["mail_"+k]; !already {
			merged["mail_"+k] = v
		}
	}
	return merged
}
