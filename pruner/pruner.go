/*
Copyright 2024 Said Sef

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/saidsef/pod-pruner/pruner/internal/auth"
	_ "github.com/saidsef/pod-pruner/pruner/internal/metrics"
	"github.com/saidsef/pod-pruner/pruner/internal/resources"
	"github.com/saidsef/pod-pruner/pruner/utils"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

// main is the entry point of the application. It sets up logging,
// retrieves environment variables, and initiates a Kubernetes client
// manager to prune specified resources (containers and jobs) in the
// defined namespaces at regular intervals.
func main() {
	log := utils.Logger()
	// Deletion is destructive, so an unset or unparseable DRY_RUN stays in dry run.
	dryRun := utils.GetEnvBool("DRY_RUN", true, log)
	// An empty namespace means every namespace to client-go, so refuse to start
	// rather than prune the whole cluster on a missing variable.
	NAMESPACES := utils.SplitAndTrim(os.Getenv("NAMESPACES"))
	if len(NAMESPACES) == 0 {
		utils.LogWithFields(logrus.FatalLevel, []string{}, "NAMESPACES environment variable is not set or empty, refusing to run against every namespace")
	}
	RESOURCES := utils.SplitAndTrim(utils.GetEnv("RESOURCES", "PODS", log))

	// Create a new Kubernetes client manager.
	k8sManager := auth.NewKubernetesClientManager(log)
	clientset, err := k8sManager.GetKubernetesClient()
	if err != nil {
		utils.LogWithFields(logrus.FatalLevel, []string{}, "Kubernetes config error", err)
	}

	// Set up a ticker to trigger every 120 seconds.
	ticker := time.NewTicker(120 * time.Second)
	defer ticker.Stop()

	utils.LogWithFields(logrus.InfoLevel, RESOURCES, "Resources to include in pruner")

	// Main loop that runs every tick.
	for range ticker.C {
		// Iterate over each namespace defined in the environment variable.
		for _, namespace := range NAMESPACES {
			// Check if "PODS" is included in the resources to prune.
			if utils.Contains(RESOURCES, "PODS") {
				// Fetch containers in the current namespace.
				containers, err := resources.GetContainers(clientset, namespace)
				if err != nil {
					utils.LogWithFields(
						logrus.ErrorLevel,
						[]string{fmt.Sprintf("namespace:%s", namespace)},
						"Error fetching containers",
						err,
					)
					continue
				}

				// Handle pruning logic for containers.
				handlePruning("containers", containers, dryRun, log, clientset)
			}

			// Check if "JOBS" is included in the resources to prune.
			if utils.Contains(RESOURCES, "JOBS") {
				// Fetch jobs in the current namespace.
				jobs, err := resources.GetJobs(clientset, namespace, log)
				if err != nil {
					utils.LogWithFields(
						logrus.ErrorLevel,
						[]string{fmt.Sprintf("namespace:%s", namespace)},
						"Error fetching jobs",
						err,
					)
					continue
				}

				// Handle pruning logic for jobs.
				handlePruning("jobs", jobs, dryRun, log, clientset)
			}
		}
	}
}

// handlePruning handles the common logic for pruning resources.
// It logs the actions taken based on the dry run mode and performs
// the deletion of specified resources if not in dry run mode.
//
// Parameters:
// - resourceType: A string indicating the type of resource being pruned (e.g., "containers" or "jobs").
// - items: A slice of ContainerInfo representing the resource identifiers to be pruned.
// - dryRun: Whether to log the resources instead of deleting them.
// - log: A pointer to a logrus.Logger instance for logging purposes.
// - clientset: A pointer to a Kubernetes Clientset for interacting with the Kubernetes API.
func handlePruning(resourceType string, items []resources.ContainerInfo, dryRun bool, log *logrus.Logger, clientset *kubernetes.Clientset) {
	var values []string
	for _, item := range items {
		values = append(values, item.Namespace, item.PodName, item.Status)
	}
	if len(items) > 0 {
		if dryRun {
			utils.LogWithFields(
				logrus.InfoLevel,
				values,
				fmt.Sprintf("Dry run mode. The following %s would be deleted", resourceType),
			)
		} else {
			utils.LogWithFields(logrus.InfoLevel,
				values,
				fmt.Sprintf("%s to be pruned", resourceType))
			if resourceType == "containers" {
				resources.DeleteContainers(clientset, items, log)
			} else if resourceType == "jobs" {
				resources.DeleteJobs(clientset, items, log)
			}
		}

	} else {
		utils.LogWithFields(
			logrus.InfoLevel,
			values,
			fmt.Sprintf("No %s to prune", resourceType),
		)
	}
}
