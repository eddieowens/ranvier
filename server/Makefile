SHELL :=/bin/bash

DOCKER_URL :=edwardrowens/ranvier-server
TAG := $(shell git rev-parse --short HEAD)
DOCKER_TAG :=$(TAG)
DOCKER_TAG_LATEST_ALPINE :=alpine
DOCKER_TAG_ALPINE :=$(TAG)-alpine

RANVIER_GIT_SSHKEY ?=~/.ssh/id_rsa

docker.build:
	docker build -f deploy/docker/Dockerfile -t $(DOCKER_URL):latest -t $(DOCKER_URL):$(DOCKER_TAG) -t ranvier .
	docker build -f deploy/docker/Dockerfile.alpine -t $(DOCKER_URL):$(DOCKER_TAG_LATEST_ALPINE) -t $(DOCKER_URL):$(DOCKER_TAG_ALPINE) -t ranvier:$(DOCKER_TAG_ALPINE) .

docker.push:
	docker push $(DOCKER_URL):latest
	docker push $(DOCKER_URL):$(DOCKER_TAG)
	docker push $(DOCKER_URL):$(DOCKER_TAG_LATEST_ALPINE)
	docker push $(DOCKER_URL):$(DOCKER_TAG_ALPINE)

docker: docker.build docker.push

docker.run.dev:
	docker run -p 8080:8080 -e RANVIER_GIT_SSHKEY=$(RANVIER_GIT_SSHKEY) -v $(RANVIER_GIT_SSHKEY):/.ssh/id_rsa $(DOCKER_URL)

swagger:
	go mod vendor
	swag init -o app/docs --parseVendor
	rm -rf vendor

compile: swagger
	go get
	go build -o server
