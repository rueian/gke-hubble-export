export IMAGE_TAG := $(if $(IMAGE_TAG),$(IMAGE_TAG),latest)
export DOCKER_BUILDKIT := 1
export COMPOSE_DOCKER_CLI_BUILD := 1
export DOCKER_CONTENT_TRUST := 1 

# Build local docker image
.PHONY: docker
docker:
	docker build -f Dockerfile -t rueian/gke-hubble:${IMAGE_TAG} .
