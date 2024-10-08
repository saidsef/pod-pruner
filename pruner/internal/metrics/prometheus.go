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

package metrics

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saidsef/pod-pruner/pruner/utils"
	"github.com/sirupsen/logrus"
)

// Define counters for metrics
var (
	// PodsPruned counts the total number of pods pruned, labelled by namespace.
	PodsPruned = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pods_pruned_total",
			Help: "Total number of pods pruned",
		},
		[]string{"namespace", "state"},
	)

	// ContainersPruned counts the total number of containers pruned, labelled by namespace.
	ContainersPruned = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "containers_pruned_total",
			Help: "Total number of containers pruned",
		},
		[]string{"namespace", "state"},
	)

	// JobsPruned counts the total number of jobs pruned, labelled by namespace.
	JobsPruned = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jobs_pruned_total",
			Help: "Total number of jobs pruned",
		},
		[]string{"namespace", "state"},
	)

	once sync.Once
)

// init registers the defined metrics with Prometheus.
func init() {
	once.Do(func() {
		logger := utils.Logger()
		utils.LogWithFields(logrus.InfoLevel, []string{}, "registering prometheus metrics count vectors")
		prometheus.MustRegister(PodsPruned, ContainersPruned, JobsPruned)
		StartMetricsServer(logger)
	})
}

// StartMetricsServer starts the metrics server and adds a handler for the /metrics endpoint.
func StartMetricsServer(log *logrus.Logger) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		port := utils.GetEnv("PORT", "8080", log)

		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
			utils.LogWithFields(logrus.FatalLevel, []string{}, "Metrics server failed to start", err)
		}
	}()
}
