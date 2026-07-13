package config

import (
	"strings"
	"testing"
)

func TestNormalizeAndValidateAPIKeyGroups(t *testing.T) {
	cfg := &Config{
		SDKConfig: SDKConfig{
			APIKeys:      []string{"client-a", "client-b"},
			APIKeyGroups: map[string]string{"client-a": " team-a ", "client-b": "   "},
		},
	}
	if err := cfg.NormalizeAndValidateAPIKeyGroups(); err != nil {
		t.Fatalf("NormalizeAndValidateAPIKeyGroups() error = %v", err)
	}
	if got := cfg.APIKeyGroups["client-a"]; got != "team-a" {
		t.Fatalf("client-a group = %q, want team-a", got)
	}
	if got := cfg.APIKeyGroups["client-b"]; got != "" {
		t.Fatalf("client-b group = %q, want empty", got)
	}
}

func TestNormalizeAndValidateAPIKeyGroupsDoesNotAlterAPIKeyIdentity(t *testing.T) {
	cfg := &Config{SDKConfig: SDKConfig{
		APIKeys:      []string{" client-a "},
		APIKeyGroups: map[string]string{"client-a": "team-a"},
	}}
	if err := cfg.NormalizeAndValidateAPIKeyGroups(); err == nil {
		t.Fatal("expected differently spaced API key to be rejected")
	}
}

func TestNormalizeAndValidateAPIKeyGroupsRejectsUnknownKey(t *testing.T) {
	cfg := &Config{SDKConfig: SDKConfig{APIKeys: []string{"client-a"}, APIKeyGroups: map[string]string{"client-b": "team-b"}}}
	err := cfg.NormalizeAndValidateAPIKeyGroups()
	if err == nil || !strings.Contains(err.Error(), "client-b") {
		t.Fatalf("NormalizeAndValidateAPIKeyGroups() error = %v, want unknown client-b error", err)
	}
}
