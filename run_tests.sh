#!/bin/bash
echo "Running tests"

go test -v -race -covermode=count -coverprofile=profile.cov ./pkg/...
goveralls -coverprofile=profile.cov -service=travis-ci

# Check for errors
if [ $? -eq 0 ]; then
	echo "Tests Passed"
else
	echo "Tests Failed"
	exit 1
fi
