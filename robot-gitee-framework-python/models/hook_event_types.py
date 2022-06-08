from dataclasses import dataclass
from hook_event_models import *


@dataclass
class NoteEvent(object):
    action: str
    comment: NoteHook
    repository: ProjectHook
    project: ProjectHook
    author: UserHook
    sender: UserHook
    url: str
    note: str
    noteable_type: str
    noteable_id: str
    title: str
    per_iid: str
    short_commit_id: str
    enterprise: EnterpriseHook
    pull_request: PullRequestHook
    issue: IssueHook
    hook_name: str
    password: str


@dataclass
class IssueEvent:
    action: str
    issue: IssueHook
    repository: ProjectHook
    project: ProjectHook
    sender: UserHook
    target_user: UserHook
    user: UserHook
    assignee: UserHook
    updated_by: UserHook
    iid: str
    title: str
    description: str
    state: str
    milestone: str
    url: str
    enterprise: EnterpriseHook
    hook_name: str
    password: str


@dataclass
class RepoInfo:
    project: ProjectHook
    repository: ProjectHook


@dataclass
class PullRequestEvent:
    action: str
    action_desc: str
    pull_request: PullRequestHook
    number: int
    iid: int
    title: str
    body: str
    state: str
    merge_status: str
    merge_commit_sha: str
    url: str
    source_branch: str
    source_repo: RepoInfo
    target_branch: str
    target_repo: RepoInfo
    project: ProjectHook
    repository: ProjectHook
    author: UserHook
    updated_by: UserHook
    sender: UserHook
    target_user: UserHook
    enterprise: EnterpriseHook
    hook_name: str
    password: str


@dataclass
class PushEvent:
    ref: str
    before: str
    after: str
    total_commits_count: int
    commits_more_than_ten: bool
    created: bool
    deleted: bool
    compare: str
    commits: List[CommitHook]
    head_commit: CommitHook
    repository: ProjectHook
    project: ProjectHook
    user_id: int
    user_name: str
    user: UserHook
    pusher: UserHook
    sender: UserHook
    enterprise: EnterpriseHook
    hook_name: str
    password: str


@dataclass
class TagPushEvent:
    action: str
