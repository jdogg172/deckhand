package resources

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodDetails struct {
	Name       string
	Namespace  string
	Node       string
	Phase      string
	HostIP     string
	PodIP      string
	Conditions []string
	Containers []string
}

type DetailService struct {
	Client kubernetes.Interface
}

func NewDetailService(client kubernetes.Interface) *DetailService {
	return &DetailService{Client: client}
}

func (s *DetailService) Pod(ctx context.Context, namespace, name string) (PodDetails, error) {
	p, err := s.Client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return PodDetails{}, fmt.Errorf("get pod details: %w", err)
	}
	return summarizeDetails(*p), nil
}

func summarizeDetails(p corev1.Pod) PodDetails {
	conditions := make([]string, 0, len(p.Status.Conditions))
	for _, c := range p.Status.Conditions {
		conditions = append(conditions, fmt.Sprintf("%s=%s", c.Type, c.Status))
	}

	containers := make([]string, 0, len(p.Spec.Containers))
	for _, c := range p.Spec.Containers {
		containers = append(containers, c.Name)
	}

	return PodDetails{
		Name:       p.Name,
		Namespace:  p.Namespace,
		Node:       p.Spec.NodeName,
		Phase:      string(p.Status.Phase),
		HostIP:     p.Status.HostIP,
		PodIP:      p.Status.PodIP,
		Conditions: conditions,
		Containers: containers,
	}
}
