package management

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v7/internal/config"
)

func TestConfiguredCredentialPatchesUpdateGroup(t *testing.T) {
	t.Setenv("MANAGEMENT_PASSWORD", "")
	h := &Handler{
		cfg: &config.Config{
			GeminiKey:           []config.GeminiKey{{APIKey: "gemini-key"}},
			InteractionsKey:     []config.GeminiKey{{APIKey: "interactions-key"}},
			ClaudeKey:           []config.ClaudeKey{{APIKey: "claude-key"}},
			CodexKey:            []config.CodexKey{{APIKey: "codex-key", BaseURL: "https://codex.example"}},
			VertexCompatAPIKey:  []config.VertexCompatKey{{APIKey: "vertex-key"}},
			OpenAICompatibility: []config.OpenAICompatibility{{Name: "compat", BaseURL: "https://compat.example"}},
		},
		configFilePath: writeTestConfigFile(t),
	}

	tests := []struct {
		name    string
		path    string
		body    string
		handler gin.HandlerFunc
		group   func() string
	}{
		{"gemini", "/gemini-api-key", `{"index":0,"value":{"group":" team-a "}}`, h.PatchGeminiKey, func() string { return h.cfg.GeminiKey[0].Group }},
		{"interactions", "/interactions-api-key", `{"index":0,"value":{"group":" team-a "}}`, h.PatchInteractionsKey, func() string { return h.cfg.InteractionsKey[0].Group }},
		{"claude", "/claude-api-key", `{"index":0,"value":{"group":" team-a "}}`, h.PatchClaudeKey, func() string { return h.cfg.ClaudeKey[0].Group }},
		{"codex", "/codex-api-key", `{"index":0,"value":{"group":" team-a "}}`, h.PatchCodexKey, func() string { return h.cfg.CodexKey[0].Group }},
		{"vertex", "/vertex-api-key", `{"index":0,"value":{"group":" team-a "}}`, h.PatchVertexCompatKey, func() string { return h.cfg.VertexCompatAPIKey[0].Group }},
		{"openai compatibility", "/openai-compatibility", `{"index":0,"value":{"group":" team-a "}}`, h.PatchOpenAICompat, func() string { return h.cfg.OpenAICompatibility[0].Group }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodPatch, tt.path, strings.NewReader(tt.body))
			ctx.Request.Header.Set("Content-Type", "application/json")
			tt.handler(ctx)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
			}
			if got := tt.group(); got != "team-a" {
				t.Fatalf("group = %q, want %q", got, "team-a")
			}
		})
	}
}

func TestPutAPIKeyGroupsRejectsUnknownKeysWithoutChangingConfiguration(t *testing.T) {
	t.Setenv("MANAGEMENT_PASSWORD", "")
	h := &Handler{
		cfg: &config.Config{SDKConfig: config.SDKConfig{
			APIKeys:      []string{"client-a"},
			APIKeyGroups: map[string]string{"client-a": "team-a"},
		}},
		configFilePath: writeTestConfigFile(t),
	}
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api-key-groups", strings.NewReader(`{"api-key-groups":{"unknown":"team-b"}}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.PutAPIKeyGroups(ctx)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if got := h.cfg.APIKeyGroups["client-a"]; got != "team-a" {
		t.Fatalf("group assignment changed to %q after rejected update", got)
	}
}

func TestDeleteAPIKeyRemovesItsGroupAssignment(t *testing.T) {
	t.Setenv("MANAGEMENT_PASSWORD", "")
	h := &Handler{
		cfg: &config.Config{SDKConfig: config.SDKConfig{
			APIKeys:      []string{"client-a", "client-b"},
			APIKeyGroups: map[string]string{"client-a": "team-a"},
		}},
		configFilePath: writeTestConfigFile(t),
	}
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/api-keys?value=client-a", nil)

	h.DeleteAPIKeys(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if _, ok := h.cfg.APIKeyGroups["client-a"]; ok {
		t.Fatal("deleted API key retained a group assignment")
	}
}
