#!/bin/bash

if [[ -z "$GITHUB_TOKEN" ]]; then
	echo "Set the GITHUB_TOKEN env variable."
	exit 1
fi


ROOK_REMOTE="https://github.com/rook/rook.git"
ROOK_BRANCH="master"
INPUT_AUTOSQUASH=false

git remote add fork "$ROOK_REMOTE"
git fetch $ROOK_REMOTE $ROOK_BRANCH

set -o xtrace

# do the rebase
git checkout -b fork/rebase-rook-master fork/rebase-rook-master
if [[ $INPUT_AUTOSQUASH == 'true' ]]; then
	GIT_SEQUENCE_EDITOR=: git rebase -i --autosquash $ROOK_REMOTE/$ROOK_BRANCH
else
	git rebase $ROOK_REMOTE/$ROOK_BRANCH
fi

# push back
git push --force-with-lease fork fork/rebase-rook-master:rebase-rook-master
