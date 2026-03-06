package actions

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type PatchService struct{ Client kubernetes.Interface }

func NewPatchService(client kubernetes.Interface) *PatchService { return &PatchService{Client: client} }

func (s *PatchService) PodMergePatch(ctx context.Context, namespace, name string, patch []byte) error {
	if _, err := s.Client.CoreV1().Pods(namespace).Patch(ctx, name, types.MergePatchType, patch, metav1.PatchOptions{}); err != nil {
		return fmt.Errorf("patch pod %s/%s: %w", namespace, name, err)
	}
	return nil
}
