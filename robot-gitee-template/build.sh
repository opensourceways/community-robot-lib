#!/bin/bash

set -euo pipefail

cd $(dirname $0)
me=$(basename $0)
robot_name=$(pwd | xargs basename)
pn=$#
all_param=( $@ )

bazel=bazel-3.4.1

tips(){
    echo -e "\n*************** $1 ***************\n"
}

update_repo(){
    tips "update repo"

    if [ -f go.mod ]; then
        go mod tidy
    else
        go mod init {REPO_PATH}
        go mod tidy
    fi
}

build(){
    update_repo

    tips "build binary"

    go build -a -o $robot_name .
}

image(){
    update_repo

    tips "build image"

    image_sh="./publish/command_status.sh"
    test -f $image_sh || image_sh="./publish/image.sh"

    image_id=$(bash $image_sh | grep IMAGE_ID | awk '{print $2}')

    docker build -t $image_id .
}

cmd_help(){
    if [ $# -eq 0 ]; then
cat << EOF
usage: $me cmd
supported cmd:
    build: build binary.
    image: build image.
    help: show the usage for each commands.
EOF
        return 0
    fi

    local cmd=$1
    case $cmd in
        "build")
            echo "$me build"
            ;;
        "image")
            echo "$me image"
            ;;
        "help")
            echo "$me help other-child-cmd"
            ;;
        *)
            echo "unknown child cmd: $cmd"
            ;;
     esac
}

fetch_parameter() {
    local index=$1
    if [ $pn -lt $index ]; then
        echo ""
    else
        echo "${all_param[@]:${index}-1}"
    fi
}

if [ $pn -lt 1 ]; then
    cmd_help
    exit 1
fi

cmd=$1
case $cmd in
    "build")
        build $(fetch_parameter 2)
        ;;
    "image")
        image
        ;;
    "--help")
        cmd_help
        ;;
    "help")
        cmd_help $(fetch_parameter 2)
        ;;
    *)
        echo "unknown cmd: $cmd"
        ;;
esac
