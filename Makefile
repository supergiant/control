DOCKER_IMAGE_NAME := supergiant/supergiant
DOCKER_IMAGE_TAG := $(shell git describe --tags --always | tr -d v || echo 'latest')

.PHONY: build test push release

clean:
	rm -r bindata/ || true && rm -rf tmp/ && mkdir tmp

generate-bindata:
	go-bindata -pkg bindata -o bindata/bindata.go config/providers/... ui/assets/... ui/views/...

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
