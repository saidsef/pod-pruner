# Pod Pruner

[![CI](https://github.com/saidsef/pod-pruner/actions/workflows/ci.yml/badge.svg)](https://github.com/saidsef/pod-pruner/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/saidsef/pod-pruner)](https://goreportcard.com/report/github.com/saidsef/pod-pruner)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/saidsef/pod-pruner)
[![GoDoc](https://godoc.org/github.com/saidsef/pod-pruner?status.svg)](https://pkg.go.dev/github.com/saidsef/pod-pruner?tab=doc)
![GitHub release(latest by date)](https://img.shields.io/github/v/release/saidsef/pod-pruner)
![Commits](https://img.shields.io/github/commits-since/saidsef/pod-pruner/latest.svg)
![GitHub](https://img.shields.io/github/license/saidsef/pod-pruner)

Pod Pruner runs inside a Kubernetes cluster and deletes pods whose containers are stuck or have stopped, along with Jobs that have finished. It sweeps the namespaces you name, or every namespace, once at start-up and on a fixed interval after that. It starts in dry run and only logs what it would delete until you turn deletion on.

- **Pods** - deletes a pod when one of its containers is waiting or terminated with a reason you list, such as `CrashLoopBackOff` or `Error`
- **Jobs** - deletes a Job when one of its conditions has a type you list, `Complete` by default, and its pods go with it
- **Scope** - sweeps the namespaces in `NAMESPACES`, or the whole cluster when it is unset
- **Dry run** - logs every candidate and deletes nothing until `DRY_RUN` is `false`
- **Metrics** - counts deletions by namespace and reason on a Prometheus endpoint

## Quick start

```sh
git clone https://github.com/saidsef/pod-pruner.git
cd pod-pruner
kubectl create namespace pod-pruner
kubectl apply -k deployment/
kubectl logs -n pod-pruner deploy/pod-pruner
```

The Deployment in `deployment/base/` sets example values for the namespaces and statuses and leaves `DRY_RUN` unset, so the pruner logs its candidates and deletes nothing. Once the log lists only the pods you expect to lose, set `DRY_RUN` to `false` in that file and apply it again.

## Documentation

[Read the Docs](https://pod-pruner.readthedocs.io/en/latest/) hosts the same pages.

| Page | Contents |
|------|----------|
| [Installation](./docs/installation.md) | Requirements, Kubernetes manifests, the container image, building from source, log output |
| [Configuration](./docs/configuration.md) | Environment variables, namespace scope, container statuses, Job statuses |
| [Architecture](./docs/architecture.md) | What one sweep does, from listing to deletion |
| [Metrics](./docs/metrics.md) | Prometheus endpoint, metric names and labels, scrape setup |

## Requirements

Pod Pruner needs a Kubernetes cluster and permission to create a ClusterRole and a ClusterRoleBinding in it. Building it from source needs the Go version that `go.mod` sets.

## Alternatives

Pod Pruner takes its idea from [pod-reaper](https://github.com/saidsef/pod-reaper/tree/master), which is the one to try if this project does not fit.

## Source

Our latest and greatest source of *pod-pruner* can be found on [GitHub](https://github.com/saidsef/pod-pruner). [Fork us](https://github.com/saidsef/pod-pruner/fork)!

## Contributing

We would :heart: you to contribute by making a [pull request](https://github.com/saidsef/pod-pruner/pulls).

Please read the official [Contribution Guide](./CONTRIBUTING.md) for more information on how you can contribute.
