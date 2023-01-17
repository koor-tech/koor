#!/bin/bash

# script for resolving patching rook to koor changes, this uses setting up a
# worktree with latest rook/master and mangling it to adopt koor conventions,
# creating a temporary commit.

set -xe

KOORDIR=$(pwd)
ROOKDIR=$(pwd)/temp_rook_merge

#setup a temporary rook latest folder
git fetch rook master
git worktree add "${ROOKDIR}" --checkout rook/master
cd "${ROOKDIR}"
git checkout -b wip-temp-rook-merge

# ordering matters here
#TODO: github workflow's
find "${ROOKDIR}" -type f ! -path "./.git/*" ! -path "rook-to-koor.sh" ! -path "./.github/*" -exec sed -i \
    -e 's|rook/rook.github.io|koor-tech/docs.koor.tech|g' \
    -e 's|rook/rook/issues|koor-tech/koor/issues|g' \
    -e 's|/work/rook/rook|/work/koor/koor|g' \
    -e 's|/work/rook|/work/koor|g' \
    -r 's|rook/rook/blob|koor-tech/koor/blob|g' \
    -e 's|quay.io/ceph:myversion|koorinc/ceph:myversion|g' \
    -e '/cmd/! s|rook/rook|koor-tech/koor|g' \
    -e 's|on our \[Rook.io Slack\](https://slack.rook.io/)|in [the Github Discussions](https://github.com/koor-tech/koor/discussions)|g' \
    -e 's|Rook version|Koor Storage Distribution version|g' \
    -e 's|rook.io/docs/rook|docs.koor.tech/docs|g' \
    -e 's|"$GITHUB_REPOSITORY_OWNER" = "rook"|"$GITHUB_REPOSITORY_OWNER" = "koor-tech"|g' \
    -e 's|rook/ceph|koorinc/ceph|g' \
    -e 's|rook.chart|koor.chart|g' \
    -e 's|version of Rook|version of Koor Storage Distribution|g' \
    -e 's/Rook(.*?)[\s] website/Koor website/g' \
    -e 's|rook-release|koor-release|g' \
    -e 's|charts.rook.io|charts.koor.tech|g' \
    -e 's|Rook Ceph Operator|Koor Operator|g' \
    -e 's|Rook Ceph toolbox|Koor toolbox|g' \
    -e 's|cncf-rook-security@lists.cncf.io|security@koor.tech|g' \
    -e 's|cncf-rook-distributors-announce@lists.cncf.io|the Koor Newletter|g' \
    -e 's|rook_io|koor_tech|g' \
    -e 's|Rook Ceph |Koor |g' \
    -e 's/https\:\/\/slack\.rook\.io.*/\[the GitHub Discussions\]\(https\:\/\/github\.com\/koor\-tech\/koor\/discussions\)/g' \
    -e 's|src/github.com/rook|src/github.com/koor-tech|g' {} \;

# places where complete replacement are apt
sed -i "s|Rook|Koor Storage Distribution|g" "${ROOKDIR}"/.github/PULL_REQUEST_TEMPLATE.md \
    "${ROOKDIR}"/CODE_OF_CONDUCT.md \
    "${ROOKDIR}"/Documentation/Contributing/development-flow.md \
    "${ROOKDIR}"/README.md \
    .github/ISSUE_TEMPLATE.md \
    .github/ISSUE_TEMPLATE/bug_report.md \
    Documentation/Getting-Started/storage-architecture.md \

# release version mangling 1.10 and 1.9 -> 1.0
sed -i "s/release\-[0-9].[0-9]/release\-1\.0/g" "${ROOKDIR}"/Documentation/Contributing/development-flow.md
sed -i 's/v1\.10\../1\.0\.0/g' "${ROOKDIR}"/Documentation/Upgrade/rook-upgrade.md
sed -i 's/v1\.9\../v1\.0\.0/g' "${ROOKDIR}"/Documentation/Upgrade/rook-upgrade.md

# others
sed -i "s|rook|koor-tech|g" "${ROOKDIR}"/.docs/macros/includes/main.py
sed -i "s|koorinc/ceph|Koor Storage Distribution|g" "${ROOKDIR}"/Documentation/Troubleshooting/common-issues.md

# run gofmt
gofmt -w .

# apply patch and merge
git add -u
git commit -sm "conflict resolution for rook upstream merge"
cd "${KOORDIR}"
if ! git merge wip-temp-rook-merge; then
  echo "Please resolve all mergee conflicts(from other terminal window)"
  # common resolutions we'd want to keep koor code here
  git checkout --ours CODE-OWNERS ROADMAP.md
  git rm ADOPTERS.md Documentation/Contributing/storage-providers.md || true
  read  -n 1 -p "Press any key to continue..."
fi

# cleanup
rm -rf "${ROOKDIR}"
git worktree prune # only deletes worktree with no path
git branch -D wip-temp-rook-merge
