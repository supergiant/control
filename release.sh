#!/usr/bin/env bash

## This should be temp, as our CI should be able to auto cut release for us.

RELEASE=$1


echo "Building bin..."
 if [ ! $(GOOS=linux GOARCH=amd64 go build -o build/supergiant ./..) ] ; then
   echo "Build successful..."
 else
   echo "Build failed..."
   exit 5
 fi

if [ $RELEASE == "" ]; then
  echo "Release not pushed to git, no release version was specified."
  exit 5
fi

github-release upload \
    --user supergiant \
    --repo supergiant \
    --tag $RELEASE \
    --name "supergiant-osx-amd64" \
    --file build/supergiant
