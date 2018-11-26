#!/bin/bash

echo "$TRAVIS_REPO_SLUG":"$TAG"
# build the docker container
echo "Building Docker container"
docker build --tag build --target build .
docker build --tag "$TRAVIS_REPO_SLUG":"$TAG" .

if [ $? -eq 0 ]; then
	echo "Complete"
else
	echo "Build Failed"
	exit 1
fi
