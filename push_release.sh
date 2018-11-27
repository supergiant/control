#!/bin/bash

echo "building release artifacts"
PROJECTDIR=${TRAVIS_HOME}/gopath/src/github.com/${TRAVIS_REPO_SLUG}
echo ${PROJECTDIR}

mkdir -p ./supergiant-build

cp -avr ${PROJECTDIR}/templates/ ./supergiant-build
cp -avr ${PROJECTDIR}/cmd/ui/assets/dist/ ./supergiant-build
chmod +x ./supergiant-build/controlplane

ls -la ./supergiant-build

tar -czpf dist/assest.gz ./supergiant-build

# if a tag has alpha or beta in the name, it will be released as a pre-release.
# if a tag does not have alpha or beta, it is pushed as a full release.
case "${TAG}" in
	*alpha* )  echo "Releasing version: ${TAG}, as pre-release"
	ghr --username supergiant --token "$GITHUB_TOKEN" --replace -b "pre-release" --prerelease --debug "$TAG"  dist/;;
	*beta* )    echo "Releasing version: ${TAG}, as pre-release"
	ghr --username supergiant --token "$GITHUB_TOKEN" --replace -b "pre-release" --prerelease --debug "$TAG"   dist/;;
	*)echo "Releasing version: ${TAG}, as latest release."
	ghr --username supergiant --token "$GITHUB_TOKEN" --replace -b "latest release" --debug "$TAG"   dist/;;
esac

# Check for errors
if [ $? -eq 0 ]; then
	echo "Release pushed"
else
	echo "Push to releases failed"
	exit 1
fi
