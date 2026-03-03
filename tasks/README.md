# Task Tracking Index

This file tracks all tasks for improving the pod-pruner codebase. When a task is completed, update its status in both the task file and this index.

## Quick Stats

| Status | Count |
|--------|-------|
| pending | 17 |
| in_progress | 0 |
| completed | 3 |
| total | 21 |

## Tasks by ID

| ID | Task | Type | Risk | Status | File |
|----|------|------|------|--------|------|
| 00 | Fix Goroutine Closure Variable Capture Bug in Jobs | bug | Critical | completed | [00-bug-goroutine-closure-jobs.md](./00-bug-goroutine-closure-jobs.md) |
| 01 | Fix Nil Pointer Dereference Risk in Container Status | bug | Critical | completed | [01-bug-nil-pointer-containers.md](./01-bug-nil-pointer-containers.md) |
| 02 | Fix Invalid Go Version (1.25 doesn't exist) | bug | Critical | completed | [02-bug-go-version-invalid.md](./02-bug-go-version-invalid.md) |
| 03 | Fix Security Workflow That Never Runs | bug | High | pending | [03-bug-security-workflow-never-runs.md](./03-bug-security-workflow-never-runs.md) |
| 04 | Investigate and Fix Suspicious YAML Import Paths | bug | High | pending | [04-bug-suspicious-yaml-imports.md](./04-bug-suspicious-yaml-imports.md) |
| 05 | Add Context Timeouts for Job Operations | improvement | High | pending | [05-improvement-unbounded-context-jobs.md](./05-improvement-unbounded-context-jobs.md) |
| 06 | Add Pagination Support for Job Listing | improvement | High | pending | [06-improvement-pagination-jobs.md](./06-improvement-pagination-jobs.md) |
| 07 | Remove or Secure Auto-Approve CI Job | security | High | pending | [07-security-auto-approve-risk.md](./07-security-auto-approve-risk.md) |
| 08 | Disable Service Account Token Automounting | security | Medium | pending | [08-security-serviceaccount-automount.md](./08-security-serviceaccount-automount.md) |
| 09 | Fix Log Field Format Inconsistency | bug | Medium | pending | [09-bug-log-field-format.md](./09-bug-log-field-format.md) |
| 10 | Add Error Handling for Delete Operations | improvement | Medium | pending | [10-improvement-delete-error-handling.md](./10-improvement-delete-error-handling.md) |
| 11 | Add Graceful Shutdown for Metrics Server | improvement | Low | pending | [11-improvement-metrics-shutdown.md](./11-improvement-metrics-shutdown.md) |
| 12 | Fix Namespace Trimming Inconsistency | improvement | Low | pending | [12-improvement-namespace-trimming.md](./12-improvement-namespace-trimming.md) |
| 13 | Add Validation for Empty Namespaces | improvement | Low | pending | [13-improvement-empty-namespace-validation.md](./13-improvement-empty-namespace-validation.md) |
| 14 | Make Tick Interval Configurable | improvement | Low | pending | [14-improvement-configurable-tick-interval.md](./14-improvement-configurable-tick-interval.md) |
| 15 | Remove Unnecessary sync.Once in init() | bug | Low | pending | [15-bug-sync-once-unnecessary.md](./15-bug-sync-once-unnecessary.md) |
| 16 | Fix Unused PodsPruned Metric | bug | Medium | pending | [16-bug-unused-pod-metric.md](./16-bug-unused-pod-metric.md) |
| 17 | Explicitly Set DRY_RUN in Deployment | improvement | Medium | pending | [17-improvement-deployment-dry-run-default.md](./17-improvement-deployment-dry-run-default.md) |
| 18 | Fix Deprecated GitHub Actions Output Syntax | bug | Medium | pending | [18-bug-deprecated-output-syntax.md](./18-bug-deprecated-output-syntax.md) |
| 19 | Add Unit Tests for All Packages | improvement | Critical | pending | [19-improvement-add-tests.md](./19-improvement-add-tests.md) |
| 20 | Update Outdated Helm Version in CI | improvement | Medium | pending | [20-improvement-helm-version.md](./20-improvement-helm-version.md) |

## Tasks by Risk Level

### Critical (4 tasks)
- [00](./00-bug-goroutine-closure-jobs.md) - Fix Goroutine Closure Variable Capture Bug in Jobs
- [01](./01-bug-nil-pointer-containers.md) - Fix Nil Pointer Dereference Risk in Container Status
- [02](./02-bug-go-version-invalid.md) - Fix Invalid Go Version (1.25 doesn't exist)
- [19](./19-improvement-add-tests.md) - Add Unit Tests for All Packages

### High (5 tasks)
- [03](./03-bug-security-workflow-never-runs.md) - Fix Security Workflow That Never Runs
- [04](./04-bug-suspicious-yaml-imports.md) - Investigate and Fix Suspicious YAML Import Paths
- [05](./05-improvement-unbounded-context-jobs.md) - Add Context Timeouts for Job Operations
- [06](./06-improvement-pagination-jobs.md) - Add Pagination Support for Job Listing
- [07](./07-security-auto-approve-risk.md) - Remove or Secure Auto-Approve CI Job

### Medium (6 tasks)
- [08](./08-security-serviceaccount-automount.md) - Disable Service Account Token Automounting
- [09](./09-bug-log-field-format.md) - Fix Log Field Format Inconsistency
- [10](./10-improvement-delete-error-handling.md) - Add Error Handling for Delete Operations
- [16](./16-bug-unused-pod-metric.md) - Fix Unused PodsPruned Metric
- [17](./17-improvement-deployment-dry-run-default.md) - Explicitly Set DRY_RUN in Deployment
- [18](./18-bug-deprecated-output-syntax.md) - Fix Deprecated GitHub Actions Output Syntax
- [20](./20-improvement-helm-version.md) - Update Outdated Helm Version in CI

### Low (6 tasks)
- [11](./11-improvement-metrics-shutdown.md) - Add Graceful Shutdown for Metrics Server
- [12](./12-improvement-namespace-trimming.md) - Fix Namespace Trimming Inconsistency
- [13](./13-improvement-empty-namespace-validation.md) - Add Validation for Empty Namespaces
- [14](./14-improvement-configurable-tick-interval.md) - Make Tick Interval Configurable
- [15](./15-bug-sync-once-unnecessary.md) - Remove Unnecessary sync.Once in init()

## Task Completion Instructions

When completing a task:

1. **Edit the task file**: Update the status field at the top of the task file from `pending` to `completed`
2. **Update this index**: Update the status column in the appropriate table(s) above
3. **Update stats**: Update the count in the Quick Stats section
4. **Commit changes**: Include the task file update in your commit

### Example Task File Update

```markdown
## Status
completed  # Changed from pending

## Type
bug

## Risk
Critical

## Complexity
Medium

## Description
Fix the classic Go concurrency bug...
```

### Example Index Update

Change:
```markdown
| 00 | Fix Goroutine Closure Variable Capture Bug in Jobs | bug | Critical | completed | [00-bug-goroutine-closure-jobs.md](./00-bug-goroutine-closure-jobs.md) |
```

To:
```markdown
| 00 | Fix Goroutine Closure Variable Capture Bug in Jobs | bug | Critical | completed | [00-bug-goroutine-closure-jobs.md](./00-bug-goroutine-closure-jobs.md) |
```

### Example Stats Update

Change:
```markdown
| Status | Count |
|--------|-------|
| pending | 19 |
| completed | 1 |
| total | 21 |
```

To:
```markdown
| Status | Count |
|--------|-------|
| pending | 19 |
| completed | 1 |
| total | 21 |
```
