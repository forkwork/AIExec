# Variables
DOCKER_REGISTRY=khulnasoft
WEB_IMAGE=$(DOCKER_REGISTRY)/aiexec-web
API_IMAGE=$(DOCKER_REGISTRY)/aiexec-api
SANDBOX_IMAGE=$(DOCKER_REGISTRY)/aiexec-sandbox
VERSION=latest

# Build Docker images
build-web:
	@echo "Building web Docker image: $(WEB_IMAGE):$(VERSION)..."
	docker build -t $(WEB_IMAGE):$(VERSION) ./web
	@echo "Web Docker image built successfully: $(WEB_IMAGE):$(VERSION)"

build-api:
	@echo "Building API Docker image: $(API_IMAGE):$(VERSION)..."
	docker build -t $(API_IMAGE):$(VERSION) ./api
	@echo "API Docker image built successfully: $(API_IMAGE):$(VERSION)"

build-sandbox:
	@echo "Building sandbox Docker image: $(SANDBOX_IMAGE):$(VERSION)..."
	cd sandbox && ./build/build_amd64.sh
	docker build -t $(SANDBOX_IMAGE):$(VERSION) -f sandbox/docker/amd64/dockerfile .
	@echo "Sandbox Docker image built successfully: $(SANDBOX_IMAGE):$(VERSION)"

# Push Docker images
push-web:
	@echo "Pushing web Docker image: $(WEB_IMAGE):$(VERSION)..."
	docker push $(WEB_IMAGE):$(VERSION)
	@echo "Web Docker image pushed successfully: $(WEB_IMAGE):$(VERSION)"

push-api:
	@echo "Pushing API Docker image: $(API_IMAGE):$(VERSION)..."
	docker push $(API_IMAGE):$(VERSION)
	@echo "API Docker image pushed successfully: $(API_IMAGE):$(VERSION)"

push-sandbox:
	@echo "Pushing sandbox Docker image: $(SANDBOX_IMAGE):$(VERSION)..."
	docker push $(SANDBOX_IMAGE):$(VERSION)
	@echo "Sandbox Docker image pushed successfully: $(SANDBOX_IMAGE):$(VERSION)"

# Build all images
build-all: build-web build-api build-sandbox

# Push all images
push-all: push-web push-api push-sandbox

build-push-api: build-api push-api
build-push-web: build-web push-web
build-push-sandbox: build-sandbox push-sandbox

# Build and push all images
build-push-all: build-all push-all
	@echo "All Docker images have been built and pushed."

# Phony targets
.PHONY: build-web build-api build-sandbox push-web push-api push-sandbox build-all push-all build-push-all build-push-sandbox