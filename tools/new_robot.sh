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

log(){
    local s=$1
    echo -e "\n${s}\n"
}

check_robot_name() {
    repo_name=$1
    local name=""

    name=$(echo "$repo_name" | awk '{print tolower($0)}')
    if [ "$name" != "$repo_name" ]; then
        log "Info: the robot name($repo_name) includes uppercase characters, and will be changed to $name."
        repo_name=$name
    fi

    name=$(underscore_to_hyphen $repo_name)
    if [ "$name" != "$repo_name" ]; then
        log "Info: the robot name($repo_name) includes '_', and will be changed to $name."
        repo_name=$name
    fi

    name=$(echo "$repo_name" | sed -e 's/[a-z0-9-]//g')
    if [ -n "$name" ]; then
        log "Error: the robot name should only include characters of letter(a-z), digitals and '-'"
        return 1
    fi

    name=$(echo "$repo_name" | sed -e 's/--*/-/g')
    if [ "$name" != "$repo_name" ]; then
        log "Info: the robot name($repo_name) includes multiple '-', and will be changed to $name."
        repo_name=$name
    fi

    name=${repo_name%-}
    if [ "$name" != "$repo_name" ]; then
        log "Info: the robot name($repo_name) ends with '-', and will be changed to $name."
        repo_name=$name
    fi

    name=${repo_name//-/}
    local prefix=${prefix_of_robot_name//-/}

    name=${name#$prefix}
    if [ -z "$name" ]; then
        log "Error: there is not real robot name. The '$prefix_of_robot_name' is the name prefix."
        return 1
    fi

    local s=${name/robot/-}
    if [ "$s" != "$name" ]; then
        log "Error: the robot name can't include reserved word 'robot'."
        return 1
    fi

    local s=${name/gitee/-}
    if [ "$s" != "$name" ]; then
        log "Error: the robot name can't include reserved word 'gitee'."
        return 1
    fi

    name=${repo_name#$prefix_of_robot_name}
    if [ "$repo_name" = "$name" ]; then
        repo_name="${prefix_of_robot_name}-${repo_name}"

        log "Info: the robot name should have prefix of '$prefix_of_robot_name', and will be changed to $repo_name."
    fi
}

build(){
    local robot_name=$1
    local repo_dir=$2
    local remote_repo=$3

    check_robot_name $robot_name
    robot_name=$repo_name

    repo_dir=$repo_dir/$robot_name
    if [ -d $repo_dir ]; then
        log "$repo_dir is exist"
        return 1
    fi

    mkdir -p $repo_dir
    cd $repo_dir

    git clone https://github.com/opensourceways/community-robot-lib.git

    cp -r community-robot-lib/robot-gitee-template/. .

    rm -fr community-robot-lib 

    cp -r basic/. .

    rm -fr basic


    repo_path=${remote_repo}/$robot_name

    git init .
    git remote add origin https://$repo_path

    repo_path=${repo_path//\//\\\/}

    sed -i -e "s/{ROBOT_NAME}/${robot_name}/g" ./Dockerfile

    sed -i -e "s/{REPO_PATH}/${repo_path}/" ./build.sh

    git add .
    git commit -m "init repo"
}

cmd_help(){
cat << EOF

Usage: $me robot-name dir-of-robot remote-repository-of-robot.

For Example: $me test . github.com/opensourceways

The command above will
generate codes at current dir with robot name of 'robot-gitee-test' and
import path of 'github.com/opensourceways/robot-gitee-test'.

EOF
}

if [ $pn -lt 3 ]; then
    cmd_help
    exit 1
fi

build $1 $2 $3
