#!/bin/bash

set -euo pipefail

cd $(dirname $0)
me=$(basename $0)
plugin_name=$(pwd | xargs basename)
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
        go mod init {PLUGIN_REPO}
        go mod tidy
    fi

    $bazel run //:gazelle -- update-repos -from_file=go.mod

    $bazel run //:gazelle
}

build(){
    update_repo

    tips "build binary"

    $bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //:$plugin_name
}

image(){
    update_repo

    tips "build image"

    $bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //:image
}

push_image(){
    update_repo

    tips "push image"

    $bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //:push_image
}

clean(){
    $bazel clean
}

cmd_help(){
    if [ $# -eq 0 ]; then
cat << EOF
usage: $me cmd
supported cmd:
    clean: clean local environment.
    build: build binary.
    image: build image.
    push_image: build and push image.
    help: show the usage for each commands.
EOF
        return 0
    fi

    local cmd=$1
    case $cmd in
        "clean")
            echo "$me clean"
            ;;
        "build")
            echo "$me build"
            ;;
        "image")
            echo "$me image"
            ;;
        "push_image")
            echo "$me push_image"
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
    "clean")
        clean
        ;;
    "build")
        build $(fetch_parameter 2)
        ;;
    "image")
        image
        ;;
    "push_image")
        push_image
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
