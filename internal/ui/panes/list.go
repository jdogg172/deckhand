package panes

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/jdogg172/deckhand/internal/resources"
)

type PodItem struct{ Pod resources.PodSummary }

func (i PodItem) FilterValue() string { return i.Pod.Name }
func (i PodItem) Title() string       { return i.Pod.Name }
func (i PodItem) Description() string {
	return fmt.Sprintf("%s  ready=%s  restarts=%d  node=%s", i.Pod.Status, i.Pod.Ready, i.Pod.Restarts, i.Pod.Node)
}

func PodItemsFromSummaries(summaries []resources.PodSummary) []list.Item {
	items := make([]list.Item, 0, len(summaries))
	for _, s := range summaries {
		items = append(items, PodItem{Pod: s})
	}
	return items
}

type NamespaceItem struct {
	Name      string
	IsProject bool
}

func (i NamespaceItem) FilterValue() string { return i.Name }
func (i NamespaceItem) Title() string       { return i.Name }
func (i NamespaceItem) Description() string {
	if i.IsProject {
		return "OpenShift project"
	}
	return "Kubernetes namespace"
}

func NamespaceItems(names []string, isProject bool) []list.Item {
	items := make([]list.Item, 0, len(names))
	for _, name := range names {
		items = append(items, NamespaceItem{Name: name, IsProject: isProject})
	}
	return items
}

type ContextItem struct {
	Name    string
	Current bool
}

func (i ContextItem) FilterValue() string { return i.Name }
func (i ContextItem) Title() string {
	if i.Current {
		return "* " + i.Name
	}
	return i.Name
}
func (i ContextItem) Description() string {
	if i.Current {
		return "current context"
	}
	return "kube context"
}

func ContextItems(names []string, current string) []list.Item {
	items := make([]list.Item, 0, len(names))
	for _, name := range names {
		items = append(items, ContextItem{Name: name, Current: name == current})
	}
	return items
}

type RouteItem struct{ Route resources.RouteSummary }

func (i RouteItem) FilterValue() string {
	return i.Route.Name + " " + i.Route.Host + " " + i.Route.ToName
}
func (i RouteItem) Title() string { return i.Route.Name }
func (i RouteItem) Description() string {
	path := i.Route.Path
	if path == "" {
		path = "/"
	}
	return fmt.Sprintf("host=%s%s -> %s/%s tls=%s", i.Route.Host, path, i.Route.ToKind, i.Route.ToName, i.Route.TLS)
}

func RouteItemsFromSummaries(summaries []resources.RouteSummary) []list.Item {
	items := make([]list.Item, 0, len(summaries))
	for _, s := range summaries {
		items = append(items, RouteItem{Route: s})
	}
	return items
}

type PipelineRunItem struct{ PipelineRun resources.PipelineRunSummary }

func (i PipelineRunItem) FilterValue() string {
	return i.PipelineRun.Name + " " + i.PipelineRun.Status + " " + i.PipelineRun.Reason
}
func (i PipelineRunItem) Title() string {
	if i.PipelineRun.Highlight {
		return "! " + i.PipelineRun.Name
	}
	return i.PipelineRun.Name
}
func (i PipelineRunItem) Description() string {
	return fmt.Sprintf("status=%s reason=%s duration=%s", i.PipelineRun.Status, i.PipelineRun.Reason, i.PipelineRun.Duration)
}

func PipelineRunItemsFromSummaries(summaries []resources.PipelineRunSummary) []list.Item {
	items := make([]list.Item, 0, len(summaries))
	for _, s := range summaries {
		items = append(items, PipelineRunItem{PipelineRun: s})
	}
	return items
}

type TaskRunItem struct{ TaskRun resources.TaskRunSummary }

func (i TaskRunItem) FilterValue() string {
	return i.TaskRun.Name + " " + i.TaskRun.Status + " " + i.TaskRun.PodName
}
func (i TaskRunItem) Title() string {
	if i.TaskRun.Highlight {
		return "! " + i.TaskRun.Name
	}
	return i.TaskRun.Name
}
func (i TaskRunItem) Description() string {
	return fmt.Sprintf("status=%s reason=%s duration=%s pod=%s", i.TaskRun.Status, i.TaskRun.Reason, i.TaskRun.Duration, i.TaskRun.PodName)
}

func TaskRunItemsFromSummaries(summaries []resources.TaskRunSummary) []list.Item {
	items := make([]list.Item, 0, len(summaries))
	for _, s := range summaries {
		items = append(items, TaskRunItem{TaskRun: s})
	}
	return items
}
