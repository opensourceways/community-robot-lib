import json
import os

import yaml


class Dict(dict):
    __setattr__ = dict.__setitem__
    __getattr__ = dict.__getitem__


def dict_to_object(dict_obj):
    # nested dict type supported
    if isinstance(dict_obj, list):
        instance_list = []
        for i in dict_obj:
            instance_list.append(dict_to_object(i))
        return instance_list

    if not isinstance(dict_obj, dict):
        return dict_obj
    instance = Dict()
    for k, v in dict_obj.items():
        instance[k] = dict_to_object(v)
    return instance


def json_to_object(json_str: str):
    return dict_to_object(json.loads(json_str))


def parse_hook_headers(hook_headers):
    hook_type = hook_headers.get('X-Gitee-Event')
    return hook_type


def parse_config(config_file):
    if not os.path.exists(config_file):
        raise Exception('configuration file does not exist.')
    with open(config_file, 'r') as config_file:
        config_options = yaml.load(config_file, Loader=yaml.SafeLoader)
    return config_options


if __name__ == '__main__':
    data = {'a': 1, 'b': 2, 'c': ['1']}
    dict_to_object()
