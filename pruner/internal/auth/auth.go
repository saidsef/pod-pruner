package auth

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// KubernetesClientManager manages the Kubernetes client creation and caching.
type KubernetesClientManager struct {
	clientset *kubernetes.Clientset
	once      sync.Once
	log       *logrus.Logger
}

// NewKubernetesClientManager creates a new instance of KubernetesClientManager.
//
// Parameters:
// - log: A pointer to a logrus.Logger instance for logging purposes.
//
// Returns:
// - A pointer to a new instance of KubernetesClientManager.
func NewKubernetesClientManager(log *logrus.Logger) *KubernetesClientManager {
	return &KubernetesClientManager{log: log}
}

// GetKubernetesClient returns a Kubernetes clientset, creating it if it doesn't exist.
//
// This method ensures that the Kubernetes clientset is created only once using sync.Once.
// It attempts to create an in-cluster Kubernetes configuration and then uses it to create
// a clientset. If any error occurs during this process, it logs the error and returns it.
//
// Returns:
// - A pointer to a kubernetes.Clientset if successful.
// - An error if there was an issue creating the clientset or retrieving the configuration.
func (m *KubernetesClientManager) GetKubernetesClient() (*kubernetes.Clientset, error) {
	var err error
	m.once.Do(func() {
		config, errConfig := rest.InClusterConfig()
		if errConfig != nil {
			err = fmt.Errorf("failed to get in-cluster Kubernetes config: %w", errConfig)
			m.log.Error(err)
			return
		}

		m.clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			err = fmt.Errorf("unable to create client set for in-cluster Kubernetes config: %w", err)
			m.log.Error(err)
			return
		}

		m.log.Info("Successfully created Kubernetes clientset")
	})

	if err != nil {
		return nil, err
	}

	return m.clientset, nil
}
