#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd -- "${SCRIPT_DIR}/.." && pwd)"

CLUSTER_NAME="${CLUSTER_NAME:-deckhand-dev}"
KIND_CONFIG="${KIND_CONFIG:-${ROOT_DIR}/kind-config.yaml}"
TEKTON_VERSION_URL="${TEKTON_VERSION_URL:-https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml}"
NS_APP="${NS_APP:-deckhand-lab}"
KUBECTL_TIMEOUT="${KUBECTL_TIMEOUT:-240s}"

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "ERROR: required command not found: $1" >&2
    exit 1
  }
}

require_cmd kind
require_cmd kubectl

cat > "${KIND_CONFIG}" <<'EOF'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: deckhand-dev
nodes:
  - role: control-plane
  - role: worker
  - role: worker
EOF

# Keep cluster name configurable while preserving default config behavior.
sed -i.bak "s/^name: deckhand-dev$/name: ${CLUSTER_NAME}/" "${KIND_CONFIG}" || true
rm -f "${KIND_CONFIG}.bak"

echo "[1/6] Creating kind cluster: ${CLUSTER_NAME}"
if kind get clusters | grep -qx "${CLUSTER_NAME}"; then
  echo "Cluster ${CLUSTER_NAME} already exists; reusing"
else
  kind create cluster --name "${CLUSTER_NAME}" --config "${KIND_CONFIG}"
fi

echo "[2/6] Creating app namespace"
kubectl create namespace "${NS_APP}" || true

echo "[3/6] Installing Tekton Pipelines"
kubectl apply -f "${TEKTON_VERSION_URL}"

echo "[4/6] Waiting for Tekton controller rollout"
kubectl -n tekton-pipelines rollout status deploy/tekton-pipelines-controller --timeout="${KUBECTL_TIMEOUT}"
kubectl -n tekton-pipelines rollout status deploy/tekton-pipelines-webhook --timeout="${KUBECTL_TIMEOUT}"

echo "[5/6] Applying sample manifests"
kubectl apply -n "${NS_APP}" -f "${ROOT_DIR}/manifests/tekton-task-hello.yaml"
kubectl apply -n "${NS_APP}" -f "${ROOT_DIR}/manifests/tekton-pipeline-hello.yaml"
kubectl apply -n "${NS_APP}" -f "${ROOT_DIR}/manifests/tekton-pipelinerun-success.yaml"
kubectl apply -n "${NS_APP}" -f "${ROOT_DIR}/manifests/tekton-pipelinerun-fail.yaml"
kubectl apply -n "${NS_APP}" -f "${ROOT_DIR}/manifests/sample-crashloop-pod.yaml"
kubectl apply -n "${NS_APP}" -f "${ROOT_DIR}/manifests/sample-pending-pod.yaml"

echo "[6/6] Done"
kubectl config current-context
kubectl get pods -A
kubectl get pipelineruns -n "${NS_APP}" || true
