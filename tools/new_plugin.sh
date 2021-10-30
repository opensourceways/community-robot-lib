#!/bin/bash

set -euo pipefail

cd $(dirname $0)
me=$(basename $0)
pn=$#

repo_name=""
prefix_of_robot_name="robot-gitee"

underscore_to_hyphen(){
    local name=$1
    echo ${name//_/-}
}

check_robot_name() {
    repo_name=$1
    local name=""

    name=$(echo "$repo_name" | awk '{print tolower($0)}')
    if [ "$name" != "$repo_name" ]; then
        echo "Info: the robot name($repo_name) includes uppercase characters, and will be changed to $name."
        repo_name=$name
    fi

    name=$(underscore_to_hyphen $repo_name)
    if [ "$name" != "$repo_name" ]; then
        echo "Info: the robot name($repo_name) includes '_', and will be changed to $name."
        repo_name=$name
    fi

    name=$(echo "$repo_name" | sed -e 's/[a-z0-9-]//g')
    if [ -n "$name" ]; then
        echo "Error: the robot name should only include characters of letter(a-z), digitals and '-'"
        return 1
    fi

    name=$(echo "$repo_name" | sed -e 's/--*/-/g')
    if [ "$name" != "$repo_name" ]; then
        echo "Info: the robot name($repo_name) includes multiple '-', and will be changed to $name."
        repo_name=$name
    fi

    name=${repo_name%-}
    if [ "$name" != "$repo_name" ]; then
        echo "Info: the robot name($repo_name) ends with '-', and will be changed to $name."
        repo_name=$name
    fi

    name=${repo_name//-/}
    local prefix=${prefix_of_robot_name//-/}

    name=${name#$prefix}
    if [ -z "$name" ]; then
        echo "Error: there is not real robot name. The '$prefix_of_robot_name' is the name prefix."
        return 1
    fi

    local s=${name/robot/-}
    if [ "$s" != "$name" ]; then
        echo "Error: the robot name can't include reserved word 'robot'."
        return 1
    fi

    local s=${name/gitee/-}
    if [ "$s" != "$name" ]; then
        echo "Error: the robot name can't include reserved word 'gitee'."
        return 1
    fi

    name=${repo_name#$prefix_of_robot_name}
    if [ "$repo_name" = "$name" ]; then
        repo_name="${prefix_of_robot_name}-${repo_name}"

        echo "Info: the robot name should have prefix of '$prefix_of_robot_name', and will be changed to $repo_name."
    fi
}

build(){
    local repo_dir=$1
    local repo_path=$2

    check_robot_name $(basename $repo_path)

    if [ -d $repo_dir ]; then
        echo "$repo_dir is exist"
        return 1
    fi

    mkdir $repo_dir
    cd $repo_dir

    git clone https://github.com/opensourceways/community-robot-lib.git

    cp -r community-robot-lib/new-plugin/. .

    rm -fr community-robot-lib 

    cp -r giteeplugin/. .

    rm -fr giteeplugin

    repo_path=${repo_path//\//\\\/}

    sed -i -e "s/{PLUGIN_REPO}/${repo_path}/" ./BUILD.bazel
    sed -i -e "s/{PLUGIN_NAME}/${repo_name}/" ./BUILD.bazel

    sed -i -e "s/{PLUGIN_REPO}/${repo_path}/" ./build.sh

    git init .
    git add .
    git commit -m "init repo"
}

cmd_help(){
cat << EOF
usage: $me dir-of-repo import-path-of-repo.
for example: $me ./test github.com/exmaple/test 
EOF
}

if [ $pn -lt 2 ]; then
    cmd_help
    exit 1
fi

build $1 $2
