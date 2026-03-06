package clients

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/deckhand/internal/config"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type KubeFactory struct {
	RestConfig       *rest.Config
	Clientset        kubernetes.Interface
	Discovery        discovery.DiscoveryInterface
	Dynamic          dynamic.Interface
	RawConfig        clientcmdapi.Config
	CurrentContext   string
	CurrentNamespace string
	HasTektonAPI     bool
	HasOpenShiftAPI  bool
	HasRouteAPI      bool
}

func NewKubeFactory(cfg config.Config) (*KubeFactory, error) {
	kubeconfigPath := cfg.Kubeconfig
	if kubeconfigPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("resolve user home: %w", err)
		}
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	overrides := &clientcmd.ConfigOverrides{}
	if cfg.Context != "" {
		overrides.CurrentContext = cfg.Context
	}
	if cfg.Namespace != "" {
		overrides.Context.Namespace = cfg.Namespace
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
	restCfg, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("build rest config: %w", err)
	}

	rawCfg, err := clientConfig.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("load raw kubeconfig: %w", err)
	}

	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		return nil, fmt.Errorf("resolve current namespace: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("build kubernetes clientset: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("build dynamic client: %w", err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("build discovery client: %w", err)
	}

	currentContext := cfg.Context
	if currentContext == "" {
		currentContext = rawCfg.CurrentContext
	}

	currentNamespace := cfg.Namespace
	if currentNamespace == "" {
		currentNamespace = namespace
	}

	hasTekton := supportsGroupVersion(discoveryClient, schema.GroupVersion{Group: "tekton.dev", Version: "v1"}.String())
	hasOpenShift := supportsGroupVersion(discoveryClient, schema.GroupVersion{Group: "project.openshift.io", Version: "v1"}.String())
	hasRoute := supportsGroupVersion(discoveryClient, schema.GroupVersion{Group: "route.openshift.io", Version: "v1"}.String())

	return &KubeFactory{
		RestConfig:       restCfg,
		Clientset:        clientset,
		Discovery:        discoveryClient,
		Dynamic:          dynamicClient,
		RawConfig:        rawCfg,
		CurrentContext:   currentContext,
		CurrentNamespace: currentNamespace,
		HasTektonAPI:     hasTekton,
		HasOpenShiftAPI:  hasOpenShift,
		HasRouteAPI:      hasRoute,
	}, nil
}

func supportsGroupVersion(discoveryClient discovery.DiscoveryInterface, groupVersion string) bool {
	if _, err := discoveryClient.ServerResourcesForGroupVersion(groupVersion); err != nil {
		return false
	}
	return true
}
