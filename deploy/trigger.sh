#!/bin/bash

set -xeo pipefail

export USER="$1"
export HOME="/home/$USER"

# shellcheck source=/dev/null
source "$HOME/.bash_profile"

cd "$HOME/projects/website"

git fetch --prune
if git status --porcelain -b -u no | grep behind; then
	git merge origin/master
	./deploy/deploy.sh
fi
