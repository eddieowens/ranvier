DOCKER_URL :=edwardrowens/ranvier/server
DOCKER_TAG :=$(shell git rev-parse --short HEAD)

docker.build:
	docker build -t $(DOCKER_URL):latest -t $(DOCKER_URL):$(DOCKER_TAG) -t ranvier .

docker.push:
	docker push $(DOCKER_URL):latest
	docker push $(DOCKER_URL):$(DOCKER_TAG)

docker: docker.build docker.push

docker.run.dev:
	docker run -p 8080:8080 $(DOCKER_URL)

swagger:
	go mod vendor
	swag init -o app/docs --parseVendor
	rm -rf vendor

compile: swagger
	go get
	go build -o server
