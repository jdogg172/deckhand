package resources

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSummarizePipelineRun_FailedHighlight(t *testing.T) {
	now := time.Date(2026, 3, 6, 12, 0, 0, 0, time.UTC)
	pr := unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{"name": "build-app", "namespace": "ci"},
		"status": map[string]any{
			"startTime":      "2026-03-06T11:58:00Z",
			"completionTime": "2026-03-06T12:00:00Z",
			"conditions":     []any{map[string]any{"status": "False", "reason": "Failed"}},
		},
	}}

	s := summarizePipelineRun(pr, now)
	if !s.Highlight {
		t.Fatalf("expected failed PipelineRun to be highlighted")
	}
	if s.Duration == "-" {
		t.Fatalf("expected duration to be calculated")
	}
}

func TestSummarizeTaskRun_PodResolution(t *testing.T) {
	now := time.Date(2026, 3, 6, 12, 0, 0, 0, time.UTC)
	tr := unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{"name": "build-app-task"},
		"status": map[string]any{
			"startTime":  "2026-03-06T11:59:30Z",
			"podName":    "build-app-task-pod",
			"conditions": []any{map[string]any{"status": "True", "reason": "Succeeded"}},
		},
	}}

	s := summarizeTaskRun(tr, now)
	if s.PodName != "build-app-task-pod" {
		t.Fatalf("expected pod resolution, got %q", s.PodName)
	}
	if s.Highlight {
		t.Fatalf("did not expect successful TaskRun to be highlighted")
	}
}
