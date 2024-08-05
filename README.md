# Pod Pruner: Kubernetes Container Pruner

This is a Kubernetes application written in Go (Golang) that periodically prunes containers in specified namespaces based on their statuses. The application can operate in a dry-run mode, allowing you to see which containers would be deleted without actually removing them.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [How It Works](#how-it-works)
- [License](#license)

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
- `NAMESPACES`: A comma-separated list of namespaces to monitor for containers to prune.
- `CONTAINER_STATUSES`: A comma-separated list of container statuses to filter by (e.g., `Waiting,Terminated`).

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
