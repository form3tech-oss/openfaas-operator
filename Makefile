.PHONY: build build-armhf push test verify-codegen ci-armhf-build ci-armhf-push ci-arm64-build ci-arm64-push

DOCKER_IMG ?= form3tech/openfaas-operator
DOCKER_TAG ?= $$(git describe --tags --dirty="-dev")

.PHONY: build
build:
	docker build --build-arg VERSION=$(DOCKER_TAG) -t $(DOCKER_IMG):$(DOCKER_TAG) . -f Dockerfile

.PHONY: push
push:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	docker push $(DOCKER_IMG):$(DOCKER_TAG)

test:
	go test ./...

verify-codegen:
	./hack/verify-codegen.sh
