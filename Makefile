export GO111MODULE=on

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
	build/golangci-lint.sh

get-tools:
	go get -u golang.org/x/tools/cmd/goimports
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1
	go get github.com/rakyll/statik

build-image:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

test:
	go test -mod=vendor ./pkg/...

build:
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -mod=vendor -o dist/controlplane-linux -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s -X main.version=${VERSION}' ./cmd/controlplane
	GOOS=darwin CGO_ENABLED=0 GOARCH=amd64 go build -mod=vendor -o dist/controlplane-osx -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s -X main.version=${VERSION}' ./cmd/controlplane
	GOOS=windows CGO_ENABLED=0 GOARCH=amd64 go build -mod=vendor -o dist/controlplane-windows -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s -X main.version=${VERSION}' ./cmd/controlplane
push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

release: build push

build-ui:
	npm install --prefix ./cmd/ui/assets
	npm run build --prefix ./cmd/ui/assets
	statik -src=./cmd/ui/assets/dist

gogen:
	go -mod=vendor generate ./pkg/account

vendor-sync:
	go mod tidy
	go mod download
	go mod vendor
