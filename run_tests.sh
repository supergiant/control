#!/bin/bash
echo "Running tests"

GOCACHE=off go test -race -covermode=atomic -coverprofile=profile.cov ./pkg/...
if [ $? -eq 0 ]; then
	echo "Tests Passed"
else
	echo "Tests Failed"
	exit 1
fi

goveralls -coverprofile=profile.cov -service=travis-ci
