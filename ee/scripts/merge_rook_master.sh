#!/bin/bash

set +x
set -e

if [[ -z "$GITHUB_TOKEN" ]]; then
    echo "Set the GITHUB_TOKEN env variable."
    exit 1
fi

set -x

export REMOTE_NAME="${REMOTE_NAME:-rook_fork}"
export SOURCE_BRANCH="${SOURCE_BRANCH:-master}"
export TARGET_BRANCH="${TARGET_BRANCH:-merge-rook-master}"

if git remote | grep --quiet "$REMOTE_NAME"; then
    git remote set-url "$REMOTE_NAME" https://github.com/rook/rook.git
else
    git remote add "$REMOTE_NAME" https://github.com/rook/rook.git
fi

# Fetch latest rook branch and switch to our rebase branch
# (if it doesn't exist locally yet, it will be reset)
git fetch "$REMOTE_NAME" "$SOURCE_BRANCH"
git checkout -B "$TARGET_BRANCH"

# Run the merge (ignore any failures)
git merge "$REMOTE_NAME/$SOURCE_BRANCH" || true

# Specific files/dirs that will be completely ours
declare -a OUR_FILES=(
    README.md
    Documentation/
)

for OUR_FILE in "${OUR_FILES[@]}"; do
    git checkout \
        --ours \
        origin/master -- "$OUR_FILE"
done

# Make sure to re-generate CRDs, Helm Chart Docs, etc., twice to be certain everything is really
# up-to-date with custom changes
make crds helm-docs generate-docs-crds
make crds helm-docs generate-docs-crds

# Stage all changed files
git add --all

git \
    -c core.editor=/bin/true \
    merge --continue || \
    { echo "Merge failed, try merging manual"; exit 1; }

# Push back
git push --force-with-lease origin "$TARGET_BRANCH"
