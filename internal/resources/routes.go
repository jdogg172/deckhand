package resources

import (
	"context"
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var routeGVR = schema.GroupVersionResource{Group: "route.openshift.io", Version: "v1", Resource: "routes"}

type RouteSummary struct {
	Name   string
	Host   string
	Path   string
	ToKind string
	ToName string
	TLS    string
}

type RouteService struct {
	DynamicClient dynamic.Interface
	HasRouteAPI   bool
}

func NewRouteService(dynamicClient dynamic.Interface, hasRouteAPI bool) *RouteService {
	return &RouteService{DynamicClient: dynamicClient, HasRouteAPI: hasRouteAPI}
}

func (s *RouteService) List(ctx context.Context, namespace string) ([]RouteSummary, error) {
	if !s.HasRouteAPI {
		return nil, ErrAPINotAvailable
	}

	list, err := s.DynamicClient.Resource(routeGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list routes: %w", err)
	}

	out := make([]RouteSummary, 0, len(list.Items))
	for _, item := range list.Items {
		out = append(out, summarizeRoute(item))
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func summarizeRoute(item unstructured.Unstructured) RouteSummary {
	host, _, _ := unstructured.NestedString(item.Object, "spec", "host")
	path, _, _ := unstructured.NestedString(item.Object, "spec", "path")
	toKind, _, _ := unstructured.NestedString(item.Object, "spec", "to", "kind")
	toName, _, _ := unstructured.NestedString(item.Object, "spec", "to", "name")

	tlsTermination := "none"
	if termination, ok, _ := unstructured.NestedString(item.Object, "spec", "tls", "termination"); ok && termination != "" {
		tlsTermination = termination
	}

	return RouteSummary{
		Name:   item.GetName(),
		Host:   host,
		Path:   path,
		ToKind: toKind,
		ToName: toName,
		TLS:    tlsTermination,
	}
}
