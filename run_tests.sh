#!/bin/bash
echo "Running tests"

go test -v -race ./pkg/...

# Check for errors
if [ $? -eq 0 ]; then
	echo "Tests Passed"
else
	echo "Tests Failed"
	exit 1
fi
