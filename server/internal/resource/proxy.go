package resource

import (
	"sync"

	"github.com/msojocs/free2api/server/internal/model"
	"gorm.io/gorm"
)

type ProxyResource struct {
	mu      sync.RWMutex
	proxies []model.Proxy
	db      *gorm.DB
	index   int
}

func NewProxyResource(db *gorm.DB) *ProxyResource {
	r := &ProxyResource{db: db}
	r.Reload()
	return r
}

func (r *ProxyResource) Reload() {
	var proxies []model.Proxy
	r.db.Where("status = ?", "active").Find(&proxies)
	r.mu.Lock()
	r.proxies = proxies
	r.index = 0
	r.mu.Unlock()
}

func (r *ProxyResource) Next() *model.Proxy {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.proxies) == 0 {
		return nil
	}
	p := r.proxies[r.index%len(r.proxies)]
	r.index++
	return &p
}

func (r *ProxyResource) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.proxies)
}
