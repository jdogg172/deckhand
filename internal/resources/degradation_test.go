package resources

import (
	"context"
	"errors"
	"testing"
)

func TestRouteServiceList_APINotAvailable(t *testing.T) {
	svc := NewRouteService(nil, false)
	_, err := svc.List(context.Background(), "default")
	if !errors.Is(err, ErrAPINotAvailable) {
		t.Fatalf("expected ErrAPINotAvailable, got %v", err)
	}
}

func TestPipelineServiceList_APINotAvailable(t *testing.T) {
	svc := NewPipelineService(nil, false)

	if _, err := svc.ListPipelineRuns(context.Background(), "default"); !errors.Is(err, ErrAPINotAvailable) {
		t.Fatalf("expected ErrAPINotAvailable from ListPipelineRuns, got %v", err)
	}

	if _, err := svc.ListTaskRunsForPipelineRun(context.Background(), "default", "any"); !errors.Is(err, ErrAPINotAvailable) {
		t.Fatalf("expected ErrAPINotAvailable from ListTaskRunsForPipelineRun, got %v", err)
	}
}
