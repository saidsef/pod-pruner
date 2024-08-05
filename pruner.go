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
	// Set up logging
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Create a Kubernetes client
	var clientset *kubernetes.Clientset
	var err error

	// Use InClusterConfig if running inside a cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.WithError(err).Panic("error getting in cluster kubernetes config")
		logrus.Panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.WithError(err).Panic("unable to get client set for in cluster kubernetes config")
	}
	if clientSet == nil {
		message := "kubernetes client set cannot be nil"
		logrus.Panic(message)
	}

	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Get the comma-separated list of namespaces from the environment variable
	namespacesEnv := os.Getenv("NAMESPACES")
	namespaces := strings.Split(namespacesEnv, ",")

	// Periodically query the API
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, namespace := range namespaces {
				containersState, err := getcontainersState(clientset, namespace)
				if err != nil {
					log.Errorf("Error fetching containers in namespace %s: %v", namespace, err)
					continue
				}
				log.Infof("Containers in error state in namespace %s: %v", namespace, containersState)
			}
		}
	}
}

func getcontainersState(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	stateEnv := os.Getenv("CONTAINER_STATUSES")
	statuses := strings.Split(stateEnv, ",")
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var containersState []string
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Waiting != nil && (contains(statuses, containerStatus.State.Waiting.Reason)) {
				containersState = append(containersState, fmt.Sprintf("%s/%s: %s", pod.Namespace, pod.Name, containerStatus.Name))
			} else if containerStatus.State.Terminated != nil && (contains(statuses, containerStatus.State.Terminated.Reason)) {
				containersState = append(containersState, fmt.Sprintf("%s/%s: %s", pod.Namespace, pod.Name, containerStatus.Name))
			}
		}
	}

	return containersState, nil
}

// Function to check if a string is in a slice
func contains(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}
