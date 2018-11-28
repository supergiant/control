#!/bin/bash

export TAG=${TRAVIS_BRANCH:-unstable}

# checks for success of the previous task
check_status () {
	if [ $? -eq 0 ]; then
		echo "Success"
	else
		echo "Something failed"
		exit 1
	fi
}

### Main
echo "Tag Name: ${TAG}"
# If a tag is pushed, tests are run, the docker container is built and pushed to
# dockerhub, and then a release is pushed to the releases page.
if [[ "$TRAVIS_TAG" =~ ^v[0-9]. ]]; then
	echo "release"
	# run tests
	./run_tests.sh
	check_status
	# Build Docker container
	./docker_build.sh
	check_status
	# Push to Dockerhub
	./docker_push.sh
	check_status
	# Push to releases page
	./push_release.sh
	check_status
# on an unstable branch, tests are run and the docker container is built and pushed.
elif [[ "$TRAVIS_BRANCH" == *release-* ]]; then
	echo "unstable branch"
	export TAG="${TAG}-unstable"
	echo "Tag Name: ${TAG}"
	# run tests
	./run_tests.sh
	check_status
	# Build docker container
	./docker_build.sh
	check_status
	# Push to Dockerhub
	./docker_push.sh
	check_status
# if a push to master happens, tests are only run
elif [[ "$TRAVIS_BRANCH" == "master" ]]; then
	echo "master branch - test will only be run"
	echo "Tag Name: ${TAG}"
	# run tests
	./run_tests.sh
	check_status
else
# any other branch is considered a testing branch and will only run tests and build the container.
	echo "testing branch - run tests and docker build"
	export TAG="${TAG}-testing"
	echo "Tag Name: ${TAG}"
	# run tests
	./run_tests.sh
	check_status
	# Build docker container
	./docker_build.sh
	check_status
fi
