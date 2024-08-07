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

	clientset, err := auth.CreateKubernetesClient(log)
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
				log.Infof("Containers to be pruned in namespace '%s': %v", namespace, containers)

				if len(containers) > 0 {
					if dryRun == "true" {
						log.Infof("Dry run enabled. The following containers would be deleted: %v", containers)
					} else {
						resources.DeleteContainers(clientset, namespace, containers, log)
					}
				} else {
					log.Infof("No containers to prune in namespace '%s'", namespace)
				}

			}
			if utils.Contains(RESOURCES, "JOBS") {
				jobs, err := resources.GetJobs(clientset, namespace)
				if err != nil {
					log.Errorf("Error fetching jobs in namespace '%s': %v", namespace, err)
					continue
				}
				log.Infof("Jobs to be pruned in namespace '%s': %v", namespace, jobs)

				if len(jobs) > 0 {
					if dryRun == "true" {
						log.Infof("Dry run enabled. The following jobd would be deleted: %v", jobs)
					} else {
						resources.DeleteJobs(clientset, namespace, jobs, log)
					}
				} else {
					log.Infof("No jobs to prune in namespace '%s'", namespace)
				}
			}
		}
	}
}
