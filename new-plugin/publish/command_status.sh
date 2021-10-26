#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

work_dir=$(pwd)
cd $(dirname $0)

commit_id=$(git describe --tags --always --dirty)
build_date=$(date -u '+%Y%m%d')
image_tag="v${build_date}-${commit_id}"
repository=$(pwd | xargs dirname | xargs basename)

image_registry=${IMAGE_REGISTRY_OVERRIDE:-swr.cn-north-4.myhuaweicloud.com}
image_repo=${IMAGE_REPO_OVERRIDE:-opensourceway/robot/$repository}
image_tag=${IMAGE_TAG_OVERRIDE:-$image_tag}

cat <<EOF
IMAGE_REGISTRY ${image_registry}
IMAGE_REPO ${image_repo}
IMAGE_TAG ${image_tag}
IMAGE_ID ${image_registry}/${image_repo}:${image_tag}
EOF

cd $work_dir
