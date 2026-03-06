# Deckhand

Deckhand is a keyboard-first TUI for Kubernetes/OpenShift operations and Tekton pipeline troubleshooting.

## Current capabilities
- Ops mode: pods list, details, events, yaml, logs
- Pipeline mode: PipelineRuns, TaskRuns, related pod drilldown, PipelineRun cancel
- Context switching and namespace/project switching
- OpenShift route listing when Route API is present
- Graceful degradation when Tekton/OpenShift APIs are unavailable
- Read-only and RBAC-aware mutating action gating with confirmations

## Build and run
```bash
go mod tidy
go test ./...
go build -o dist/deckhand ./cmd/deckhand

# Linux/macOS
./dist/deckhand
./dist/deckhand --readonly
./dist/deckhand --mode pipeline

# Windows (PowerShell)
.\dist\deckhand.exe
.\dist\deckhand.exe --readonly
.\dist\deckhand.exe --mode pipeline
```

## Local test strategy
Deckhand uses a two-lane local test ladder:

1. **Primary dev lane:** kind + Tekton Pipelines
2. **OpenShift validation lane:** OpenShift Local + OpenShift Pipelines

See:
- [docs/test-strategy.md](docs/test-strategy.md)
- [docs/test-matrix.md](docs/test-matrix.md)
- [docs/openshift-local-notes.md](docs/openshift-local-notes.md)
- [docs/local-test-quickstart.md](docs/local-test-quickstart.md)

## Test environment assets
- [scripts/setup-kind-tekton.sh](scripts/setup-kind-tekton.sh)
- [scripts/setup-kind-tekton.ps1](scripts/setup-kind-tekton.ps1)
- [scripts/reset-kind-tekton.sh](scripts/reset-kind-tekton.sh)
- [scripts/reset-kind-tekton.ps1](scripts/reset-kind-tekton.ps1)
- [manifests/tekton-task-hello.yaml](manifests/tekton-task-hello.yaml)
- [manifests/tekton-pipeline-hello.yaml](manifests/tekton-pipeline-hello.yaml)
- [manifests/tekton-pipelinerun-success.yaml](manifests/tekton-pipelinerun-success.yaml)
- [manifests/tekton-pipelinerun-fail.yaml](manifests/tekton-pipelinerun-fail.yaml)
- [manifests/sample-crashloop-pod.yaml](manifests/sample-crashloop-pod.yaml)
- [manifests/sample-pending-pod.yaml](manifests/sample-pending-pod.yaml)

## Immediate kind + Tekton smoke path
```bash
./scripts/setup-kind-tekton.sh
./dist/deckhand --context kind-deckhand-dev --namespace deckhand-lab

# Windows (PowerShell)
.\scripts\setup-kind-tekton.ps1
.\dist\deckhand.exe --context kind-deckhand-dev --namespace deckhand-lab
```

Then verify in Deckhand:
- Ops mode with pending/crashloop highlighting
- Pipeline mode with success/fail PipelineRuns and TaskRuns
- logs/yaml/events/details panes
- namespace switching and read-only/RBAC behavior
