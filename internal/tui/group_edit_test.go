package tui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestKeysTabStartsGroupEdit(t *testing.T) {
	m := newKeysTabModel(nil)
	m.keys = []string{"client-key"}
	m.groups = map[string]string{"client-key": "team-a"}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("g")})
	if !updated.editingGroup {
		t.Fatal("group editing was not enabled")
	}
	if got := updated.editInput.Value(); got != "team-a" {
		t.Fatalf("group input = %q, want %q", got, "team-a")
	}
}

func TestAuthTabIncludesEditableGroup(t *testing.T) {
	if got := authEditableFields[3]; got.label != "Group" || got.key != "group" {
		t.Fatalf("group field = %#v, want Group/group", got)
	}

	m := newAuthTabModel(nil)
	m.files = []map[string]any{{"name": "credential.json", "group": "team-a"}}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("4")})
	if !updated.editing || updated.editField != 3 {
		t.Fatalf("group edit state = editing:%v field:%d", updated.editing, updated.editField)
	}
	if got := updated.editInput.Value(); got != "team-a" {
		t.Fatalf("group input = %q, want %q", got, "team-a")
	}
	if detail := m.renderDetail(m.files[0]); !strings.Contains(detail, "Group") || !strings.Contains(detail, "team-a") {
		t.Fatalf("group was not rendered in auth-file details: %q", detail)
	}
}

func TestClientPutAPIKeyGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/v0/management/api-key-groups" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		var body struct {
			Groups map[string]string `json:"api-key-groups"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := body.Groups["client-key"]; got != "team-a" {
			t.Fatalf("group = %q, want %q", got, "team-a")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{baseURL: server.URL, http: server.Client()}
	if err := client.PutAPIKeyGroups(map[string]string{"client-key": "team-a"}); err != nil {
		t.Fatalf("PutAPIKeyGroups returned error: %v", err)
	}
}
