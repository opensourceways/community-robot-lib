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

cat <<EOF
IMAGE_REGISTRY ${IMAGE_REGISTRY_OVERRIDE:-swr.ap-southeast-1.myhuaweicloud.com}
IMAGE_REPO ${IMAGE_REPO_OVERRIDE:-opensourceway/robot/$repository}
IMAGE_TAG ${IMAGE_TAG_OVERRIDE:-$image_tag}
EOF

cd $work_dir
