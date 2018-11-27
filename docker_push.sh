#!/bin/bash

# log into docker
echo "Log into docker"
docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"
# push to docker
echo "Pushing to Docker"
docker push "$TRAVIS_REPO_SLUG":"$TAG"
if [[ "$TRAVIS_TAG" =~ ^v[0-9]. ]]; then
    docker push "$TRAVIS_REPO_SLUG":"latest"
fi
# Check for errors
if [ $? -eq 0 ]; then
	echo "Push Complete"
else
	echo "Push Failed"
	exit 1
fi
