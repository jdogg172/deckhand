# Test Strategy for Deckhand

## Objective
Create a local test ladder that balances:
- fast iteration
- low cost
- realistic OpenShift validation
- repeatable failure cases

## Lane 1: kind + Tekton
Use for:
- resource list rendering
- pod status highlighting
- details/yaml/events/logs panes
- namespace switching
- delete/patch safety flows
- PipelineRun / TaskRun rendering
- failure and stuck pipeline states

Execution path:
1. Bring up `kind` + Tekton
2. Apply sample fixtures (success/fail PipelineRuns, crashloop pod, pending pod)
3. Run Deckhand against `kind-deckhand-dev` / `deckhand-lab`
4. Verify Ops Mode and Pipeline Mode matrix before OpenShift lane

## Lane 2: OpenShift Local + OpenShift Pipelines
Use for:
- project-aware namespace handling
- Route support
- OpenShift Pipelines integration assumptions
- missing/present API handling
- OpenShift user workflow expectations

Execution path:
1. Start OpenShift Local and login
2. Ensure OpenShift Pipelines is installed and healthy
3. Run Deckhand against CRC context
4. Validate project semantics, Route support, and graceful API behavior

## Design Rules for Deckhand
1. The app must work well against plain Kubernetes + Tekton.
2. The app must detect and gracefully enhance behavior when OpenShift APIs exist.
3. The app must not assume Tekton CRDs are present.
4. The app must not assume OpenShift-specific resources are present.
5. All resource/action handling should degrade cleanly under RBAC denial.
