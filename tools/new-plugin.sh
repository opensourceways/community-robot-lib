#!/usr/bin/bash

set -euo pipefail

cd $(dirname $0)
me=$(basename $0)
pn=$#

build(){
    local repo_dir=$1
    local repo_path=$2
    local repo_name=$(basename $repo_path)

    if [ -d $repo_dir ]; then
	    echo "$repo_dir is exist"
	    return 1
    fi

    mkdir $repo_dir
    cd $repo_dir

    git clone https://github.com/opensourceways/community-robot-lib.git

    cp -r community-robot-lib/new-plugin/. .

    rm -fr community-robot-lib 

    repo_path=${repo_path//\//\\\/}

    sed -i -e "s/{PLUGIN_REPO}/${repo_path}/" ./BUILD.bazel
    sed -i -e "s/{PLUGIN_NAME}/${repo_name}/" ./BUILD.bazel

    git init .
    git add .
    git commit -m "init repo"
}

cmd_help(){
cat << EOF
usage: $me absolute-dir-of-repo import-path-of-repo. for example: $me /test github.com/exmaple/test 
EOF
}

if [ $pn -lt 2 ]; then
    cmd_help
    exit 1
fi

build $1 $2
