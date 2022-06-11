from handlers import Handlers
from converter import *
from models.hook_event_helper import *


class Dispatcher(object):

    def __init__(self, handlers: Handlers = Handlers()):
        self.handlers = handlers

    def dispatch(self, event_type, payload):
        if event_type == EVENT_TYPE_NOTE:
            if not self.handlers.noteEventHandler:
                return None
            event = convert_to_note_event(payload)
        elif event_type == EVENT_TYPE_ISSUE:
            if not self.handlers.issueHandler:
                return None
            event = convert_to_issue_event(payload)
        elif event_type == EVENT_TYPE_PR:
            if not self.handlers.pullRequestHandler:
                return None
            event = convert_to_pr_event(payload)
        elif event_type == EVENT_TYPE_PUSH:
            if not self.handlers.pushEventHandler:
                return None
            event = convert_to_push_event(payload)
        else:
            pass

    def get_config(self):
        pass

    def handle_note_event(self, note_event: NoteEvent):
        self.handlers.noteEventHandler(note_event)

    def handle_pull_request_event(self, pull_request_event: PullRequestEvent):
        self.handlers.pullRequestHandler(pull_request_event)

    def handle_issue_event(self, issue_event: IssueEvent):
        self.handlers.issueHandler(issue_event)

    def handle_push_event(self, push_event: PushEvent):
        self.handlers.pushEventHandler(push_event)
