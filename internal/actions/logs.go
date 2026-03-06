package actions

import (
	"bytes"
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type LogService struct{ Client kubernetes.Interface }

func NewLogService(client kubernetes.Interface) *LogService { return &LogService{Client: client} }

func (s *LogService) Pod(ctx context.Context, namespace, name, container string, tailLines int64) (string, error) {
	opts := &corev1.PodLogOptions{Container: container, TailLines: &tailLines}
	req := s.Client.CoreV1().Pods(namespace).GetLogs(name, opts)
	rc, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("stream logs: %w", err)
	}
	defer rc.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, rc); err != nil {
		return "", fmt.Errorf("read logs: %w", err)
	}
	return buf.String(), nil
}
