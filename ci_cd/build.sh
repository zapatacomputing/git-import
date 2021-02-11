#!/bin/bash
# We make sure to exit if the script fails and print the executed commands
set -e -x
# Get current directory to make it easier to reference
ROOT_DIR=`pwd`
cd git-import
# Build all the neeeded binaries
echo "Building with $(go version)"
go get -u github.com/mitchellh/gox
for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    export GOOS GOARCH
    go build -v -o git-import-$GOOS-$GOARCH
  done
  echo "------------------------------------------------"
done
# Move them back to ROOT to make it easier to gunzip
mv  git-import-* $ROOT_DIR
# Get commit_sha for future reference by Github Releases
cat .git/ref > $ROOT_DIR/commit_sha
# Set the BuildNumber using Date_Timestamp
BUILD_NUMBER=$(date -u +%Y%m%d.%H%M%S)
if [ -z ${MANUAL_VERSION+x} ]
then 
    RELEASE_VERSION=${STAR_VERSION}
else
    # If MANUAL_VERSION is set, override the star name.
    RELEASE_VERSION=`echo ${MANUAL_VERSION}`
fi
# Create new version with ReleaseVersion and BuildNumber
UPDATED_VERSION="${RELEASE_VERSION}.${BUILD_NUMBER}"
echo "New Release Tag: $UPDATED_VERSION"
echo "$UPDATED_VERSION" > $ROOT_DIR/tag
cd $ROOT_DIR
# put everything related to the build into one gunzip file:
tar -P -cvzf release-artifacts/build-$UPDATED_VERSION.tgz git-import-* commit_sha tag git-import
echo "Build Archive Created: release-artifacts/build-$UPDATED_VERSION.tgz"