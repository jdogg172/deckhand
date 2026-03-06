# Local Test Quickstart

This guide prioritizes practical local validation for Deckhand.

## Prerequisites

### kind + Tekton lane
- `kind`
- `kubectl`
- `bash` (or Git Bash)

### OpenShift Local lane
- OpenShift Local (`crc`) installed and running
- `oc`
- OpenShift Pipelines operator installed in cluster

## Lane A: kind + Tekton (primary)

### Bash
```bash
./scripts/setup-kind-tekton.sh
```

### PowerShell
```powershell
./scripts/setup-kind-tekton.ps1
```

After setup:
```bash
kubectl config current-context
kubectl -n deckhand-lab get pods
kubectl -n deckhand-lab get pipelineruns
```

Run Deckhand:
```bash
./deckhand.exe --context kind-deckhand-dev --namespace deckhand-lab
```

What to validate:
- Ops mode pod list/details/yaml/events/logs
- Pending and CrashLoopBackOff visibility
- Pipeline mode PipelineRuns/TaskRuns/status/reason/duration
- failed/stuck highlighting
- TaskRun -> related pod logs drilldown
- namespace switching and refresh behavior
- read-only mode and RBAC action gating

## Lane B: OpenShift Local + OpenShift Pipelines (validation)

Start OpenShift Local and login:
```bash
crc start
oc login -u kubeadmin -p <password>
```

Install/verify OpenShift Pipelines as needed, then run Deckhand against CRC context:
```bash
./deckhand.exe --context crc-admin --namespace default
```

What to validate:
- project-aware namespace/project switching
- Route listing when route API exists
- Pipeline mode behavior with OpenShift Pipelines resources
- graceful fallback when Route or Tekton APIs are absent

## Reset kind lane

### Bash
```bash
./scripts/reset-kind-tekton.sh
```

### PowerShell
```powershell
./scripts/reset-kind-tekton.ps1
```
