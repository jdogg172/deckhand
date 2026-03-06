#!/usr/bin/env bash
set -euo pipefail
































































kubectl get pipelineruns -n $Namespacekubectl get pods -Akubectl config current-contextWrite-Host "[6/6] Done"kubectl apply -n $Namespace -f (Join-Path $rootDir "manifests/sample-pending-pod.yaml")kubectl apply -n $Namespace -f (Join-Path $rootDir "manifests/sample-crashloop-pod.yaml")kubectl apply -n $Namespace -f (Join-Path $rootDir "manifests/tekton-pipelinerun-fail.yaml")kubectl apply -n $Namespace -f (Join-Path $rootDir "manifests/tekton-pipelinerun-success.yaml")kubectl apply -n $Namespace -f (Join-Path $rootDir "manifests/tekton-pipeline-hello.yaml")kubectl apply -n $Namespace -f (Join-Path $rootDir "manifests/tekton-task-hello.yaml")Write-Host "[5/6] Applying sample manifests"kubectl -n tekton-pipelines rollout status deploy/tekton-pipelines-webhook --timeout=$timeoutkubectl -n tekton-pipelines rollout status deploy/tekton-pipelines-controller --timeout=$timeout$timeout = "$TimeoutSeconds" + "s"Write-Host "[4/6] Waiting for Tekton controller rollout"kubectl apply -f $TektonUrlWrite-Host "[3/6] Installing Tekton Pipelines"kubectl create namespace $Namespace 2>$null | Out-NullWrite-Host "[2/6] Creating app namespace"}    kind create cluster --name $ClusterName --config $kindConfigPath} else {    Write-Host "Cluster $ClusterName already exists; reusing"if ($existingClusters -contains $ClusterName) {$existingClusters = kind get clustersWrite-Host "[1/6] Creating or reusing kind cluster: $ClusterName"Set-Content -Path $kindConfigPath -Value $kindConfig -Encoding utf8"@  - role: worker  - role: worker  - role: control-planenodes:name: $ClusterNameapiVersion: kind.x-k8s.io/v1alpha4kind: Cluster$kindConfig = @"$kindConfigPath = Join-Path $rootDir "kind-config.yaml"$rootDir = Split-Path -Parent $scriptDir$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.PathRequire-Command kubectlRequire-Command kind}    }        throw "Required command not found: $Name"    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {    param([string]$Name)function Require-Command {$ErrorActionPreference = "Stop")    [int]$TimeoutSeconds = 240    [string]$TektonUrl = "https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml",    [string]$Namespace = "deckhand-lab",
command -v kind >/dev/null 2>&1 || {
	echo "ERROR: required command not found: kind" >&2
	exit 1
}

CLUSTER_NAME="${CLUSTER_NAME:-deckhand-dev}"
echo "Deleting kind cluster ${CLUSTER_NAME}"
kind delete cluster --name "${CLUSTER_NAME}" || true
echo "Done"
