#!/bin/bash

set -ex

# Set to something sane, so the script can be run outside of the CI as well
if [ -z ${GITHUB_STEP_SUMMARY+x} ]; then
    export GITHUB_STEP_SUMMARY="/tmp/ci-summary"
fi

export REMOTE_NAME="rook_fork"
export TARGET_REPO="$GITHUB_REPOSITORY"
export SOURCE_BRANCH="master"
export TARGET_BRANCH="merge-rook-master"

./ee/scripts/merge_rook_master.sh

# List current open PRs
pr_list=$(gh pr list --repo "$TARGET_REPO" --base master --head merge-rook-master --state open --json url)
prs=$(echo "$pr_list" | jq -r '. | length')

# Create rebase PR if none exist yet
if (( prs == 0 )); then
    if pr_url=$(gh pr create --title 'Merge Rook upstream to Koor' --body 'Created by Github action Backport rook/rook upstream cronjob'); then
        echo ":new_moon: Rebase PR opened $pr_url" >> $GITHUB_STEP_SUMMARY
    else
        echo "gh pr create Output: $pr_url"
        echo ":x: Failed to create Rebase PR!" >> $GITHUB_STEP_SUMMARY
    fi
else
    pr_url="$(echo "$pr_list" | jq -r '. | first | .url')"
    echo ":pushpin: Rebase PR already opened $pr_url" >> $GITHUB_STEP_SUMMARY
fi
