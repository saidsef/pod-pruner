# Pod Pruner: Kubernetes Container Pruner

[![Go Report Card](https://goreportcard.com/badge/github.com/saidsef/pod-pruner)](https://goreportcard.com/report/github.com/saidsef/pod-pruner)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/saidsef/pod-pruner)
[![GoDoc](https://godoc.org/github.com/saidsef/pod-pruner?status.svg)](https://pkg.go.dev/github.com/saidsef/pod-pruner?tab=doc)
![GitHub release(latest by date)](https://img.shields.io/github/v/release/saidsef/pod-pruner)
![Commits](https://img.shields.io/github/commits-since/saidsef/pod-pruner/latest.svg)
![GitHub](https://img.shields.io/github/license/saidsef/pod-pruner)

This is a Kubernetes application written in Go (Golang) that periodically prunes containers in specified namespaces based on their statuses. The application can operate in a dry-run mode, allowing you to see which containers would be deleted without actually removing them.

## Alternatives

This application was inspired by [pod-reaper](https://github.com/saidsef/pod-reaper/tree/master). If you need an alternative, I suggest using [pod-reaper](https://github.com/saidsef/pod-reaper/tree/master).

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [How It Works](#how-it-works)

## Prerequisites

- Go (version 1.22 or later)
- Kubernetes cluster
- Access to the Kubernetes API from within the cluster

## Installation

1. Clone the repository:
```bash
git clone https://github.com/saidsef/pod-pruner.git
```

2. Build the application:
```bash
go build pruner.go
```

3. Ensure that the application is packaged into a Docker image and pushed to a container registry if you plan to deploy it in a Kubernetes environment.

## Configuration

The application requires certain environment variables to be set:

- `DRY_RUN`: Set to `"true"` to enable dry-run mode (default is `"true"`).
- `RESOURCES`: A comma-separated list of Kubernetes resources (default is `"PODS"`)
- `NAMESPACES`: A comma-separated list of namespaces to monitor for containers to prune.
- `CONTAINER_STATUSES`: A comma-separated list of container statuses to filter by (e.g., `Error,ContainerStatusUnknown,Unknown,Completed`).
- `JOB_STATUSES`: A comma-separated list of jobs statuses to filter by (default is `Complete`).

Example of setting environment variables in a Kubernetes deployment spec:

```bash
kubectl apply -k deployment/ -n pod-pruner
```

## Usage

Once the application is deployed, it will start monitoring the specified namespaces every `60 seconds`. It will log the containers that are eligible for pruning based on their statuses. If dry-run mode is disabled, it will proceed to delete the identified containers.

## How It Works

1. **Environment Variables**: The application retrieves configuration values from environment variables.
2. **Kubernetes Client**: It creates a Kubernetes client using in-cluster configuration to interact with the Kubernetes API.
3. **Container Monitoring**: Every 60 seconds, it checks the specified namespaces for containers that are in the defined states (e.g., `Waiting`, `Terminated`).
4. **Pruning Logic**: If containers are found, it either logs the containers that would be deleted (in dry-run mode) or deletes them from the cluster.

## Source

Our latest and greatest source of *Reverse Geocoding* can be found on [GitHub]. [Fork us](https://github.com/saidsef/pod-pruner/fork)!

## Contributing

We would :heart: you to contribute by making a [pull request](https://github.com/saidsef/pod-pruner/pulls).

Please read the official [Contribution Guide](./CONTRIBUTING.md) for more information on how you can contribute.
