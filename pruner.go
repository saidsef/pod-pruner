package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	dryRun := getEnv("DRY_RUN", "true", log)
	namespaces := strings.Split(os.Getenv("NAMESPACES"), ",")

	clientset, err := createKubernetesClient(log)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, namespace := range namespaces {
			containers, err := getContainers(clientset, namespace)
			if err != nil {
				log.Errorf("Error fetching containers in namespace '%s': %v", namespace, err)
				continue
			}
			log.Infof("Containers to be pruned in namespace '%s': %v", namespace, containers)

			if len(containers) > 0 {
				if dryRun == "true" {
					log.Infof("Dry run enabled. The following containers would be deleted: %v", containers)
				} else {
					deleteContainers(clientset, namespace, containers, log)
				}
			} else {
				log.Infof("No containers to prune in namespace '%s'", namespace)
			}
		}
	}
}

// getEnv retrieves the value of the environment variable specified by key.
// If the variable is not set, it returns the defaultValue and logs a warning.
//
// Parameters:
// - key: The name of the environment variable to retrieve.
// - defaultValue: The value to return if the environment variable is not set.
// - log: A logger instance for logging warnings.
//
// Returns:
// - The value of the environment variable or the default value if not set.
func getEnv(key, defaultValue string, log *logrus.Logger) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Warnf("%s environment variable not set, defaulting to %s", key, defaultValue)
		return defaultValue
	}
	return value
}

// createKubernetesClient creates a new Kubernetes client using in-cluster configuration.
// It returns a pointer to the Kubernetes clientset or an error if the creation fails.
//
// Parameters:
// - log: A logger instance for logging errors.
//
// Returns:
// - A pointer to the Kubernetes clientset.
// - An error if the clientset could not be created.
func createKubernetesClient(log *logrus.Logger) (*kubernetes.Clientset, error) {
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

// getContainers retrieves a list of container names in a given namespace that are in specified states.
// It checks the environment variable CONTAINER_STATUSES to determine which states to filter by.
//
// Parameters:
// - clientset: A pointer to the Kubernetes clientset.
// - namespace: The namespace from which to retrieve the containers.
//
// Returns:
// - A slice of strings containing the names of the containers in the specified states.
// - An error if the retrieval of pods fails.
func getContainers(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	statuses := strings.Split(os.Getenv("CONTAINER_STATUSES"), ",")
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var containers []string
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if isContainerInState(containerStatus, statuses) {
				containers = append(containers, fmt.Sprintf("%s/%s: %s", pod.Namespace, pod.Name, containerStatus.Name))
			}
		}
	}
	return containers, nil
}

// isContainerInState checks if a container's status matches any of the specified states.
//
// Parameters:
// - containerStatus: The status of the container to check.
// - statuses: A slice of strings representing the states to check against.
//
// Returns:
// - A boolean indicating whether the container is in one of the specified states.
func isContainerInState(containerStatus metav1.ContainerStatus, statuses []string) bool {
	if containerStatus.State.Waiting != nil && contains(statuses, containerStatus.State.Waiting.Reason) {
		return true
	}
	if containerStatus.State.Terminated != nil && contains(statuses, containerStatus.State.Terminated.Reason) {
		return true
	}
	return false
}

// contains checks if a string is present in a slice of strings.
//
// Parameters:
// - list: A slice of strings to search through.
// - str: The string to search for.
//
// Returns:
// - A boolean indicating whether the string is found in the slice.
func contains(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

// deleteContainers deletes the specified containers in a given namespace.
// It logs the success or failure of each deletion attempt.
//
// Parameters:
// - clientset: A pointer to the Kubernetes clientset.
// - namespace: The namespace from which to delete the containers.
// - containers: A slice of strings representing the containers to delete.
// - log: A logger instance for logging the results of the deletion attempts.
func deleteContainers(clientset *kubernetes.Clientset, namespace string, containers []string, log *logrus.Logger) {
	for _, container := range containers {
		parts := strings.Split(container, ": ")
		if len(parts) != 2 {
			log.Warnf("Unexpected format for container state: '%s'", container)
			continue
		}
		podInfo := parts[0]
		containerName := strings.Split(podInfo, "/")[1]

		err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), containerName, metav1.DeleteOptions{})
		if err != nil {
			log.Errorf("Failed to delete pod '%s' in namespace '%s': %v", podInfo, namespace, err)
		} else {
			log.Infof("Successfully deleted pod '%s' in namespace '%s'", podInfo, namespace)
		}
	}
}
