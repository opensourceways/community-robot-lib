class Handlers(object):
    def __init__(self):
        self.issueHandler = None
        self.pullRequestHandler = None
        self.pushEventHandler = None
        self.noteEventHandler = None

    def register_issue_handler(self, fn):
        self.issueHandler = fn

    def register_pull_request_handler(self, fn):
        self.pullRequestHandler = fn

    def register_push_event_handler(self, fn):
        self.pushEventHandler = fn

    def register_note_event_handler(self, fn):
        self.noteEventHandler = fn
