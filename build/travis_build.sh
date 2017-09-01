#!/usr/bin/env bash

# Travis deployment script. After test success actions go here.

TAG=${TRAVIS_BRANCH:-unstable}


echo "Tag Name: ${TAG}"
if [[ "$TAG" =~ ^v[0-100]. ]]; then
  echo "global deploy"
  ./packer build build/build_release.json
else
  echo "private unstable"
  docker login -u $DOCKER_USER -p $DOCKER_PASS

  ## UI Docker Build
  REPO=supergiant/supergiant-ui
  cp dist/supergiant-ui-linux-amd64 build/docker/ui/linux-amd64/
  cp dist/supergiant-ui-darwin-amd64 build/docker/ui/darwin-amd64/
  cp dist/supergiant-ui-windows-amd64 build/docker/ui/windows-amd64/
  cp dist/supergiant-ui-linux-arm64 build/docker/ui/linux-arm64/
  docker build -t $REPO:unstable-$TAG-linux build/docker/ui/linux-amd64/
  docker build -t $REPO:unstable-$TAG-darwin build/docker/ui/linux-amd64/
  docker build -t $REPO:unstable-$TAG-windows build/docker/ui/linux-amd64/
  docker push $REPO
  docker build -t $REPO:unstable-$TAG-arm build/docker/ui/linux-arm64/
  docker push $REPO

  ## API Docker Build
  REPO=supergiant/supergiant-api
  cp dist/supergiant-server-linux-amd64 build/docker/api/linux-amd64/
  cp dist/supergiant-server-darwin-amd64 build/docker/api/darwin-amd64/
  cp dist/supergiant-server-windows-amd64 build/docker/api/windows-amd64/
  cp dist/supergiant-server-linux-arm64 build/docker/api/linux-arm64/
  docker build -t $REPO:unstable-$TAG-linux build/docker/api/linux-amd64/
  docker build -t $REPO:unstable-$TAG-darwin build/docker/api/linux-amd64/
  docker build -t $REPO:unstable-$TAG-windows build/docker/api/linux-amd64/
  docker push $REPO
  docker build -t $REPO:unstable-$TAG-arm build/docker/api/linux-arm64/
  docker push $REPO
fi
