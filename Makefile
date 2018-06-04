DOCKER_IMAGE_NAME := supergiant/supergiant
DOCKER_IMAGE_TAG := $(shell git describe --tags --always | tr -d v || echo 'latest')

.PHONY: build test push release

clean:
	rm -r bindata/ || true && rm -rf tmp/ && mkdir tmp

generate-bindata:
	go-bindata -pkg bindata -o bindata/bindata.go config/providers/... ui/assets/... ui/views/...

fmt: gofmt goimports

gofmt:
	@FLAGS="-w" build/verify/gofmt.sh

goimports:
	@FLAGS="-w -local github.com/supergiant/supergiant" build/verify/goimports.sh

lint:
	# TODO: enable the test directory when e2e tests will be updated
	@build/verify/gometalinter.sh

get-tools:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/kardianos/govendor
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

build-builder:
	docker build -t supergiant-builder --file build/Dockerfile.build .
	docker create --name supergiant-builder supergiant-builder
	rm -rf build/dist
	docker cp supergiant-builder:/go/src/github.com/supergiant/supergiant/dist build/dist
	docker rm supergiant-builder

build-image:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) --file build/Dockerfile .
	docker build -t $(DOCKER_IMAGE_NAME):latest --file build/Dockerfile .

test:
	docker build -t supergiant --file build/Dockerfile.build .
	docker run --rm supergiant govendor test +local

build: build-builder build-image

push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

release: build push
