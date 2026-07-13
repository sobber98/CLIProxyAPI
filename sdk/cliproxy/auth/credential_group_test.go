package auth

import (
	"testing"

	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/executor"
)

func TestAuthMatchesCredentialGroup(t *testing.T) {
	strictTeamA := map[string]any{
		cliproxyexecutor.CredentialGroupStrictMetadataKey: true,
		cliproxyexecutor.CredentialGroupMetadataKey:       "team-a",
	}
	if !authMatchesCredentialGroup(&Auth{Group: " team-a "}, strictTeamA) {
		t.Fatal("matching credential group was rejected")
	}
	if authMatchesCredentialGroup(&Auth{Group: "team-b"}, strictTeamA) {
		t.Fatal("different credential group was accepted")
	}
	if authMatchesCredentialGroup(&Auth{Group: "team-a"}, map[string]any{
		cliproxyexecutor.CredentialGroupStrictMetadataKey: true,
		cliproxyexecutor.CredentialGroupMetadataKey:       "",
	}) {
		t.Fatal("grouped credential was accepted for ungrouped downstream key")
	}
	if !authMatchesCredentialGroup(&Auth{Group: "team-b"}, nil) {
		t.Fatal("group filtering changed compatibility mode")
	}
}
