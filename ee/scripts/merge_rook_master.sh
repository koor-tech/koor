#!/bin/bash

set -x

ROOK_BRANCH="master"

if [[ -z "$GITHUB_TOKEN" ]]; then
	echo "Set the GITHUB_TOKEN env variable."
	exit 1
fi

git remote add rook_fork https://github.com/rook/rook.git
git fetch rook_fork $ROOK_BRANCH

set -o xtrace

# do the merge
git checkout -b merge-rook-master

git merge rook_fork/$ROOK_BRANCH
#specific files that will be completely ours
git checkout --ours rook_fork/master -- README.md
(git -c core.editor=/bin/true merge --continue) || (echo "Merge failed, try merging manual" && exit 1)

# push back
git push --force-with-lease origin merge-rook-master
