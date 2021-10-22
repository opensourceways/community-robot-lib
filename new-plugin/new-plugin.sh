#!/usr/bin/bash

set -euo pipefail

cd $(dirname $0)
me=$(basename $0)
pn=$#

build(){
    local repo_path=$1
    local repo_name=$(basename $repo_path)

    if [ -d $repo_path ]; then
	    echo "$repo_path is exist"
	    return 1
    fi

    mkdir $repo_path
    cd $repo_path

    git clone https://github.com/opensourceways/community-robot-lib.git

    cp -r community-robot-lib/new-plugin/. .

    rm -fr community-robot-lib 

    git init .
    git add .
    git commit -m "init repo"
}

cmd_help(){
cat << EOF
usage: $me repo-absolute-dir
EOF
}

if [ $pn -lt 1 ]; then
    cmd_help
    exit 1
fi

build $1

