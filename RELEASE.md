# Release

## Create New Major/ Minor Release Branch

```console
branch_name=release-1.11
git checkout master
git fetch --all
git checkout -b $branch_name
git push origin $branch_name
```

### Create Alpha/Beta Release Version/ Tag

```console
tag_name=v1.11.0-alpha.0
git tag -a $tag_name -m "$tag_name minor alpha release tag"
git push origin $tag_name
```

## Create Release Version/ Tag

```console
tag_name=v1.11.0
git tag -a $tag_name -m "$tag_name minor release tag"
git push origin $tag_name
```

## Create Patch Release Version/ Tag

```console
tag_name=v1.11.1
git tag -a $tag_name -m "$tag_name patch release tag"
git push origin $tag_name
```
