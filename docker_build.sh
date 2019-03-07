#!/bin/bash

echo "$TRAVIS_REPO_SLUG":"$TAG"
# build the docker container
echo "Building Docker container"

docker build --tag "$TRAVIS_REPO_SLUG":"$TAG" .
if [[ "$TRAVIS_TAG" =~ ^v[0-9]. ]]; then
    docker build --build-arg TAG=$(git describe --tags) --tag "$TRAVIS_REPO_SLUG":"latest" .
fi

if [ $? -eq 0 ]; then
	echo "Complete"
else
	echo "Build Failed"
	exit 1
fi
