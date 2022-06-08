import threading
from flask import Flask, request
from handler.event_handler import PushHandler, TagHandler, NoteHandler, IssueHandler, MergeRequestHandler
from utils import parse_utils
from utils.parse_utils import parse_config, dict_to_object

MERGE_REQUEST_HOOK = 'Merge Request Hook'
PUSH_HOOK = 'Push Hook'
TAG_PUSH_HOOK = 'Tag Push Hook'
NOTE_HOOK = 'Note Hook'
ISSUE_HOOK = 'Issue Hook'

# supported hook types
HOOK_TYPE = [MERGE_REQUEST_HOOK, PUSH_HOOK, TAG_PUSH_HOOK, NOTE_HOOK, ISSUE_HOOK]
ACCESS_TOKEN = ''


def create_app():
    app = Flask(__name__)
    return app


def gitee_hook():
    headers = request.headers
    hook_type = parse_utils.parse_hook_headers(headers)
    print(hook_type)
    if hook_type not in HOOK_TYPE:
        print(f'{hook_type} not in allowed_hooks of your config file, end receiving webhook.')
        return 'webhook successfully'
    request_payload = request.json
    print(request_payload)
    request_payload['access_token'] = ACCESS_TOKEN
    payload_obj = dict_to_object(request_payload)
    if hook_type == PUSH_HOOK and GiteeRobot.PUSH_HANDLER is not None:
        GiteeRobot.PUSH_HANDLER.get_func()(payload_obj)
    elif hook_type == TAG_PUSH_HOOK and GiteeRobot.TAG_HANDLER is not None:
        GiteeRobot.TAG_HANDLER.get_func()(payload_obj)
    elif hook_type == MERGE_REQUEST_HOOK and GiteeRobot.MERGE_REQUEST_HANDLER is not None:
        GiteeRobot.MERGE_REQUEST_HANDLER.get_func()(payload_obj)
    elif hook_type == NOTE_HOOK and GiteeRobot.NOTE_HANDLER is not None:
        GiteeRobot.NOTE_HANDLER.get_func()(payload_obj)
    elif hook_type == ISSUE_HOOK and GiteeRobot.ISSUE_HANDLER is not None:
        GiteeRobot.ISSUE_HANDLER.get_func()(payload_obj)
    return 'webhook successfully'


class GiteeRobot(object):
    PUSH_HANDLER = None
    TAG_HANDLER = None
    NOTE_HANDLER = None
    ISSUE_HANDLER = None
    MERGE_REQUEST_HANDLER = None

    def __init__(self, config_file='../conf/conf.yaml'):
        self.config_file = config_file

    @staticmethod
    def add_handler(handler):
        handler.check_registered_func()
        if isinstance(handler, PushHandler):
            GiteeRobot.PUSH_HANDLER = handler
        elif isinstance(handler, TagHandler):
            GiteeRobot.TAG_HANDLER = handler
        elif isinstance(handler, NoteHandler):
            GiteeRobot.NOTE_HANDLER = handler
        elif isinstance(handler, IssueHandler):
            GiteeRobot.ISSUE_HANDLER = handler
        elif isinstance(handler, MergeRequestHandler):
            GiteeRobot.MERGE_REQUEST_HANDLER = handler
        else:
            raise ValueError('No handler of the correct type.')

    @staticmethod
    def check_handler():
        handler_list = [GiteeRobot.PUSH_HANDLER, GiteeRobot.TAG_HANDLER, GiteeRobot.NOTE_HANDLER,
                        GiteeRobot.ISSUE_HANDLER, GiteeRobot.MERGE_REQUEST_HANDLER]
        if not any(handler_list):
            raise Exception('No registered handler found.')
        return True

    def load_config(self):
        global HOOK_TYPE
        global ACCESS_TOKEN
        options = parse_config(self.config_file)
        access_token = options.get('access_token')
        if not access_token:
            raise Exception('access_token in config file cannot be empty')
        ACCESS_TOKEN = access_token
        allowed_hooks = options.get('allowed_hooks', HOOK_TYPE)
        if allowed_hooks:
            if not isinstance(allowed_hooks, list):
                raise Exception('allowed_hooks in config file need a list type value.')
            unexpected_hooks = set(allowed_hooks) - set(HOOK_TYPE)
            if unexpected_hooks:
                raise Exception(
                    'allowed_hooks in config file only support Merge Request Hook, Push Hook, Tag Push Hook, '
                    'Note Hook and Issue Hook.')
        HOOK_TYPE = allowed_hooks
        return options

    def run(self):
        options = self.load_config()
        ip = options.get('ip', '127.0.0.1')
        port = options.get('port', 5000)
        route = options.get('route', '/gitee-hook')
        debug = options.get('debug', False)
        app = create_app()
        app.add_url_rule(route, view_func=gitee_hook, methods=['POST'])
        threading.Thread(target=app.run, args=(ip, port, debug)).start()
