#!/usr/bin/env bash

# Travis deployment script. After test success actions go here.

TAG=${TRAVIS_BRANCH:-unstable}
REPO=supergiant/supergiant-ui

echo "Tag Name: ${TAG}"
if [[ "$TAG" =~ ^v[0-100]. ]]; then
  echo "global deploy"
  ./packer build build/build_release.json
else
  echo "private unstable"
  docker login -e $DOCKER_EMAIL -u $DOCKER_USER -p $DOCKER_PASS
  cp dist/supergiant-ui-linux-amd64 build/docker/ui/linux-amd64/dist/supergiant-ui-linux-amd64
  docker build -t $REPO:$TAG build/docker/ui/linux-amd64/
  docker push $REPO
fi
