package mailprovider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTempMailReceive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/inbox/create":
			_, _ = w.Write([]byte(`{"address":"mock@tempmail.lol","token":"mocktoken123"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v2/inbox":
			_, _ = w.Write([]byte(`{"emails":[{"_id":"msg001","subject":"Your code","body":"code: 123456","date":1775451744}]}`))
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()

	provider := NewTempMailLol(map[string]string{
		"api_url": ts.URL + "/v2",
	})

	ctx := context.Background()
	email, err := provider.GetEmail(ctx)
	if err != nil {
		t.Fatalf("GetEmail: %v", err)
	}
	t.Logf("Created email: %s", email.Email)

	code, err := provider.WaitForCode(ctx, email, "code", 10)
	if err != nil {
		t.Fatalf("WaitForCode: %v", err)
	}
	t.Logf("Received code: %s", code)
	if code != "123456" {
		t.Fatalf("Expected code 123456, got %s", code)
	}
}
