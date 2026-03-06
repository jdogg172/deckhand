# Test Matrix

## Quick setup commands
- `bash ./scripts/setup-kind-tekton.sh`
- `./deckhand.exe --context kind-deckhand-dev --namespace deckhand-lab`
- `./deckhand.exe --readonly --context kind-deckhand-dev --namespace deckhand-lab`

## Environment A: kind + Tekton
### Cluster sanity
- [ ] kubeconfig loads correctly
- [ ] current context detected
- [ ] namespace switch works
- [ ] refresh works
- [ ] no crash on empty namespace

### Ops mode
- [ ] list pods
- [ ] show pod details
- [ ] show yaml
- [ ] show recent events
- [ ] stream logs
- [ ] highlight Running/Pending/CrashLoopBackOff
- [ ] delete flow prompts correctly
- [ ] patch flow prompts correctly
- [ ] read-only mode blocks delete/patch actions
- [ ] RBAC-denied actions show disabled reasons

### Pipeline mode
- [ ] list PipelineRuns
- [ ] list TaskRuns for selected PipelineRun
- [ ] show status/reason/duration
- [ ] jump TaskRun -> pod
- [ ] fetch logs
- [ ] highlight failed pipeline
- [ ] highlight running pipeline
- [ ] cancel PipelineRun action gated correctly

### Failure fixtures
- [ ] failed PipelineRun visible
- [ ] crashlooping pod visible
- [ ] pending pod visible

## Environment B: OpenShift Local + OpenShift Pipelines
### OpenShift specifics
- [ ] detect OpenShift/project semantics
- [ ] project switch works
- [ ] Routes list works
- [ ] no panic if Route API unavailable
- [ ] OpenShift Pipelines resources render if installed

### UX truth checks
- [ ] terminal workflow feels correct for OpenShift users
- [ ] errors are readable
- [ ] RBAC-denied actions are visibly disabled
- [ ] read-only mode is obvious and enforced
