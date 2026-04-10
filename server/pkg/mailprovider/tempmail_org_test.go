package mailprovider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTempMailOrgReceive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/mailbox":
			_, _ = w.Write([]byte(`{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiZWJiYmI3N2RkMGYzNGI5MjkyZmQzMzVhOWJkMWI3NjciLCJtYWlsYm94Ijoic2FyaWRlOTM0MkBueXNwcmluZy5jb20iLCJpYXQiOjE3NzU4MzA5Nzl9.D-RTuIB5rJE3hFacEACgfBEQb_IiBz91KTPPbvgOELs",
    "mailbox": "saride9342@nyspring.com"
}`))
		case r.Method == http.MethodGet && r.URL.Path == "/messages":
			_, _ = w.Write([]byte(`{
    "mailbox": "saride9342@nyspring.com",
    "messages": [
        {
            "_id": "69d90822aa82a000e3e70adb",
            "receivedAt": 1775831075,
            "from": "Anonymousemail <noreply@anonymousemail.se>",
            "subject": "hi",
            "bodyPreview": " Powered by Anonymousemail himsg ",
            "attachmentsCount": 0
        }
    ]
}`))
		case r.Method == http.MethodGet && r.URL.Path == "/messages/69d90822aa82a000e3e70adb":
			_, _ = w.Write([]byte(`{
    "_id": "69d905b652d75f008baae825",
    "receivedAt": 1775830454,
    "user": "5c5ce8b7868d4dcdb8be94af100a2d64",
    "mailbox": "leciv98548@bpotogo.com",
    "from": "Anonymousemail <noreply@anonymousemail.se>",
    "subject": "hi",
    "bodyPreview": " Powered by Anonymousemail hi msg ",
    "bodyHtml": "<p><span style=\"color:#c0392b\">Powered by <strong>Anonymousemail</strong></span></p><p>code:123456</p>\n\n",
    "attachmentsCount": 0,
    "attachments": [],
    "createdAt": "2026-04-10T14:14:14.109Z"
}`))
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer ts.Close()

	provider := NewTempMailOrg(map[string]string{
		"api_url": ts.URL,
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
