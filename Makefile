APP=deckhand
DIST=dist

.PHONY: build run tidy test clean fmt env-kind env-kind-reset env-samples

build:
	mkdir -p $(DIST)
	go build -o $(DIST)/$(APP) ./cmd/deckhand

run:
	go run ./cmd/deckhand

tidy:
	go mod tidy

test:
	go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf $(DIST)

env-kind:
	bash ./scripts/setup-kind-tekton.sh

env-kind-reset:
	bash ./scripts/reset-kind-tekton.sh

env-samples:
	kubectl apply -n deckhand-lab -f manifests/tekton-task-hello.yaml
	kubectl apply -n deckhand-lab -f manifests/tekton-pipeline-hello.yaml
	kubectl apply -n deckhand-lab -f manifests/tekton-pipelinerun-success.yaml
	kubectl apply -n deckhand-lab -f manifests/tekton-pipelinerun-fail.yaml
	kubectl apply -n deckhand-lab -f manifests/sample-crashloop-pod.yaml
	kubectl apply -n deckhand-lab -f manifests/sample-pending-pod.yaml
