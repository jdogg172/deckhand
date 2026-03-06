# Deckhand Architecture

## Overview

Deckhand is a terminal UI application for OpenShift and Tekton troubleshooting.

Primary design goals:
- portable
- cross-platform
- no installer
- keyboard-first
- operationally safe
- works with existing kubeconfig and cluster RBAC
- unifies workloads and pipelines

## Architectural Principles

### 1. Vertical slice first
The app should deliver a working thin slice quickly:
- load kubeconfig
- select context/namespace
- list pods
- view details
- switch modes

### 2. Clear boundaries
The codebase is split into:
- config
- clients
- resources
- actions
- watchers
- UI
- RBAC
- diagnostics

### 3. UI does not talk directly to API clients
The Bubble Tea model should depend on services/interfaces, not raw client-go calls everywhere.

### 4. Graceful degradation
The app must tolerate:
- missing Tekton CRDs
- OpenShift-specific APIs not present
- RBAC-denied resources or verbs
- partial visibility

### 5. Safe-by-default
Destructive actions:
- should be confirmable
- should be disable-able via read-only mode
- should surface meaningful error messages

## Future Growth

### OpenShift-specific
- Routes
- Projects
- DeploymentConfigs
- ImageStreams
- Operators / CSVs
- SCC awareness

### Pipeline-specific
- aggregated logs
- failure dashboard
- TaskRun -> pod -> events correlation
- workspaces / PVC relationships

