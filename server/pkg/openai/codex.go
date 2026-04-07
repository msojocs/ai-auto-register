package openai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

type CodexClient struct {
	accountId    string
	accessToken  string
	refreshToken string
	serverURL    string
	client       *http.Client
}

type CodexConfig struct {
	AccountId    string
	AccessToken  string
	RefreshToken string
	ServerURL    string
}

type CodexUsageResponse struct {
	UserId    string `json:"user_id"`
	AccountId string `json:"account_id"`
	Email     string `json:"email"`
	PlanType  string `json:"plan_type"`
	// We only care about the primary window for now.
	RateLimit struct {
		Allowed       bool `json:"allowed"`
		LimitReached  bool `json:"limit_reached"`
		PrimaryWindow struct {
			UsedPercent        float64 `json:"used_percent"`
			LimitWindowSeconds int     `json:"limit_window_seconds"`
			ResetAfterSeconds  int     `json:"reset_after_seconds"`
			ResetAt            int64   `json:"reset_at"`
		} `json:"primary_window"`
	} `json:"rate_limit"`
	CodeReviewRateLimit struct {
		Allowed       bool `json:"allowed"`
		LimitReached  bool `json:"limit_reached"`
		PrimaryWindow struct {
			UsedPercent        float64 `json:"used_percent"`
			LimitWindowSeconds int     `json:"limit_window_seconds"`
			ResetAfterSeconds  int     `json:"reset_after_seconds"`
			ResetAt            int64   `json:"reset_at"`
		} `json:"primary_window"`
	} `json:"code_review_rate_limit"`
}

type CodexAccountCheckResponse struct {
	Accounts         []CodexAccount `json:"accounts"`
	DefaultAccountId string         `json:"default_account_id"`
	AccountOrdering  []string       `json:"account_ordering"`
}

type CodexAccount struct {
	Id                string  `json:"id"`
	AccountUserId     string  `json:"account_user_id"`
	Structure         string  `json:"structure"`
	PlanType          string  `json:"plan_type"`
	Name              *string `json:"name"`
	ProfilePictureURL *string `json:"profile_picture_url"`
}

func NewCodexClient(config CodexConfig) (*CodexClient, error) {
	if config.ServerURL == "" {
		config.ServerURL = "https://chatgpt.com/backend-api"
	}
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	return &CodexClient{
		accountId:    config.AccountId,
		accessToken:  config.AccessToken,
		refreshToken: config.RefreshToken,
		serverURL:    config.ServerURL,
		client:       &http.Client{Jar: jar},
	}, nil
}

func (c *CodexClient) buildHeaders() http.Header {
	headers := http.Header{}
	headers.Set("sentry-trace", "00000000000000000000000000000000-0000000000000000")
	headers.Set("baggage", "sentry-environment=prod,sentry-release=codex%4026.325.31654,sentry-public_key=6719eaa18601933a26ac21499dcaba2f,sentry-trace_id=00000000000000000000000000000000,sentry-org_id=33249,sentry-sampled=false")
	headers.Set("authorization", "Bearer "+c.accessToken)
	headers.Set("chatgpt-account-id", c.accountId)
	headers.Set("originator", "Codex Desktop")
	headers.Set("user-agent", "Codex Desktop/26.325.31654 (win32; x64)")
	headers.Set("sec-fetch-site", "none")
	headers.Set("sec-fetch-mode", "no-cors")
	headers.Set("sec-fetch-dest", "empty")
	headers.Set("accept-encoding", "gzip, deflate, br, zstd")
	headers.Set("accept-language", "zh-CN")
	headers.Set("priority", "u=4, i")
	return headers
}

func (c *CodexClient) QueryUsage() (*CodexUsageResponse, error) {
	url := c.serverURL + "/wham/usage"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = c.buildHeaders()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var usageResp CodexUsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&usageResp); err != nil {
		return nil, err
	}
	// We use the code review rate limit as the usage metric since it's more likely to be the one with a limit.
	if !usageResp.CodeReviewRateLimit.Allowed {
		return &usageResp, nil
	}
	return &usageResp, nil
}

func (c *CodexClient) CheckAccount() (*CodexAccountCheckResponse, error) {
	url := c.serverURL + "/wham/accounts/check"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = c.buildHeaders()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var accountResp CodexAccountCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountResp); err != nil {
		return nil, err
	}
	return &accountResp, nil
}

func (c *CodexClient) RefreshToken() error {
	// TODO: Implement token refresh logic using the refresh token.
	return nil
}
