package resources

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var (
	pipelineRunGVR = schema.GroupVersionResource{Group: "tekton.dev", Version: "v1", Resource: "pipelineruns"}
	taskRunGVR     = schema.GroupVersionResource{Group: "tekton.dev", Version: "v1", Resource: "taskruns"}
)

type PipelineRunSummary struct {
	Name      string
	Namespace string
	Status    string
	Reason    string
	Duration  string
	StartTime string
	Highlight bool
}

type TaskRunSummary struct {
	Name      string
	Status    string
	Reason    string
	Duration  string
	PodName   string
	Highlight bool
}

type PipelineService struct {
	DynamicClient dynamic.Interface
	HasTektonAPI  bool
	Now           func() time.Time
}

func NewPipelineService(dynamicClient dynamic.Interface, hasTektonAPI bool) *PipelineService {
	return &PipelineService{DynamicClient: dynamicClient, HasTektonAPI: hasTektonAPI, Now: time.Now}
}

func (s *PipelineService) ListPipelineRuns(ctx context.Context, namespace string) ([]PipelineRunSummary, error) {
	if !s.HasTektonAPI {
		return nil, ErrAPINotAvailable
	}

	prs, err := s.DynamicClient.Resource(pipelineRunGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pipelineruns: %w", err)
	}

	out := make([]PipelineRunSummary, 0, len(prs.Items))
	for _, pr := range prs.Items {
		out = append(out, summarizePipelineRun(pr, s.Now()))
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (s *PipelineService) ListTaskRunsForPipelineRun(ctx context.Context, namespace, pipelineRunName string) ([]TaskRunSummary, error) {
	if !s.HasTektonAPI {
		return nil, ErrAPINotAvailable
	}

	trs, err := s.DynamicClient.Resource(taskRunGVR).Namespace(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("tekton.dev/pipelineRun=%s", pipelineRunName),
	})
	if err != nil {
		return nil, fmt.Errorf("list taskruns: %w", err)
	}

	out := make([]TaskRunSummary, 0, len(trs.Items))
	for _, tr := range trs.Items {
		out = append(out, summarizeTaskRun(tr, s.Now()))
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func summarizePipelineRun(pr unstructured.Unstructured, now time.Time) PipelineRunSummary {
	status, reason := tektonCondition(pr.Object)
	startTimeRaw, _, _ := unstructured.NestedString(pr.Object, "status", "startTime")
	completionTimeRaw, _, _ := unstructured.NestedString(pr.Object, "status", "completionTime")

	duration := durationString(startTimeRaw, completionTimeRaw, now)
	start := "-"
	if startTimeRaw != "" {
		start = startTimeRaw
	}

	highlight := isFailedOrStuck(status, reason, duration)
	return PipelineRunSummary{
		Name:      pr.GetName(),
		Namespace: pr.GetNamespace(),
		Status:    status,
		Reason:    reason,
		Duration:  duration,
		StartTime: start,
		Highlight: highlight,
	}
}

func summarizeTaskRun(tr unstructured.Unstructured, now time.Time) TaskRunSummary {
	status, reason := tektonCondition(tr.Object)
	startTimeRaw, _, _ := unstructured.NestedString(tr.Object, "status", "startTime")
	completionTimeRaw, _, _ := unstructured.NestedString(tr.Object, "status", "completionTime")
	duration := durationString(startTimeRaw, completionTimeRaw, now)
	podName, _, _ := unstructured.NestedString(tr.Object, "status", "podName")
	highlight := isFailedOrStuck(status, reason, duration)

	return TaskRunSummary{
		Name:      tr.GetName(),
		Status:    status,
		Reason:    reason,
		Duration:  duration,
		PodName:   podName,
		Highlight: highlight,
	}
}

func tektonCondition(obj map[string]any) (string, string) {
	conditions, found, _ := unstructured.NestedSlice(obj, "status", "conditions")
	if !found || len(conditions) == 0 {
		return "Unknown", "NoConditions"
	}

	cond, ok := conditions[0].(map[string]any)
	if !ok {
		return "Unknown", "InvalidCondition"
	}

	status, _, _ := unstructured.NestedString(cond, "status")
	reason, _, _ := unstructured.NestedString(cond, "reason")
	if status == "" {
		status = "Unknown"
	}
	if reason == "" {
		reason = "Unknown"
	}
	return status, reason
}

func durationString(startTimeRaw, completionTimeRaw string, now time.Time) string {
	if startTimeRaw == "" {
		return "-"
	}

	start, err := time.Parse(time.RFC3339, startTimeRaw)
	if err != nil {
		return "-"
	}

	end := now
	if completionTimeRaw != "" {
		if completion, err := time.Parse(time.RFC3339, completionTimeRaw); err == nil {
			end = completion
		}
	}

	d := end.Sub(start).Round(time.Second)
	if d < 0 {
		d = 0
	}
	return d.String()
}

func isFailedOrStuck(status, reason, duration string) bool {
	if strings.EqualFold(status, "False") || strings.Contains(strings.ToLower(reason), "failed") {
		return true
	}
	if strings.EqualFold(status, "Unknown") && duration != "-" {
		return true
	}
	if strings.Contains(strings.ToLower(reason), "timeout") || strings.Contains(strings.ToLower(reason), "cancel") {
		return true
	}
	return false
}
