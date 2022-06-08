from dataclasses import dataclass
from typing import List
from model_issue import Issue


@dataclass
class LabelHook:
    id: int
    name: str
    color: str


@dataclass
class UserHook:
    id: int
    name: str
    email: str
    user_name: str
    url: str
    login: str
    avatar_url: str
    html_url: str
    type: str
    site_admin: str
    time: str
    remark: str


# EnterpriseHook : 企业信息
@dataclass
class EnterpriseHook:
    name: str
    url: str


# NoteHook : 评论信息
@dataclass
class NoteHook:
    id: int
    body: str
    user: UserHook
    created_at: str
    updated_at: str
    html_url: str
    position: str
    commit_id: str


# ProjectHook : project 信息
@dataclass
class ProjectHook:
    id: int
    name: str
    path: str
    full_name: str
    owner: UserHook
    private: bool
    html_url: str
    url: str
    description: str
    fork: bool
    pushed_at: str
    created_at: str
    updated_at: str
    ssh_url: str
    git_url: str
    clone_url: str
    svn_url: str
    git_http_url: str
    git_ssh_url: str
    git_svn_url: str
    homepage: str
    stargazers_count: int
    watchers_count: int
    forks_count: int
    language: str

    has_issues: bool
    has_wiki: bool
    has_pages: bool
    license: str

    open_issues_count: int
    default_branch: str
    namespace: str

    name_with_namespace: str
    path_with_namespace: str


# MilestoneHook : 里程碑信息
@dataclass
class MilestoneHook:
    id: int
    html_url: str
    number: int
    title: str
    description: str
    open_issues: int
    closed_issues: int
    state: str
    created_at: str
    updated_at: str
    due_on: str


# BranchHook : 分支信息
@dataclass
class BranchHook:
    label: str
    ref: str
    sha: str
    user: UserHook
    repo: ProjectHook


# PullRequestHook : PR 信息
@dataclass
class PullRequestHook:
    id: int
    number: int
    state: str
    html_url: str
    diff_url: str
    patch_url: str
    title: str
    body: str
    stale_labels: List[LabelHook]
    labels: List[LabelHook]
    created_at: str
    updated_at: str
    closed_at: str
    merged_at: str
    merge_commit_sha: str
    merge_reference_name: str
    user: UserHook
    assignee: UserHook
    assignees: List[UserHook]
    tester: List[UserHook]
    testers: List[UserHook]
    need_test: bool
    need_review: bool
    milestone: MilestoneHook
    head: BranchHook
    base: BranchHook
    merged: bool
    mergeable: bool
    merge_status: str
    updated_by: UserHook
    comments: int
    commits: int
    additions: int
    deletions: int
    changed_files: int
    issues: List[Issue]
    stale_issues: List[Issue]


# IssueHook : issue 信息
@dataclass
class IssueHook:
    id: int
    html_url: str
    number: str
    title: str
    user: UserHook
    labels: List[LabelHook]
    state: str
    state_name: str
    type_name: str
    assignee: UserHook
    collaborators: List[UserHook]
    milestone: MilestoneHook
    comments: int
    created_at: str
    updated_at: str
    body: str


# CommitHook : git commit 中的信息
@dataclass
class CommitHook:
    id: str
    tree_id: str
    parent_ids: List[str]
    message: str
    timestamp: str
    url: str
    author: UserHook
    committer: UserHook
    distinct: bool
    added: List[str]
    removed: List[str]
    modified: List[str]
