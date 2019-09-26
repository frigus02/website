#!/bin/bash

export USER="$1"
export HOME="/home/$USER"

# shellcheck source=/dev/null
source "$HOME/.bash_profile"

set -xeo pipefail

cd "$HOME/projects/website"

git fetch --prune
if git status --porcelain -b -u no | grep behind; then
	echo "Updates available. Starting deployment"
	git merge origin/master
	./deploy/deploy.sh
else
	echo "Already up to date"
fi
