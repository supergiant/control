#!/bin/bash
if [[ ! "$TRAVIS_TAG" =~ ^v[0-100]. ]]; then
  echo "Releasing supergiant version: ${VERSION}, pre-release"
  ghr --username supergiant --token $GITHUB_TOKEN --replace --prerelease --debug unstable-$TRAVIS_TAG dist/
  exit 0
elif [ ! -z "$TRAVIS_TAG" ]; then
  echo "Releasing supergiant version: ${TRAVIS_TAG}, as latest release."
  ghr --username supergiant --token $GITHUB_TOKEN --replace --debug $TRAVIS_TAG dist/
  exit 0
fi
echo "Unable to determine tag."
exit 5
