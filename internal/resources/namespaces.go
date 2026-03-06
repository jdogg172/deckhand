package resources

import (
	"context"
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var projectGVR = schema.GroupVersionResource{Group: "project.openshift.io", Version: "v1", Resource: "projects"}

type NamespaceService struct {
	Client          kubernetes.Interface
	DynamicClient   dynamic.Interface
	HasOpenShiftAPI bool
}

func NewNamespaceService(client kubernetes.Interface, dynamicClient dynamic.Interface, hasOpenShiftAPI bool) *NamespaceService {
	return &NamespaceService{Client: client, DynamicClient: dynamicClient, HasOpenShiftAPI: hasOpenShiftAPI}
}

func (s *NamespaceService) List(ctx context.Context) ([]string, bool, error) {
	if s.HasOpenShiftAPI {
		projects, err := s.DynamicClient.Resource(projectGVR).List(ctx, metav1.ListOptions{})
		if err == nil {
			names := namesFromUnstructured(projects.Items)
			sort.Strings(names)
			return names, true, nil
		}
	}

	nss, err := s.Client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, false, fmt.Errorf("list namespaces: %w", err)
	}

	names := make([]string, 0, len(nss.Items))
	for _, ns := range nss.Items {
		names = append(names, ns.Name)
	}
	sort.Strings(names)
	return names, false, nil
}

func namesFromUnstructured(items []unstructured.Unstructured) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item.GetName() != "" {
			out = append(out, item.GetName())
		}
	}
	return out
}
