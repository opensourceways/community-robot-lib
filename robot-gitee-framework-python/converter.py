from typing import Dict

import dacite

from models.hook_event_types import NoteEvent, IssueEvent, PullRequestEvent, PushEvent


def convert_to_note_event(payload: Dict) -> NoteEvent:
    note_event = dacite.from_dict(NoteEvent, payload)
    return note_event


def convert_to_issue_event(payload: Dict) -> IssueEvent:
    issue_event = dacite.from_dict(IssueEvent, payload)
    return issue_event


def convert_to_pr_event(payload: Dict) -> PullRequestEvent:
    pull_request_event = dacite.from_dict(PullRequestEvent, payload)
    return pull_request_event


def convert_to_push_event(payload: Dict) -> PushEvent:
    push_event = dacite.from_dict(PushEvent, payload)
    return push_event
