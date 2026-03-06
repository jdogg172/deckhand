package app

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jdogg172/deckhand/internal/actions"
	"github.com/jdogg172/deckhand/internal/clients"
	"github.com/jdogg172/deckhand/internal/config"
	"github.com/jdogg172/deckhand/internal/rbac"
	"github.com/jdogg172/deckhand/internal/resources"
	"github.com/jdogg172/deckhand/internal/ui/modes"
	"github.com/jdogg172/deckhand/internal/ui/panes"
	"github.com/jdogg172/deckhand/internal/ui/styles"
)

type listScope string

const (
	listScopePods       listScope = "pods"
	listScopeNamespaces listScope = "namespaces"
	listScopeContexts   listScope = "contexts"
	listScopeRoutes     listScope = "routes"
	listScopePipelines  listScope = "pipelineruns"
	listScopeTaskRuns   listScope = "taskruns"
)

type pendingConfirmation struct {
	Action  string
	Target  string
	Message string
}

type permissionState struct {
	DeleteAllowed bool
	DeleteReason  string
	PatchAllowed  bool
	PatchReason   string
	CancelAllowed bool
	CancelReason  string
}

type Model struct {
	ctx          context.Context
	cfg          config.Config
	kube         *clients.KubeFactory
	podSvc       *resources.PodService
	eventSvc     *resources.EventService
	detailSvc    *resources.DetailService
	namespaceSvc *resources.NamespaceService
	routeSvc     *resources.RouteService
	pipelineSvc  *resources.PipelineService
	yamlSvc      *actions.YAMLService
	logSvc       *actions.LogService
	deleteSvc    *actions.DeleteService
	patchSvc     *actions.PatchService
	prActionSvc  *actions.PipelineRunActionService
	authorizer   *rbac.Authorizer

	keys   KeyMap
	help   help.Model
	list   list.Model
	theme  styles.Theme
	mode   modes.Mode
	width  int
	height int

	activeTab  string
	detailText string
	statusText string
	errText    string

	listScope         listScope
	previousListScope listScope
	currentNamespace  string
	selectedContext   string
	selectedPipeline  string
	permissions       permissionState
	pendingConfirm    *pendingConfirmation
}

func NewModel(ctx context.Context, cfg config.Config, kube *clients.KubeFactory) Model {
	keys := DefaultKeyMap()
	h := help.New()
	h.ShowAll = false

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Pods"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowPagination(true)
	l.Styles.Title = lipgloss.NewStyle().Bold(true)

	return Model{
		ctx:               ctx,
		cfg:               cfg,
		kube:              kube,
		podSvc:            resources.NewPodService(kube.Clientset),
		eventSvc:          resources.NewEventService(kube.Clientset),
		detailSvc:         resources.NewDetailService(kube.Clientset),
		namespaceSvc:      resources.NewNamespaceService(kube.Clientset, kube.Dynamic, kube.HasOpenShiftAPI),
		routeSvc:          resources.NewRouteService(kube.Dynamic, kube.HasRouteAPI),
		pipelineSvc:       resources.NewPipelineService(kube.Dynamic, kube.HasTektonAPI),
		yamlSvc:           actions.NewYAMLService(kube.Clientset),
		logSvc:            actions.NewLogService(kube.Clientset),
		deleteSvc:         actions.NewDeleteService(kube.Clientset),
		patchSvc:          actions.NewPatchService(kube.Clientset),
		prActionSvc:       actions.NewPipelineRunActionService(kube.Dynamic, kube.HasTektonAPI),
		authorizer:        rbac.NewAuthorizer(kube.Clientset),
		keys:              keys,
		help:              h,
		list:              l,
		theme:             styles.DefaultTheme(cfg.NoColor),
		mode:              modes.Normalize(cfg.Mode),
		activeTab:         "details",
		detailText:        "Loading...",
		statusText:        "Ready",
		listScope:         listScopePods,
		previousListScope: listScopePods,
		currentNamespace:  cfg.Namespace,
		selectedContext:   cfg.Context,
	}
}

func (m Model) Init() tea.Cmd {
	if m.currentNamespace == "" {
		m.currentNamespace = "default"
	}
	return tea.Batch(
		m.refreshCurrentScopeCmd(),
		m.refreshPermissionsCmd(),
		tickRefresh(time.Duration(m.cfg.UI.RefreshIntervalSeconds)*time.Second),
	)
}

func (m Model) selectedPodName() string {
	if item, ok := m.list.SelectedItem().(panes.PodItem); ok {
		return item.Pod.Name
	}
	return ""
}

func (m Model) headerView() string {
	mode := string(m.mode)
	ro := "RW"
	if m.cfg.ReadOnly {
		ro = "RO"
	}
	text := fmt.Sprintf("Deckhand | ctx=%s | ns=%s | mode=%s | %s", m.cfg.Context, m.currentNamespace, mode, ro)
	return m.theme.Header.Render(text)
}

func (m Model) footerView() string {
	actionHints := ""
	if m.cfg.ReadOnly {
		actionHints = " | readonly: x/p/k disabled"
	} else {
		hints := make([]string, 0, 3)
		if !m.permissions.DeleteAllowed {
			hints = append(hints, "x disabled")
		}
		if !m.permissions.PatchAllowed {
			hints = append(hints, "p disabled")
		}
		if !m.permissions.CancelAllowed {
			hints = append(hints, "k disabled")
		}
		if len(hints) > 0 {
			actionHints = " | " + fmt.Sprintf("%s", hints[0])
			for i := 1; i < len(hints); i++ {
				actionHints += ", " + hints[i]
			}
		}
	}

	return m.theme.Footer.Render("1 ops  2 pipeline  c context  n namespace/project  u routes  t taskruns  o pod  d details  l logs  y yaml  e events  p patch  x delete  k cancel  r refresh  q quit" + actionHints)
}

func (m Model) withKube(kube *clients.KubeFactory) Model {
	m.kube = kube
	m.podSvc = resources.NewPodService(kube.Clientset)
	m.eventSvc = resources.NewEventService(kube.Clientset)
	m.detailSvc = resources.NewDetailService(kube.Clientset)
	m.namespaceSvc = resources.NewNamespaceService(kube.Clientset, kube.Dynamic, kube.HasOpenShiftAPI)
	m.routeSvc = resources.NewRouteService(kube.Dynamic, kube.HasRouteAPI)
	m.pipelineSvc = resources.NewPipelineService(kube.Dynamic, kube.HasTektonAPI)
	m.yamlSvc = actions.NewYAMLService(kube.Clientset)
	m.logSvc = actions.NewLogService(kube.Clientset)
	m.deleteSvc = actions.NewDeleteService(kube.Clientset)
	m.patchSvc = actions.NewPatchService(kube.Clientset)
	m.prActionSvc = actions.NewPipelineRunActionService(kube.Dynamic, kube.HasTektonAPI)
	m.authorizer = rbac.NewAuthorizer(kube.Clientset)
	return m
}
