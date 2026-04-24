package service

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/msojocs/ai-auto-register/server/internal/model"
)

func TestExtractAccessTokenExpiresAtAndAccountIDFromJWT(t *testing.T) {
	exp := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	token := buildTestJWT(t, map[string]any{
		"exp": exp.Unix(),
		"https://api.openai.com/auth": map[string]any{
			"chatgpt_account_id": "acct_123",
		},
	})

	if got := extractAccessTokenExpiresAt(token); got != exp.Format(time.RFC3339) {
		t.Fatalf("expected expires_at %q, got %q", exp.Format(time.RFC3339), got)
	}

	if got := extractChatGPTAccountIDFromAccessToken(token); got != "acct_123" {
		t.Fatalf("expected account id %q, got %q", "acct_123", got)
	}
}

func TestShouldRefreshTokenPrefersAccessTokenExp(t *testing.T) {
	token := buildTestJWT(t, map[string]any{
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	})

	extra := model.JSONMap{
		"access_token":            token,
		"access_token_expires_at": time.Now().Add(7 * 24 * time.Hour).UTC().Format(time.RFC3339),
	}

	if !shouldRefreshToken(extra) {
		t.Fatal("expected shouldRefreshToken to use access_token exp")
	}
}

func TestShouldRefreshTokenFallsBackToLegacyExpiresAt(t *testing.T) {
	extra := model.JSONMap{
		"access_token":            "invalid-token",
		"access_token_expires_at": time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339),
	}

	if !shouldRefreshToken(extra) {
		t.Fatal("expected shouldRefreshToken to fall back to legacy access_token_expires_at")
	}
}

func buildTestJWT(t *testing.T, payload map[string]any) string {
	t.Helper()

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	return header + "." + base64.RawURLEncoding.EncodeToString(payloadBytes) + ".sig"
}
