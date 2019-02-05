DOCKER_IMAGE_NAME := supergiant/control
DOCKER_IMAGE_TAG := $(shell git describe --tags --always | tr -d v || echo 'latest')
VERSION := $(shell git describe --always --long --dirty)

.PHONY: build test push release

fmt: gofmt goimports

gofmt:
	@FLAGS="-w" build/gofmt.sh

goimports:
	@FLAGS="-w -local github.com/supergiant/control" build/goimports.sh

lint:
	# TODO: enable the test directory when e2e tests will be updated
	build/gometalinter.sh

get-tools:
	go get -u github.com/kardianos/govendor
	go get -u github.com/alecthomas/gometalinter
	go get github.com/rakyll/statik
	gometalinter --install

build-image:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .
	docker build -t $(DOCKER_IMAGE_NAME):latest .

test:
	go test ./pkg/...

build:
	go get -u github.com/hpcloud/tail/...
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o dist/controlplane-linux -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s -X main.version=${VERSION}' ./cmd/controlplane
	GOOS=darwin CGO_ENABLED=0 GOARCH=amd64 go build -o dist/controlplane-osx -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s -X main.version=${VERSION}' ./cmd/controlplane
	GOOS=windows CGO_ENABLED=0 GOARCH=amd64 go build -o dist/controlplane-windows -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s -X main.version=${VERSION}' ./cmd/controlplane
push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

release: build push

build-ui:
	npm install --prefix ./cmd/ui/assets
	npm run build --prefix ./cmd/ui/assets
	statik -src=./cmd/ui/assets/dist
