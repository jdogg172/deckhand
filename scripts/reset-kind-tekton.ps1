param(
    [string]$ClusterName = "deckhand-dev"
)

$ErrorActionPreference = "Stop"

if (-not (Get-Command kind -ErrorAction SilentlyContinue)) {
    throw "Required command not found: kind"
}

Write-Host "Deleting kind cluster $ClusterName"
kind delete cluster --name $ClusterName
Write-Host "Done"
