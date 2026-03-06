package app

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/example/deckhand/internal/resources"
	"github.com/example/deckhand/internal/ui/modes"
	"github.com/example/deckhand/internal/ui/panes"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && m.pendingConfirm != nil {
		return m.handleConfirmationInput(keyMsg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 1
		footerHeight := 1
		m.list.SetSize(msg.Width/3, msg.Height-headerHeight-footerHeight-2)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			if m.listScope == listScopeNamespaces || m.listScope == listScopeContexts {
				m.listScope = m.previousListScope
				if m.mode == modes.PipelineMode {
					m.list.Title = "PipelineRuns"
				} else {
					m.list.Title = "Pods"
				}
				return m, m.refreshCurrentScopeCmd()
			}
			return m, tea.Quit
		case key.Matches(msg, m.keys.OpsMode):
			m.mode = modes.OpsMode
			m.listScope = listScopePods
			m.list.Title = "Pods"
			m.detailText = "Switching to Ops mode..."
			return m, tea.Batch(m.refreshPodsCmd(), m.refreshPermissionsCmd())
		case key.Matches(msg, m.keys.PipelineMode):
			m.mode = modes.PipelineMode
			m.listScope = listScopePipelines
			m.list.Title = "PipelineRuns"
			m.detailText = "Loading PipelineRuns..."
			return m, tea.Batch(m.refreshPipelineRunsCmd(), m.refreshPermissionsCmd())
		case key.Matches(msg, m.keys.Namespace):
			m.previousListScope = m.listScope
			m.listScope = listScopeNamespaces
			m.list.Title = "Namespaces/Projects"
			m.detailText = "Select namespace/project and press enter"
			return m, m.loadNamespacesCmd()
		case key.Matches(msg, m.keys.Context):
			m.previousListScope = m.listScope
			m.listScope = listScopeContexts
			m.list.Title = "Contexts"
			m.detailText = "Select context and press enter"
			return m, m.loadContextsCmd()
		case key.Matches(msg, m.keys.Routes):
			m.listScope = listScopeRoutes
			m.list.Title = "Routes"
			return m, m.refreshRoutesCmd()
		case key.Matches(msg, m.keys.Refresh):
			return m, tea.Batch(m.refreshCurrentScopeCmd(), m.refreshPermissionsCmd())
		case key.Matches(msg, m.keys.TaskRuns):
			if m.mode == modes.PipelineMode {
				m.listScope = listScopeTaskRuns
				m.list.Title = "TaskRuns"
				return m, m.loadTaskRunsCmd()
			}
			return m, nil
		case key.Matches(msg, m.keys.OpenPod):
			if m.mode == modes.PipelineMode && m.listScope == listScopeTaskRuns {
				m.activeTab = "logs"
				return m, m.loadLogsCmd()
			}
			return m, nil
		case key.Matches(msg, m.keys.Details):
			if m.listScope == listScopeNamespaces {
				if selected, ok := m.list.SelectedItem().(panes.NamespaceItem); ok {
					m.currentNamespace = selected.Name
					m.cfg.Namespace = selected.Name
					m.listScope = m.previousListScope
					if m.mode == modes.PipelineMode {
						m.list.Title = "PipelineRuns"
					} else {
						m.list.Title = "Pods"
					}
					m.statusText = "Switched to namespace/project: " + selected.Name
					return m, tea.Batch(m.refreshCurrentScopeCmd(), m.refreshPermissionsCmd())
				}
				return m, nil
			}
			if m.listScope == listScopeContexts {
				if selected, ok := m.list.SelectedItem().(panes.ContextItem); ok {
					m.statusText = "Switching context: " + selected.Name
					return m, m.switchContextCmd(selected.Name)
				}
				return m, nil
			}
			m.activeTab = "details"
			if m.mode == modes.OpsMode || m.listScope == listScopeTaskRuns {
				return m, m.loadDetailsCmd()
			}
			return m, nil
		case key.Matches(msg, m.keys.Events):
			m.activeTab = "events"
			if m.mode == modes.OpsMode || m.listScope == listScopeTaskRuns {
				return m, m.loadEventsCmd()
			}
			return m, nil
		case key.Matches(msg, m.keys.YAML):
			m.activeTab = "yaml"
			if m.mode == modes.OpsMode || m.listScope == listScopeTaskRuns {
				return m, m.loadYAMLCmd()
			}
			return m, nil
		case key.Matches(msg, m.keys.Logs):
			m.activeTab = "logs"
			if m.mode == modes.OpsMode || m.listScope == listScopeTaskRuns {
				return m, m.loadLogsCmd()
			}
			return m, nil
		case key.Matches(msg, m.keys.Delete):
			if m.cfg.ReadOnly {
				m.statusText = "Read-only mode: delete blocked"
				return m, nil
			}
			if !m.permissions.DeleteAllowed {
				m.statusText = "Delete not allowed"
				m.detailText = "RBAC: " + m.permissions.DeleteReason
				return m, nil
			}
			if pod := m.selectedPodName(); pod != "" {
				m.pendingConfirm = &pendingConfirmation{Action: "delete", Target: pod, Message: "Delete pod " + pod + "?"}
			}
			return m, nil
		case key.Matches(msg, m.keys.Patch):
			if m.cfg.ReadOnly {
				m.statusText = "Read-only mode: patch blocked"
				return m, nil
			}
			if !m.permissions.PatchAllowed {
				m.statusText = "Patch not allowed"
				m.detailText = "RBAC: " + m.permissions.PatchReason
				return m, nil
			}
			if pod := m.selectedPodName(); pod != "" {
				m.pendingConfirm = &pendingConfirmation{Action: "patch", Target: pod, Message: "Patch pod " + pod + " annotation?"}
			}
			return m, nil
		case key.Matches(msg, m.keys.Cancel):
			if m.cfg.ReadOnly {
				m.statusText = "Read-only mode: cancel blocked"
				return m, nil
			}
			if !m.permissions.CancelAllowed {
				m.statusText = "Cancel not allowed"
				m.detailText = "RBAC: " + m.permissions.CancelReason
				return m, nil
			}
			if pr := m.selectedPipelineRunName(); pr != "" {
				m.pendingConfirm = &pendingConfirmation{Action: "cancel", Target: pr, Message: "Cancel PipelineRun " + pr + "?"}
			}
			return m, nil
		}

	case podsLoadedMsg:
		if msg.Err != nil {
			m.errText = msg.Err.Error()
			m.statusText = "Failed to load pods"
			return m, nil
		}
		m.list.SetItems(panes.PodItemsFromSummaries(msg.Pods))
		m.statusText = fmt.Sprintf("Loaded %d pods", len(msg.Pods))
		if len(msg.Pods) > 0 {
			m.listScope = listScopePods
			return m, m.loadDetailsCmd()
		}
		m.detailText = "No pods found."
		return m, nil

	case routesLoadedMsg:
		if msg.Err != nil {
			if errors.Is(msg.Err, resources.ErrAPINotAvailable) {
				m.errText = ""
				m.statusText = "OpenShift Route API not available"
				m.detailText = "This cluster does not expose route.openshift.io/v1 routes."
				m.list.SetItems([]list.Item{})
				return m, nil
			}
			m.errText = msg.Err.Error()
			m.statusText = "Failed to load routes"
			return m, nil
		}
		m.listScope = listScopeRoutes
		m.list.SetItems(panes.RouteItemsFromSummaries(msg.Routes))
		m.statusText = fmt.Sprintf("Loaded %d routes", len(msg.Routes))
		if len(msg.Routes) == 0 {
			m.detailText = "No routes found in this namespace/project."
		}
		return m, nil

	case pipelinesLoadedMsg:
		if msg.Err != nil {
			if errors.Is(msg.Err, resources.ErrAPINotAvailable) {
				m.errText = ""
				m.statusText = "Tekton API not available"
				m.detailText = "This cluster does not expose tekton.dev/v1 resources."
				return m, nil
			}
			m.errText = msg.Err.Error()
			m.statusText = "Failed to load PipelineRuns"
			return m, nil
		}
		m.listScope = listScopePipelines
		m.list.SetItems(panes.PipelineRunItemsFromSummaries(msg.PipelineRuns))
		m.statusText = fmt.Sprintf("Loaded %d PipelineRuns", len(msg.PipelineRuns))
		if len(msg.PipelineRuns) > 0 {
			m.selectedPipeline = msg.PipelineRuns[0].Name
		}
		return m, nil

	case taskRunsLoadedMsg:
		if msg.Err != nil {
			m.errText = msg.Err.Error()
			m.statusText = "Failed to load TaskRuns"
			return m, nil
		}
		m.listScope = listScopeTaskRuns
		m.list.SetItems(panes.TaskRunItemsFromSummaries(msg.TaskRuns))
		m.statusText = fmt.Sprintf("Loaded %d TaskRuns", len(msg.TaskRuns))
		return m, nil

	case namespacesLoadedMsg:
		if msg.Err != nil {
			m.errText = msg.Err.Error()
			m.statusText = "Failed to load namespaces/projects"
			return m, nil
		}
		m.listScope = listScopeNamespaces
		m.list.SetItems(panes.NamespaceItems(msg.Namespaces, msg.IsProject))
		if msg.IsProject {
			m.statusText = fmt.Sprintf("Loaded %d OpenShift projects", len(msg.Namespaces))
		} else {
			m.statusText = fmt.Sprintf("Loaded %d namespaces", len(msg.Namespaces))
		}
		return m, nil

	case contextsLoadedMsg:
		if msg.Err != nil {
			m.errText = msg.Err.Error()
			m.statusText = "Failed to load contexts"
			return m, nil
		}
		m.listScope = listScopeContexts
		m.list.SetItems(panes.ContextItems(msg.Contexts, m.cfg.Context))
		m.statusText = fmt.Sprintf("Loaded %d contexts", len(msg.Contexts))
		return m, nil

	case contextSwitchedMsg:
		if msg.Err != nil {
			m.errText = msg.Err.Error()
			m.statusText = "Failed to switch context"
			return m, nil
		}
		m = m.withKube(msg.Kube)
		m.cfg.Context = msg.Context
		m.selectedContext = msg.Context
		m.currentNamespace = msg.Namespace
		m.cfg.Namespace = msg.Namespace
		m.listScope = m.previousListScope
		if m.mode == modes.PipelineMode {
			m.list.Title = "PipelineRuns"
		} else {
			m.list.Title = "Pods"
		}
		m.statusText = fmt.Sprintf("Switched context to %s (ns=%s)", msg.Context, msg.Namespace)
		return m, tea.Batch(m.refreshCurrentScopeCmd(), m.refreshPermissionsCmd())

	case textLoadedMsg:
		if msg.Err != nil {
			m.errText = msg.Err.Error()
			m.detailText = "Error: " + msg.Err.Error()
			m.statusText = "Action failed"
			return m, nil
		}
		m.errText = ""
		m.activeTab = msg.Tab
		m.detailText = msg.Text
		if msg.Tab != "" {
			m.statusText = strings.ToUpper(msg.Tab[:1]) + msg.Tab[1:] + " loaded"
		}
		return m, nil

	case permissionsLoadedMsg:
		if msg.Err != nil {
			m.errText = msg.Err.Error()
			m.statusText = "Failed to evaluate permissions"
			return m, nil
		}
		m.permissions.DeleteAllowed = msg.DeleteAllowed
		m.permissions.DeleteReason = msg.DeleteReason
		m.permissions.PatchAllowed = msg.PatchAllowed
		m.permissions.PatchReason = msg.PatchReason
		m.permissions.CancelAllowed = msg.CancelAllowed
		m.permissions.CancelReason = msg.CancelReason
		return m, nil

	case actionDoneMsg:
		if msg.Err != nil {
			m.errText = msg.Err.Error()
			m.statusText = "Action failed: " + msg.Action
			m.pendingConfirm = nil
			return m, nil
		}
		m.pendingConfirm = nil
		m.statusText = "Action completed: " + msg.Action
		return m, m.refreshCurrentScopeCmd()

	case time.Time:
		return m, tea.Batch(
			m.refreshCurrentScopeCmd(),
			m.refreshPermissionsCmd(),
			tickRefresh(time.Duration(m.cfg.UI.RefreshIntervalSeconds)*time.Second),
		)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if m.listScope == listScopePipelines {
		m.selectedPipeline = m.selectedPipelineRunName()
	}
	return m, cmd
}

func (m Model) handleConfirmationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y", "enter":
		action := m.pendingConfirm.Action
		switch action {
		case "delete":
			return m, m.deletePodCmd()
		case "patch":
			return m, m.patchPodCmd()
		case "cancel":
			return m, m.cancelPipelineRunCmd()
		default:
			m.pendingConfirm = nil
			return m, nil
		}
	case "n", "N", "esc":
		m.pendingConfirm = nil
		m.statusText = "Action canceled"
		return m, nil
	default:
		return m, nil
	}
}
