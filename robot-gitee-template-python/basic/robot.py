import gitee
from robotgitee.client import Client
from robotgitee.framework.handlers import NoteEventHandler
from robotgitee.logutil import LogUtil
#from robotgitee.framework.handlers import PullRequestEventHandler
#from robotgitee.framework.handlers import IssueEventHandler
#from robotgitee.framework.handlers import PushEventHandler


class _Robot(NoteEventHandler):
    def __init__(self, cli: Client):
        self._cli = cli

    def handle_note_event(self, event: gitee.NoteEvent, log: LogUtil) -> None:
        log.info("receive note event")
