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
	"strings"
	"time"

	"github.com/saidsef/pod-pruner/pruner/internal/auth"
	"github.com/saidsef/pod-pruner/pruner/internal/metrics"
	"github.com/saidsef/pod-pruner/pruner/internal/resources"
	"github.com/saidsef/pod-pruner/pruner/utils"
	"github.com/sirupsen/logrus"
)

// main is the entry point of the application. It sets up logging,
// retrieves environment variables,
// and initiates a Kubernetes client manager to prune specified
// resources (containers and jobs)
// in the defined namespaces at regular intervals.
func main() {
	log := utils.Logger()
	// Retrieve the dry run mode from environment variables, defaulting to "true".
	dryRun := utils.GetEnv("DRY_RUN", "true", log)
	// Split the NAMESPACES environment variable into a slice.
	NAMESPACES := strings.Split(os.Getenv("NAMESPACES"), ",")
	// Split the RESOURCES environment variable into a slice, defaulting to "PODS".
	RESOURCES := strings.Split(utils.GetEnv("RESOURCES", "PODS", log), ",")

	// Create a new Kubernetes client manager.
	k8sManager := auth.NewKubernetesClientManager(log)
	clientset, err := k8sManager.GetKubernetesClient()
	if err != nil {
		utils.LogWithFields(logrus.FatalLevel, []string{}, "Kubenetes config error", err)
	}

	// Set up a ticker to trigger every 120 seconds.
	ticker := time.NewTicker(120 * time.Second)
	defer ticker.Stop()

	metrics.StartMetricsServer(log)
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
					utils.LogWithFields(logrus.ErrorLevel, append([]string{}, fmt.Sprintf("namespace:%s", namespace)), "Error fetching containers", err)
					continue
				}

				// If there are containers to prune, log the action based on dry run mode.
				if len(containers) > 0 {
					if dryRun == "true" {
						utils.LogWithFields(logrus.InfoLevel, append(containers, fmt.Sprintf("namespace:%s", namespace)), "Dry run mode. The following containers would be deleted")
					} else {
						utils.LogWithFields(logrus.InfoLevel, append(containers, namespace), "Containers to be pruned")
						resources.DeleteContainers(clientset, namespace, containers, log)
						metrics.ContainersPruned.WithLabelValues(namespace).Add(float64(len(containers))) // Increment the counter
					}
				} else {
					utils.LogWithFields(logrus.InfoLevel, append([]string{}, fmt.Sprintf("namespace:%s", namespace)), "No containers to prune")
				}
			}

			// Check if "JOBS" is included in the resources to prune.
			if utils.Contains(RESOURCES, "JOBS") {
				// Fetch jobs in the current namespace.
				jobs, err := resources.GetJobs(clientset, namespace, log)
				if err != nil {
					utils.LogWithFields(logrus.ErrorLevel, append([]string{}, fmt.Sprintf("namespace:%s", namespace)), "Error fetching jobs", err)
					continue
				}

				// If there are jobs to prune, log the action based on dry run mode.
				if len(jobs) > 0 {
					if dryRun == "true" {
						utils.LogWithFields(logrus.InfoLevel, jobs, "Dry run enabled. The following jobs would be deleted")
					} else {
						utils.LogWithFields(logrus.InfoLevel, append(jobs, fmt.Sprintf("namespace:%s", namespace)), "Jobs to be pruned")
						resources.DeleteJobs(clientset, namespace, jobs, log)
						metrics.JobsPruned.WithLabelValues(namespace).Add(float64(len(jobs))) // Increment the counter
					}
				} else {
					utils.LogWithFields(logrus.InfoLevel, append([]string{}, fmt.Sprintf("namespace:%s", namespace)), "No jobs to prune")
				}
			}
		}
	}
}
