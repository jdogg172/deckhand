package app

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/example/deckhand/internal/config"
	"github.com/example/deckhand/internal/resources"
	"github.com/example/deckhand/internal/ui/panes"
)

func TestHandleConfirmationInput_Cancel(t *testing.T) {
	m := Model{pendingConfirm: &pendingConfirmation{Action: "delete"}}

	updated, _ := m.handleConfirmationInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	um, ok := updated.(Model)
	if !ok {
		t.Fatalf("expected model type")
	}
	if um.pendingConfirm != nil {
		t.Fatalf("expected confirmation to be cleared")
	}
}

func TestUpdate_DeleteBlockedInReadOnly(t *testing.T) {
	keys := DefaultKeyMap()
	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)

	m := Model{
		cfg:       config.Config{ReadOnly: true},
		keys:      keys,
		list:      l,
		listScope: listScopePods,
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	um, ok := updated.(Model)
	if !ok {
		t.Fatalf("expected model type")
	}
	if um.statusText == "" {
		t.Fatalf("expected read-only status message")
	}
}

func TestUpdate_ContextPickerScope(t *testing.T) {
	keys := DefaultKeyMap()
	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)

	m := Model{
		cfg:       config.Config{},
		keys:      keys,
		list:      l,
		listScope: listScopePods,
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	um, ok := updated.(Model)
	if !ok {
		t.Fatalf("expected model type")
	}
	if um.listScope != listScopeContexts {
		t.Fatalf("expected context scope, got %q", um.listScope)
	}
	if um.list.Title != "Contexts" {
		t.Fatalf("expected list title to be Contexts")
	}
}

func TestActionStateViewIncludesDisabledReasons(t *testing.T) {
	m := Model{
		cfg: config.Config{ReadOnly: false},
		permissions: permissionState{
			DeleteAllowed: false,
			DeleteReason:  "rbac deny",
			PatchAllowed:  true,
			CancelAllowed: false,
			CancelReason:  "Tekton API not available",
		},
	}

	text := m.actionStateView()
	if !strings.Contains(text, "x=disabled(rbac deny)") {
		t.Fatalf("expected delete disabled reason in action state view, got %q", text)
	}
	if !strings.Contains(text, "k=disabled(Tekton API not available)") {
		t.Fatalf("expected cancel disabled reason in action state view, got %q", text)
	}
}

func TestUpdate_MutatingActions_ReadOnlyAndRBAC(t *testing.T) {
	delegate := list.NewDefaultDelegate()
	baseList := list.New([]list.Item{}, delegate, 0, 0)
	baseList.SetItems(panes.PodItemsFromSummaries([]resources.PodSummary{{Name: "pod-1", Namespace: "ns1"}}))

	tests := []struct {
		name        string
		keyRune     rune
		model       Model
		wantStatus  string
		wantDetail  string
		wantConfirm bool
		wantAction  string
	}{
		{
			name:    "readonly delete blocked",
			keyRune: 'x',
			model: Model{
				cfg:       config.Config{ReadOnly: true},
				keys:      DefaultKeyMap(),
				list:      baseList,
				listScope: listScopePods,
			},
			wantStatus:  "Read-only mode: delete blocked",
			wantConfirm: false,
		},
		{
			name:    "readonly patch blocked",
			keyRune: 'p',
			model: Model{
				cfg:       config.Config{ReadOnly: true},
				keys:      DefaultKeyMap(),
				list:      baseList,
				listScope: listScopePods,
			},
			wantStatus:  "Read-only mode: patch blocked",
			wantConfirm: false,
		},
		{
			name:    "readonly cancel blocked",
			keyRune: 'k',
			model: Model{
				cfg:              config.Config{ReadOnly: true},
				keys:             DefaultKeyMap(),
				list:             baseList,
				listScope:        listScopePipelines,
				selectedPipeline: "pr-1",
			},
			wantStatus:  "Read-only mode: cancel blocked",
			wantConfirm: false,
		},
		{
			name:    "rbac delete blocked",
			keyRune: 'x',
			model: Model{
				cfg:  config.Config{ReadOnly: false},
				keys: DefaultKeyMap(),
				list: baseList,
				permissions: permissionState{
					DeleteAllowed: false,
					DeleteReason:  "delete denied",
				},
				listScope: listScopePods,
			},
			wantStatus:  "Delete not allowed",
			wantDetail:  "RBAC: delete denied",
			wantConfirm: false,
		},
		{
			name:    "rbac patch blocked",
			keyRune: 'p',
			model: Model{
				cfg:  config.Config{ReadOnly: false},
				keys: DefaultKeyMap(),
				list: baseList,
				permissions: permissionState{
					PatchAllowed: false,
					PatchReason:  "patch denied",
				},
				listScope: listScopePods,
			},
			wantStatus:  "Patch not allowed",
			wantDetail:  "RBAC: patch denied",
			wantConfirm: false,
		},
		{
			name:    "rbac cancel blocked",
			keyRune: 'k',
			model: Model{
				cfg:              config.Config{ReadOnly: false},
				keys:             DefaultKeyMap(),
				list:             baseList,
				listScope:        listScopePipelines,
				selectedPipeline: "pr-1",
				permissions: permissionState{
					CancelAllowed: false,
					CancelReason:  "cancel denied",
				},
			},
			wantStatus:  "Cancel not allowed",
			wantDetail:  "RBAC: cancel denied",
			wantConfirm: false,
		},
		{
			name:    "delete confirmation when allowed",
			keyRune: 'x',
			model: Model{
				cfg:  config.Config{ReadOnly: false},
				keys: DefaultKeyMap(),
				list: baseList,
				permissions: permissionState{
					DeleteAllowed: true,
				},
				listScope: listScopePods,
			},
			wantConfirm: true,
			wantAction:  "delete",
		},
		{
			name:    "patch confirmation when allowed",
			keyRune: 'p',
			model: Model{
				cfg:  config.Config{ReadOnly: false},
				keys: DefaultKeyMap(),
				list: baseList,
				permissions: permissionState{
					PatchAllowed: true,
				},
				listScope: listScopePods,
			},
			wantConfirm: true,
			wantAction:  "patch",
		},
		{
			name:    "cancel confirmation when allowed",
			keyRune: 'k',
			model: Model{
				cfg:              config.Config{ReadOnly: false},
				keys:             DefaultKeyMap(),
				list:             baseList,
				listScope:        listScopePipelines,
				selectedPipeline: "pr-1",
				permissions: permissionState{
					CancelAllowed: true,
				},
			},
			wantConfirm: true,
			wantAction:  "cancel",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			updated, _ := tc.model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tc.keyRune}})
			um, ok := updated.(Model)
			if !ok {
				t.Fatalf("expected model type")
			}

			if tc.wantStatus != "" && um.statusText != tc.wantStatus {
				t.Fatalf("expected status %q, got %q", tc.wantStatus, um.statusText)
			}
			if tc.wantDetail != "" && um.detailText != tc.wantDetail {
				t.Fatalf("expected detail %q, got %q", tc.wantDetail, um.detailText)
			}

			if tc.wantConfirm {
				if um.pendingConfirm == nil {
					t.Fatalf("expected confirmation to be pending")
				}
				if um.pendingConfirm.Action != tc.wantAction {
					t.Fatalf("expected confirm action %q, got %q", tc.wantAction, um.pendingConfirm.Action)
				}
			} else if um.pendingConfirm != nil {
				t.Fatalf("did not expect pending confirmation, got %+v", *um.pendingConfirm)
			}
		})
	}
}
