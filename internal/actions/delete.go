package actions

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeleteService struct{ Client kubernetes.Interface }

func NewDeleteService(client kubernetes.Interface) *DeleteService {
	return &DeleteService{Client: client}
}

func (s *DeleteService) Pod(ctx context.Context, namespace, name string) error {
	if err := s.Client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("delete pod %s/%s: %w", namespace, name, err)
	}
	return nil
}
