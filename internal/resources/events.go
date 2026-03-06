package resources

import (
	"context"
	"fmt"
	"sort"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type EventSummary struct {
	Type      string
	Reason    string
	Message   string
	Timestamp string
}

type EventService struct {
	Client kubernetes.Interface
}

func NewEventService(client kubernetes.Interface) *EventService {
	return &EventService{Client: client}
}

func (s *EventService) ForPod(ctx context.Context, namespace, podName string) ([]EventSummary, error) {
	evts, err := s.Client.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", podName),
	})
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	sort.Slice(evts.Items, func(i, j int) bool { return eventTime(evts.Items[i]).Before(eventTime(evts.Items[j])) })

	out := make([]EventSummary, 0, len(evts.Items))
	for _, e := range evts.Items {
		out = append(out, EventSummary{
			Type:      e.Type,
			Reason:    e.Reason,
			Message:   e.Message,
			Timestamp: eventTime(e).Format("2006-01-02 15:04:05"),
		})
	}
	return out, nil
}

func eventTime(e corev1.Event) time.Time {
	if !e.LastTimestamp.IsZero() {
		return e.LastTimestamp.Time
	}
	if !e.EventTime.IsZero() {
		return e.EventTime.Time
	}
	return e.FirstTimestamp.Time
}
