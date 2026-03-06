package app

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jdogg172/deckhand/internal/clients"
	"github.com/jdogg172/deckhand/internal/ui/modes"
	"github.com/jdogg172/deckhand/internal/ui/panes"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (m Model) refreshCurrentScopeCmd() tea.Cmd {
	if m.mode == modes.PipelineMode {
		m.listScope = listScopePipelines
		return m.refreshPipelineRunsCmd()
	}
	m.listScope = listScopePods
	return m.refreshPodsCmd()
}

func (m Model) refreshPodsCmd() tea.Cmd {
	return func() tea.Msg {
		pods, err := m.podSvc.List(m.ctx, m.currentNamespace)
		return podsLoadedMsg{Pods: pods, Err: err}
	}
}

func (m Model) refreshRoutesCmd() tea.Cmd {
	return func() tea.Msg {
		routes, err := m.routeSvc.List(m.ctx, m.currentNamespace)
		return routesLoadedMsg{Routes: routes, Err: err}
	}
}

func (m Model) refreshPipelineRunsCmd() tea.Cmd {
	return func() tea.Msg {
		items, err := m.pipelineSvc.ListPipelineRuns(m.ctx, m.currentNamespace)
		return pipelinesLoadedMsg{PipelineRuns: items, Err: err}
	}
}

func (m Model) loadTaskRunsCmd() tea.Cmd {
	pipelineRun := m.selectedPipelineRunName()
	if pipelineRun == "" {
		return nil
	}
	return func() tea.Msg {
		taskRuns, err := m.pipelineSvc.ListTaskRunsForPipelineRun(m.ctx, m.currentNamespace, pipelineRun)
		return taskRunsLoadedMsg{TaskRuns: taskRuns, Err: err}
	}
}

func (m Model) loadNamespacesCmd() tea.Cmd {
	return func() tea.Msg {
		namespaces, isProject, err := m.namespaceSvc.List(m.ctx)
		return namespacesLoadedMsg{Namespaces: namespaces, IsProject: isProject, Err: err}
	}
}

func (m Model) loadContextsCmd() tea.Cmd {
	return func() tea.Msg {
		contexts := make([]string, 0, len(m.kube.RawConfig.Contexts))
		for name := range m.kube.RawConfig.Contexts {
			contexts = append(contexts, name)
		}
		sort.Strings(contexts)
		return contextsLoadedMsg{Contexts: contexts}
	}
}

func (m Model) switchContextCmd(contextName string) tea.Cmd {
	if contextName == "" {
		return nil
	}

	return func() tea.Msg {
		nextCfg := m.cfg
		nextCfg.Context = contextName
		nextCfg.Namespace = ""

		nextKube, err := clients.NewKubeFactory(nextCfg)
		if err != nil {
			return contextSwitchedMsg{Err: err}
		}

		namespace := nextKube.CurrentNamespace
		if namespace == "" {
			namespace = "default"
		}

		return contextSwitchedMsg{Context: contextName, Namespace: namespace, Kube: nextKube}
	}
}

func (m Model) loadDetailsCmd() tea.Cmd {
	name := m.selectedRelatedPodName()
	if name == "" {
		return nil
	}
	return func() tea.Msg {
		d, err := m.detailSvc.Pod(context.Background(), m.currentNamespace, name)
		if err != nil {
			return textLoadedMsg{Tab: "details", Err: err}
		}
		text := "Pod Details\n\n" +
			"Name: " + d.Name + "\n" +
			"Namespace: " + d.Namespace + "\n" +
			"Phase: " + d.Phase + "\n" +
			"Node: " + d.Node + "\n" +
			"PodIP: " + d.PodIP + "\n" +
			"HostIP: " + d.HostIP + "\n\n" +
			"Containers:\n- " + strings.Join(d.Containers, "\n- ") + "\n\n" +
			"Conditions:\n- " + strings.Join(d.Conditions, "\n- ")
		return textLoadedMsg{Tab: "details", Text: text}
	}
}

func (m Model) loadEventsCmd() tea.Cmd {
	name := m.selectedRelatedPodName()
	if name == "" {
		return nil
	}
	return func() tea.Msg {
		events, err := m.eventSvc.ForPod(context.Background(), m.currentNamespace, name)
		if err != nil {
			return textLoadedMsg{Tab: "events", Err: err}
		}
		if len(events) == 0 {
			return textLoadedMsg{Tab: "events", Text: "No events found."}
		}
		var b strings.Builder
		b.WriteString("Events\n\n")
		for _, e := range events {
			b.WriteString(e.Timestamp + " [" + e.Type + "] " + e.Reason + " - " + e.Message + "\n")
		}
		return textLoadedMsg{Tab: "events", Text: b.String()}
	}
}

