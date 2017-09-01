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
  REPO=supergiant/supergiant-ui
  cp dist/supergiant-ui-linux-amd64 build/docker/ui/linux-amd64/
  cp dist/supergiant-ui-linux-arm64 build/docker/ui/linux-arm64/
  ls -l build/docker/ui/linux-amd64/
  ls -l build/docker/ui/linux-arm64/
  docker build -t $REPO:$TAG-amd build/docker/ui/linux-amd64/
  docker push $REPO
  docker build -t $REPO:$TAG-arm build/docker/ui/linux-arm64/
  docker push $REPO
fi
