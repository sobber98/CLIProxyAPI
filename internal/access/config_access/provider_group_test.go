package configaccess

import (
	"net/http/httptest"
	"testing"
)

func TestProviderAuthenticateAddsCredentialGroupConstraint(t *testing.T) {
	p := newProvider("test", []string{"client-a", "client-b"}, map[string]string{"client-a": " team-a "})
	request := httptest.NewRequest("GET", "http://example.test", nil)
	request.Header.Set("Authorization", "Bearer client-a")
	result, err := p.Authenticate(request.Context(), request)
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
	if got := result.Metadata["credential_group"]; got != "team-a" {
		t.Fatalf("credential group = %q, want team-a", got)
	}
	if got := result.Metadata["credential_group_strict"]; got != "true" {
		t.Fatalf("strict = %q, want true", got)
	}

	request.Header.Set("Authorization", "Bearer client-b")
	result, err = p.Authenticate(request.Context(), request)
	if err != nil {
		t.Fatalf("Authenticate() unassigned error = %v", err)
	}
	if got := result.Metadata["credential_group"]; got != "" {
		t.Fatalf("unassigned group = %q, want empty", got)
	}
}
