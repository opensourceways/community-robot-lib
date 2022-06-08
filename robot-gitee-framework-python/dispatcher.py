from handlers import Handlers
from hook_event_helper import *


class Dispatcher(object):

    def __init__(self, handlers: Handlers = Handlers()):
        self.handlers = handlers

    def dispatch(self, event_type, payload):
        if event_type == EVENT_TYPE_NOTE:
            if not self.handlers.noteEventHandler:
                return None
            pass
        elif event_type == EVENT_TYPE_ISSUE:
            if not self.handlers.issueHandler:
                return None
            pass
        elif event_type == EVENT_TYPE_PR:
            if not self.handlers.pullRequestHandler:
                return None
            pass
        elif event_type == EVENT_TYPE_PUSH:
            if not self.handlers.pushEventHandler:
                return None
            pass
        else:
            pass

    def handleNoteEvent(self, note_event):
        pass

    def handlePullRequestEvent(self, PullRequestEvent):
        pass
