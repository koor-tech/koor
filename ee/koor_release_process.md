## How to create a Koor release?

We follow the release process of Rook upstream and also sync our latest releases with rook releases.

* Create a latest release branch for Koor with major version in it's name, eg: release-1.11
* Pull the latest Rook release changes to this branch and make sure all Koor PRs are merged.

  ```
  # 1. checkout Koor master
  git checkout koor/master
  # 2. merge all relevant Koor PRs
  # 3. pull rook corresponding release changes
  git pull rook release-$VERSION
  # 4. merge all relevant Koor PRs
  # 5. create new release branch
  git checkout <branch> # e.g. release-1.10
  # set to the new release
  tag_name=<release version> # e.g., v1.10.9
  git tag -a "$tag_name" -m "$tag_name release tag"
  git push upstream "$tag_name"
  ```

* Follow the steps documented in [Rook release process](https://github.com/rook/rook/tree/master/build/release) from step 5 i.e. about generating release notes.
