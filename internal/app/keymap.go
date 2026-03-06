package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	OpsMode      key.Binding
	PipelineMode key.Binding
	Filter       key.Binding
	Namespace    key.Binding
	Context      key.Binding
	Refresh      key.Binding
	Help         key.Binding
	Quit         key.Binding
	Logs         key.Binding
	YAML         key.Binding
	Events       key.Binding
	Details      key.Binding
	Routes       key.Binding
	Delete       key.Binding
	Patch        key.Binding
	TaskRuns     key.Binding
	OpenPod      key.Binding
	Cancel       key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		OpsMode:      key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "ops mode")),
		PipelineMode: key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "pipeline mode")),
		Filter:       key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Namespace:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "namespace")),
		Context:      key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "context")),
		Refresh:      key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
		Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:         key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Logs:         key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logs")),
		YAML:         key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "yaml")),
		Events:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "events")),
		Details:      key.NewBinding(key.WithKeys("d", "enter"), key.WithHelp("d/enter", "details")),
		Routes:       key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "routes")),
		Delete:       key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete")),
		Patch:        key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "patch")),
		TaskRuns:     key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "taskruns")),
		OpenPod:      key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "related pod")),
		Cancel:       key.NewBinding(key.WithKeys("k"), key.WithHelp("k", "cancel pipelinerun")),
	}
}
