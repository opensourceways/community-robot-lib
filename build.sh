#!/usr/bin/bash

set -euo pipefail

cd $(dirname $0)
me=$(basename $0)
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
        go mod init
        go mod tidy
    fi

    $bazel run //:gazelle -- update-repos -from_file=go.mod

    $bazel run //:gazelle
}

build(){
    if [ $# -ne 1 ]; then
        cmd_help build
        return 1
    fi

    local path=$1

    if [ ! -d $path ]; then
        echo "invalid path: $path"
        return 1
    fi

    update_repo

    tips "build lib"

    $bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //$path:go_default_library
}

clean(){
    $bazel clean
}

cmd_help(){
    if [ $# -eq 0 ]; then
cat << EOF
usage: $me cmd
supported cmd:
    clean: clean local environment
    build: build lib for each components.
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
            echo "$me build component-path"
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
