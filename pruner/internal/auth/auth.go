package auth

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// CreateKubernetesClient creates a new Kubernetes client using in-cluster configuration.
// It returns a pointer to the Kubernetes clientset or an error if the creation fails.
//
// Parameters:
// - log: A logger instance for logging errors.
//
// Returns:
// - A pointer to the Kubernetes clientset.
// - An error if the clientset could not be created.
func CreateKubernetesClient(log *logrus.Logger) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster Kubernetes config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create client set for in-cluster Kubernetes config: %w", err)
	}
	if clientset == nil {
		return nil, fmt.Errorf("Kubernetes client set cannot be nil")
	}
	return clientset, nil
}
