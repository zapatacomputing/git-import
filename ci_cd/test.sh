#!/bin/bash

# We make sure to exit if the script fails and print the executed commands
set -e -x

mkdir ~/.ssh
eval "$(ssh-agent)"
echo "$ZAPATABOT_SSH" > ssh-secret
chmod 0600 ssh-secret
ssh-add ssh-secret
ssh-keyscan github.com >> ~/.ssh/known_hosts

# Ensure the git cli is setup with correct credentials
# This enables golang to get all the dependencies required to test the code
git config --global credential.helper store
echo "https://zapatabot:$GITHUB_PERSONAL_TOKEN@github.com" > ~/.git-credentials

cd git-import

# Runs unit tests as part of CI/CD framework.
echo "Testing with $(go version)"
ginkgo -v ./...
