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
  cp dist/supergiant-ui-darwin-10.6-amd64 build/docker/ui/darwin-amd64/
  cp dist/supergiant-ui-windows-4.0-amd64.exe build/docker/ui/windows-amd64/
  cp dist/supergiant-ui-linux-arm64 build/docker/ui/linux-arm64/
  docker build -t $REPO:$TAG-linux-x64 build/docker/ui/linux-amd64/
  docker build -t $REPO:$TAG-darwin-x64 build/docker/ui/linux-amd64/
  docker build -t $REPO:$TAG-windows-x64 build/docker/ui/linux-amd64/
  docker build -t $REPO:$TAG-linux-arm64 build/docker/ui/linux-arm64/
  docker push $REPO

  ## API Docker Build
  REPO=supergiant/supergiant-api
  cp dist/supergiant-server-linux-amd64 build/docker/api/linux-amd64/
  cp dist/supergiant-server-darwin-10.6-amd64 build/docker/api/darwin-amd64/
  cp dist/supergiant-server-windows-4.0-amd64.exe build/docker/api/windows-amd64/
  cp dist/supergiant-server-linux-arm64 build/docker/api/linux-arm64/
  docker build -t $REPO:$TAG-linux-x64 build/docker/api/linux-amd64/
  docker build -t $REPO:$TAG-darwin-x64 build/docker/api/linux-amd64/
  docker build -t $REPO:$TAG-windows-x64 build/docker/api/linux-amd64/
  docker build -t $REPO:$TAG-linux-arm64 build/docker/api/linux-arm64/
  docker push $REPO
fi
