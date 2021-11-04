#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

work_dir=$(pwd)
cd $(dirname $0)

branch=$(git rev-parse --abbrev-ref HEAD)
branch=${BRANCH_OVERRIDE:-$branch}

commit_id=$(git describe --tags --always --dirty)
commit_id=${COMMIT_ID_OVERRIDE:-$commit_id}

repository=$(git remote -v | tail -1 | awk -F '/' '{print $NF}' | awk -F ' ' '{print $1}' | awk -F '.' '{print $1}')
if [ -z "$repository" ]; then
    repository=$(pwd | xargs dirname | xargs basename)
fi
repository=${REPOSITORY_OVERRIDE:-$repository}

image_tag="${branch}-${commit_id}"
image_tag=${IMAGE_TAG_OVERRIDE:-$image_tag}

image_registry=${IMAGE_REGISTRY_OVERRIDE:-swr.cn-north-4.myhuaweicloud.com}
image_repo=${IMAGE_REPO_OVERRIDE:-opensourceway/robot/$repository}

cat <<EOF
IMAGE_REGISTRY ${image_registry}
IMAGE_REPO ${image_repo}
IMAGE_TAG ${image_tag}
IMAGE_ID ${image_registry}/${image_repo}:${image_tag}
CODE_REPOSITORY ${repository}
EOF

cd $work_dir
