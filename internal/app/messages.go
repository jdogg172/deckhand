package app

import (
	"github.com/example/deckhand/internal/clients"
	"github.com/example/deckhand/internal/resources"
)

type podsLoadedMsg struct {
	Pods []resources.PodSummary
	Err  error
}

type routesLoadedMsg struct {
	Routes []resources.RouteSummary
	Err    error
}

type pipelinesLoadedMsg struct {
	PipelineRuns []resources.PipelineRunSummary
	Err          error
}

type taskRunsLoadedMsg struct {
	TaskRuns []resources.TaskRunSummary
	Err      error
}

type namespacesLoadedMsg struct {
	Namespaces []string
	IsProject  bool
	Err        error
}

type contextsLoadedMsg struct {
	Contexts []string
	Err      error
}

type contextSwitchedMsg struct {
	Context   string
	Namespace string
	Kube      *clients.KubeFactory
	Err       error
}

type textLoadedMsg struct {
	Tab  string
	Text string
	Err  error
}

type actionDoneMsg struct {
	Action string
	Err    error
}

type permissionsLoadedMsg struct {
	DeleteAllowed bool
	DeleteReason  string
	PatchAllowed  bool
	PatchReason   string
	CancelAllowed bool
	CancelReason  string
	Err           error
}
