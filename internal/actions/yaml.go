package actions

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

type YAMLService struct{ Client kubernetes.Interface }

func NewYAMLService(client kubernetes.Interface) *YAMLService { return &YAMLService{Client: client} }

func (s *YAMLService) Pod(ctx context.Context, namespace, name string) (string, error) {
	p, err := s.Client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("get pod for yaml: %w", err)
	}
	b, err := yaml.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("marshal yaml: %w", err)
	}
	return string(b), nil
}
