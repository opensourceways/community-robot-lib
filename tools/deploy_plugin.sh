#!/bin/bash

set -euo pipefail

cd $(dirname $0)
me=$(basename $0)
pn=$#

ph_bot_name="{{BOT-NAME}}"
ph_platform="{{PLATFORM}}"
ph_image="{{IMAGE}}"
ph_port="{{PORT}}"

replace() {
    local from=$1
    local to=$2
    local file=$3

    sed -i -e "s/${from}/${to}/g" $file
}

underscore_to_hyphen(){
    local name=$1
    echo ${name//_/-}
}

convert_backslash(){
    local s=$1
    echo ${s//\//\\\/}
}

gen_deployment() {
    local bot_name=$1
    local platform=$2
    local image=$3
    local port=$4

    bot_name=$(underscore_to_hyphen $bot_name)
    platform=$(underscore_to_hyphen $platform)
    image=$(convert_backslash $image)

    local file=deployment.yaml

    replace "$ph_bot_name" "$bot_name" $file
    replace "$ph_platform" "$platform" $file
    replace "$ph_image" "$image" $file
    replace "$ph_port" "$port" $file
}


gen_service() {
    local bot_name=$1
    local port=$2

    bot_name=$(underscore_to_hyphen $bot_name)

    local file=service.yaml

    replace "$ph_bot_name" "$bot_name" $file
    replace "$ph_port" "$port" $file
}

gen_app_config() {
    local bot_name=$1
    local platform=$2

    bot_name=$(underscore_to_hyphen $bot_name)
    platform=$(underscore_to_hyphen $platform)

    local file=app_config.yaml

    replace "$ph_bot_name" "$bot_name" $file
    replace "$ph_platform" "$platform" $file
}

cmd_help(){
cat << EOF
usage: $me dir-of-repo platform bot_name port image.
for example: $me ./checkpr gitee checkpr 8888 swr.ap-southeast-1.myhuaweicloud.com/opensourceway/robot/robot-gitee-checkpr:bc480df
EOF
}

deploy(){
    if [ $pn -lt 5 ]; then
        cmd_help
        exit 1
    fi

    local repo_dir=$1
    local platform=$2
    local bot_name=$3
    local port=$4
    local image=$5

    if [ -d $repo_dir ]; then
        echo "$repo_dir is exist"
        return 1
    fi

    mkdir $repo_dir
    cd $repo_dir

    git clone https://github.com/opensourceways/community-robot-lib.git

    cp -r community-robot-lib/deploy-plugin/. .

    rm -fr community-robot-lib

    gen_service $bot_name $port

    gen_deployment $bot_name $platform $image $port

    gen_app_config $bot_name $platform
}

deploy $@
