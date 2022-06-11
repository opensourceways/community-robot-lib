from dataclasses import dataclass
from typing import Optional

import dacite
from dacite import Config

from hook_event_models import *


@dataclass
class NoteEvent(object):
    action: Optional[str]
    comment: Optional[NoteHook]
    repository: Optional[ProjectHook]
    project: Optional[ProjectHook]
    author: Optional[UserHook]
    sender: Optional[UserHook]
    url: Optional[str]
    note: Optional[str]
    noteable_type: Optional[str]
    noteable_id: Optional[str]
    title: Optional[str]
    per_iid: Optional[str]
    short_commit_id: Optional[str]
    enterprise: Optional[EnterpriseHook]
    pull_request: Optional[PullRequestHook]
    issue: Optional[IssueHook]
    hook_name: Optional[str]
    password: Optional[str]


@dataclass
class IssueEvent:
    action: Optional[str]
    issue: Optional[IssueHook]
    repository: Optional[ProjectHook]
    project: Optional[ProjectHook]
    sender: Optional[UserHook]
    target_user: Optional[UserHook]
    user: Optional[UserHook]
    assignee: Optional[UserHook]
    updated_by: Optional[UserHook]
    iid: Optional[str]
    title: Optional[str]
    description: Optional[str]
    state: Optional[str]
    milestone: Optional[str]
    url: Optional[str]
    enterprise: Optional[EnterpriseHook]
    hook_name: Optional[str]
    password: Optional[str]


@dataclass
class RepoInfo:
    project: Optional[ProjectHook]
    repository: Optional[ProjectHook]


@dataclass
class PullRequestEvent:
    action: Optional[str]
    action_desc: Optional[str]
    pull_request: Optional[PullRequestHook]
    number: Optional[int]
    iid: Optional[int]
    title: Optional[str]
    body: Optional[str]
    state: Optional[str]
    merge_status: Optional[str]
    merge_commit_sha: Optional[str]
    url: Optional[str]
    source_branch: Optional[str]
    source_repo: Optional[RepoInfo]
    target_branch: Optional[str]
    target_repo: Optional[RepoInfo]
    project: Optional[ProjectHook]
    repository: Optional[ProjectHook]
    author: Optional[UserHook]
    updated_by: Optional[UserHook]
    sender: Optional[UserHook]
    target_user: Optional[UserHook]
    enterprise: Optional[EnterpriseHook]
    hook_name: Optional[str]
    password: Optional[str]


@dataclass
class PushEvent:
    ref: Optional[str]
    before: Optional[str]
    after: Optional[str]
    total_commits_count: Optional[int]
    commits_more_than_ten: Optional[bool]
    created: Optional[bool]
    deleted: Optional[bool]
    compare: Optional[str]
    commits: Optional[List[CommitHook]]
    head_commit: Optional[CommitHook]
    repository: Optional[ProjectHook]
    project: Optional[ProjectHook]
    user_id: Optional[int]
    user_name: Optional[str]
    user: Optional[UserHook]
    pusher: Optional[UserHook]
    sender: Optional[UserHook]
    enterprise: Optional[EnterpriseHook]
    hook_name: Optional[str]
    password: Optional[str]

    def get_ref(self) -> str:
        if self is None or self.ref is None:
            return ""
        return self.ref

    def get_before(self) -> str:
        if self is None or self.before is None:
            return ""
        return self.before

    def get_after(self) -> str:
        if self is None or self.after is None:
            return ""
        return self.after

    def get_commits_more_than_ten(self) -> bool:
        if self is None or self.commits_more_than_ten is None:
            return False
        return self.commits_more_than_ten

    def get_created(self) -> bool:
        if self is None or self.created is None:
            return False
        return self.created


@dataclass
class TagPushEvent:
    action: Optional[str]


if __name__ == '__main__':
    data = {'ref': 'ref', 'before': 'before', 'after': 'after', 'total_commits_count': 10,
            'commits_more_than_ten': False, 'created': True, 'deleted': False, 'compare': 'compare', 'commits': [],
            'head_commit': None, 'repository': None, 'project': None, 'user_id': 1, 'user_name': 'wwl', 'user': None,
            'pusher': None, 'sender': None, 'enterprise': None, 'hook_name': 'push', 'password': '123456'}
    event = dacite.from_dict(PushEvent, data)
    print(event.get_ref())
    print(type(event))
    # dacite.from_dict()
    print(event)
