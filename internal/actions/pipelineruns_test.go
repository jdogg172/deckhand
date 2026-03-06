package actions

import (
	"context"
	"errors"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestPipelineRunActionServiceCancel_TektonUnavailable(t *testing.T) {
	svc := NewPipelineRunActionService(nil, false)
	err := svc.Cancel(context.Background(), "ci", "pr-1")
	if err == nil {
		t.Fatalf("expected error when Tekton API is unavailable")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "tekton api not available") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPipelineRunActionServiceCancel_GetFailure(t *testing.T) {
	client := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme())
	svc := NewPipelineRunActionService(client, true)

	err := svc.Cancel(context.Background(), "ci", "missing")
	if err == nil {
		t.Fatalf("expected get failure for missing PipelineRun")
	}
	if !strings.Contains(err.Error(), "get pipelinerun ci/missing") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPipelineRunActionServiceCancel_UpdateFailure(t *testing.T) {
	obj := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "tekton.dev/v1",
		"kind":       "PipelineRun",
		"metadata": map[string]any{
			"name":      "pr-1",
			"namespace": "ci",
		},
		"spec": map[string]any{},
	}}

	client := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), obj)
	client.PrependReactor("update", "pipelineruns", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("update boom")
	})

	svc := NewPipelineRunActionService(client, true)
	err := svc.Cancel(context.Background(), "ci", "pr-1")
	if err == nil {
		t.Fatalf("expected update failure")
	}
	if !strings.Contains(err.Error(), "cancel pipelinerun ci/pr-1") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPipelineRunActionServiceCancel_SetsCancelledStatus(t *testing.T) {
	obj := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "tekton.dev/v1",
		"kind":       "PipelineRun",
		"metadata": map[string]any{
			"name":      "pr-2",
			"namespace": "ci",
		},
		"spec": map[string]any{},
	}}

	client := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), obj)
	svc := NewPipelineRunActionService(client, true)

	if err := svc.Cancel(context.Background(), "ci", "pr-2"); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	updated, err := client.Resource(pipelineRunGVR).Namespace("ci").Get(context.Background(), "pr-2", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get updated PipelineRun: %v", err)
	}

	status, found, err := unstructured.NestedString(updated.Object, "spec", "status")
	if err != nil {
		t.Fatalf("failed to read status: %v", err)
	}
	if !found || status != "Cancelled" {
		t.Fatalf("expected spec.status=Cancelled, got found=%t value=%q", found, status)
	}
}
