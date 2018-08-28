DOCKER_IMAGE_NAME := supergiant/supergiant
DOCKER_IMAGE_TAG := $(shell git describe --tags --always | tr -d v || echo 'latest')

.PHONY: build test push release

fmt: gofmt goimports

gofmt:
	@FLAGS="-w" build/gofmt.sh

goimports:
	@FLAGS="-w -local github.com/supergiant/supergiant" build/goimports.sh

lint:
	# TODO: enable the test directory when e2e tests will be updated
	build/gometalinter.sh

get-tools:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

build-image:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .
	docker build -t $(DOCKER_IMAGE_NAME):latest .

test:
	go test ./pkg/...

build: build-image

push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

release: build push
