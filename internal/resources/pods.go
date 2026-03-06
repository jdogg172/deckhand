package resources

import (
	"context"
	"fmt"
	"sort"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodSummary struct {
	Name         string
	Namespace    string
	Status       string
	Ready        string
	Restarts     int32
	Node         string
	AgeTimestamp metav1.Time
}

type PodService struct {
	Client kubernetes.Interface
}

func NewPodService(client kubernetes.Interface) *PodService {
	return &PodService{Client: client}
}

func (s *PodService) List(ctx context.Context, namespace string) ([]PodSummary, error) {
	pods, err := s.Client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	out := make([]PodSummary, 0, len(pods.Items))
	for _, p := range pods.Items {
		out = append(out, summarizePod(p))
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func summarizePod(p corev1.Pod) PodSummary {
	var restarts int32
	var readyCount int
	for _, cs := range p.Status.ContainerStatuses {
		restarts += cs.RestartCount
		if cs.Ready {
			readyCount++
		}
	}

	status := string(p.Status.Phase)
	for _, cs := range p.Status.ContainerStatuses {
		if cs.State.Waiting != nil && cs.State.Waiting.Reason != "" {
			status = cs.State.Waiting.Reason
			break
		}
		if cs.State.Terminated != nil && cs.State.Terminated.Reason != "" {
			status = cs.State.Terminated.Reason
			break
		}
	}

	return PodSummary{
		Name:         p.Name,
		Namespace:    p.Namespace,
		Status:       status,
		Ready:        fmt.Sprintf("%d/%d", readyCount, len(p.Status.ContainerStatuses)),
		Restarts:     restarts,
		Node:         p.Spec.NodeName,
		AgeTimestamp: p.CreationTimestamp,
	}
}
