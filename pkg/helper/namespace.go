package helper

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateNamespaceIfNotExist creates a namespace in case it doesn't exist.
func CreateNamespaceIfNotExist(ctx context.Context, k8sClient client.Client, namespace string) error {
	if namespace == "" {
		return nil // objects with no namespace such as ManagedClusterSet
	}

	namespaceObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{ Name: namespace },
	}

	if err := k8sClient.Create(ctx, namespaceObj); err != nil && !strings.HasSuffix(err.Error(), "already exists") {
		return fmt.Errorf("failed to create namespace - %w", err)
	}

	return nil
}

