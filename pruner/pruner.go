package main

import (
	"os"
	"strings"
	"time"

	"github.com/saidsef/pod-pruner/pruner/internal/auth"
	"github.com/saidsef/pod-pruner/pruner/internal/resources"
	"github.com/saidsef/pod-pruner/pruner/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: false})

	dryRun := utils.GetEnv("DRY_RUN", "true", log)
	NAMESPACES := strings.Split(os.Getenv("NAMESPACES"), ",")
	RESOURCES := strings.Split(utils.GetEnv("RESOURCES", "PODS", log), ",")

	k8sManager := auth.NewKubernetesClientManager(log)
	clientset, err := k8sManager.GetKubernetesClient()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Infof("Resources to include in pruner are %v", RESOURCES)
		for _, namespace := range NAMESPACES {
			if utils.Contains(RESOURCES, "PODS") {
				containers, err := resources.GetContainers(clientset, namespace)
				if err != nil {
					log.Errorf("Error fetching containers in namespace '%s': %v", namespace, err)
					continue
				}

				if len(containers) > 0 {
					if dryRun == "true" {
						log.Infof("Dry run enabled. The following containers would be deleted: %v", containers)
					} else {
						log.Infof("Containers to be pruned in namespace '%s': %v", namespace, containers)
						resources.DeleteContainers(clientset, namespace, containers, log)
					}
				} else {
					log.Infof("No containers to prune in namespace '%s'", namespace)
				}

			}
			if utils.Contains(RESOURCES, "JOBS") {
				jobs, err := resources.GetJobs(clientset, namespace, log)
				if err != nil {
					log.Errorf("Error fetching jobs in namespace '%s': %v", namespace, err)
					continue
				}

				if len(jobs) > 0 {
					if dryRun == "true" {
						log.Infof("Dry run enabled. The following jobd would be deleted: %v", jobs)
					} else {
						log.Infof("Jobs to be pruned in namespace '%s': %v", namespace, jobs)
						resources.DeleteJobs(clientset, namespace, jobs, log)
					}
				} else {
					log.Infof("No jobs to prune in namespace '%s': %v", namespace, jobs)
				}
			}
		}
	}
}
