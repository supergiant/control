#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail
gometalinter --deadline=50s --vendor \
    --cyclo-over=50 --dupl-threshold=100 \
    --disable-all \
    --enable=vet \
    --enable=deadcode \
    --enable=golint \
    --enable=vetshadow \
    --enable=gocyclo \
    --enable=misspell \
    --skip=test \
    --skip=bindata \
    --skip=vendor \
    --tests \
    ./...
