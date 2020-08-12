#!/bin/bash

# We make sure to exit if the script fails and print the executed commands
set -e -x

TAG=`cat current-release/tag`

echo "Current Release Tag: $TAG"

if [ -z ${NEW_TAG+x} ]
then 
    NEW_TAG=`echo $TAG | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g'`
fi

echo "New Release Tag: $NEW_TAG"

echo "$NEW_TAG" > release-artifacts/tag
echo "Git Import $NEW_TAG" > release-artifacts/name

cd git-import
echo "Building with $(go version)"
for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    echo "--------------------------------"
    export GOOS GOARCH
    go build -mod=vendor -v -o git-import-$NEW_TAG-$GOOS-$GOARCH
  done
done

cp git-import-* ../release-artifacts/
