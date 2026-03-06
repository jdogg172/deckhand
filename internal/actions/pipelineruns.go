package actions

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var pipelineRunGVR = schema.GroupVersionResource{Group: "tekton.dev", Version: "v1", Resource: "pipelineruns"}

type PipelineRunActionService struct {
	DynamicClient dynamic.Interface
	HasTektonAPI  bool
}

func NewPipelineRunActionService(dynamicClient dynamic.Interface, hasTektonAPI bool) *PipelineRunActionService {
	return &PipelineRunActionService{DynamicClient: dynamicClient, HasTektonAPI: hasTektonAPI}
}

func (s *PipelineRunActionService) Cancel(ctx context.Context, namespace, name string) error {
	if !s.HasTektonAPI {
		return fmt.Errorf("tekton api not available")
	}

	pr, err := s.DynamicClient.Resource(pipelineRunGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get pipelinerun %s/%s: %w", namespace, name, err)
	}

	if err := unstructured.SetNestedField(pr.Object, "Cancelled", "spec", "status"); err != nil {
		return fmt.Errorf("set pipelinerun status: %w", err)
	}

	if _, err := s.DynamicClient.Resource(pipelineRunGVR).Namespace(namespace).Update(ctx, pr, metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("cancel pipelinerun %s/%s: %w", namespace, name, err)
	}

	return nil
}
