SHELL :=/bin/bash

DOCKER_URL :=edwardrowens/ranvier-server
DOCKER_TAG :=$(shell git rev-parse --short HEAD)

RANVIER_GIT_SSHKEY ?=~/.ssh/id_rsa

docker.build:
	docker build -t $(DOCKER_URL):latest -t $(DOCKER_URL):$(DOCKER_TAG) -t ranvier .

docker.push:
	docker push $(DOCKER_URL):latest
	docker push $(DOCKER_URL):$(DOCKER_TAG)

docker: docker.build docker.push

docker.run.dev:
	docker run -p 8080:8080 -e RANVIER_GIT_SSHKEY=$(RANVIER_GIT_SSHKEY) -v `dirname $(RANVIER_GIT_SSHKEY)`:/.ssh $(DOCKER_URL)

swagger:
	go mod vendor
	swag init -o app/docs --parseVendor
	rm -rf vendor

compile: swagger
	go get
	go build -o server