func (m Model) loadYAMLCmd() tea.Cmd {
	name := m.selectedRelatedPodName()
	if name == "" {
		return nil
	}
	return func() tea.Msg {
		text, err := m.yamlSvc.Pod(context.Background(), m.currentNamespace, name)
		return textLoadedMsg{Tab: "yaml", Text: text, Err: err}
	}
}

func (m Model) loadLogsCmd() tea.Cmd {
	name := m.selectedRelatedPodName()
	if name == "" {
		return nil
	}
	return func() tea.Msg {
		text, err := m.logSvc.Pod(context.Background(), m.currentNamespace, name, "", 200)
		return textLoadedMsg{Tab: "logs", Text: text, Err: err}
	}
}

func (m Model) deletePodCmd() tea.Cmd {
	name := m.selectedPodName()
	if name == "" {
		return nil
	}
	return func() tea.Msg {
		err := m.deleteSvc.Pod(context.Background(), m.currentNamespace, name)
		return actionDoneMsg{Action: "delete", Err: err}
	}
}

func (m Model) patchPodCmd() tea.Cmd {
	name := m.selectedPodName()
	if name == "" {
		return nil
	}

	patchDoc := map[string]any{
		"metadata": map[string]any{
			"annotations": map[string]any{
				"deckhand.io/last-patch-at": time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	patchBytes, err := json.Marshal(patchDoc)
	if err != nil {
		return func() tea.Msg { return actionDoneMsg{Action: "patch", Err: fmt.Errorf("marshal patch: %w", err)} }
	}

	return func() tea.Msg {
		err := m.patchSvc.PodMergePatch(context.Background(), m.currentNamespace, name, patchBytes)
		return actionDoneMsg{Action: "patch", Err: err}
	}
}

func (m Model) cancelPipelineRunCmd() tea.Cmd {
	name := m.selectedPipelineRunName()
	if name == "" {
		return nil
	}
	return func() tea.Msg {
		err := m.prActionSvc.Cancel(context.Background(), m.currentNamespace, name)
		return actionDoneMsg{Action: "cancel", Err: err}
	}
}

func (m Model) refreshPermissionsCmd() tea.Cmd {
	if m.cfg.ReadOnly {
		return nil
	}

	return func() tea.Msg {
		deleteAllowed, deleteReason, deleteErr := m.authorizer.Allowed(context.Background(), m.currentNamespace, "", "pods", "delete")
		if deleteErr != nil {
			return permissionsLoadedMsg{Err: deleteErr}
		}

		patchAllowed, patchReason, patchErr := m.authorizer.Allowed(context.Background(), m.currentNamespace, "", "pods", "patch")
		if patchErr != nil {
			return permissionsLoadedMsg{Err: patchErr}
		}

		cancelAllowed, cancelReason, cancelErr := m.authorizer.Allowed(context.Background(), m.currentNamespace, "tekton.dev", "pipelineruns", "update")
		if cancelErr != nil {
			if errors.IsNotFound(cancelErr) || strings.Contains(strings.ToLower(cancelErr.Error()), "the server could not find the requested resource") {
				cancelAllowed = false
				cancelReason = "Tekton API not available"
			} else {
				return permissionsLoadedMsg{Err: cancelErr}
			}
		}

		return permissionsLoadedMsg{
			DeleteAllowed: deleteAllowed,
			DeleteReason:  deleteReason,
			PatchAllowed:  patchAllowed,
			PatchReason:   patchReason,
			CancelAllowed: cancelAllowed,
			CancelReason:  cancelReason,
		}
	}
}

func tickRefresh(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return t })
}

func (m Model) selectedPipelineRunName() string {
	if item, ok := m.list.SelectedItem().(panes.PipelineRunItem); ok {
		return item.PipelineRun.Name
	}
	return m.selectedPipeline
}

func (m Model) selectedRelatedPodName() string {
	if m.listScope == listScopeTaskRuns {
		if item, ok := m.list.SelectedItem().(panes.TaskRunItem); ok {
			return item.TaskRun.PodName
		}
	}
	return m.selectedPodName()
}
