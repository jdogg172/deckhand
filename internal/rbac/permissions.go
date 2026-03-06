package rbac

import (
	"context"
	"fmt"

	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Authorizer struct {
	Client kubernetes.Interface
}

func NewAuthorizer(client kubernetes.Interface) *Authorizer {
	return &Authorizer{Client: client}
}

func (a *Authorizer) Allowed(ctx context.Context, namespace, group, resource, verb string) (bool, string, error) {
	review, err := a.Client.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace: namespace,
				Group:     group,
				Resource:  resource,
				Verb:      verb,
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return false, "", fmt.Errorf("check access %s %s/%s: %w", verb, group, resource, err)
	}

	if review.Status.Allowed {
		return true, "", nil
	}

	if review.Status.Reason != "" {
		return false, review.Status.Reason, nil
	}
	if review.Status.EvaluationError != "" {
		return false, review.Status.EvaluationError, nil
	}
	return false, "action not permitted by current RBAC", nil
}
