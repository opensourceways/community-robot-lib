# type handlers struct {
# 	issueHandlers      IssueHandler
# 	pullRequestHandler PullRequestHandler
# 	pushEventHandler   PushEventHandler
# 	noteEventHandler   NoteEventHandler
# }

class Handlers(object):
    def __init__(self, issueHandler=None, pullRequestHandler=None, pushEventHandler=None, noteEventHandler=None):
        self.issueHandler = issueHandler
        self.pullRequestHandler = pullRequestHandler
        self.pushEventHandler = pushEventHandler
        self.noteEventHandler = noteEventHandler

    def register_issue_handler(self, fn):
        pass

    def register_pull_request_handler(self, fn):
        pass

    def register_push_event_handler(self, fn):
        pass

    def register_note_event_handler(self, fn):
        pass
