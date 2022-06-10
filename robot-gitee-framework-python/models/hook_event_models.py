from dataclasses import dataclass
from typing import List, Optional
from model_issue import Issue


@dataclass
class LabelHook:
    id: Optional[int]
    name: Optional[str]
    color: Optional[str]


@dataclass
class UserHook:
    id: Optional[int]
    name: Optional[str]
    email: Optional[str]
    user_name: Optional[str]
    url: Optional[str]
    login: Optional[str]
    avatar_url: Optional[str]
    html_url: Optional[str]
    type: Optional[str]
    site_admin: Optional[str]
    time: Optional[str]
    remark: Optional[str]


# EnterpriseHook : 企业信息
@dataclass
class EnterpriseHook:
    name: Optional[str]
    url: Optional[str]


# NoteHook : 评论信息
@dataclass
class NoteHook:
    id: Optional[int]
    body: Optional[str]
    user: Optional[UserHook]
    created_at: Optional[str]
    updated_at: Optional[str]
    html_url: Optional[str]
    position: Optional[str]
    commit_id: Optional[str]


# ProjectHook : project 信息
@dataclass
class ProjectHook:
    id: Optional[int]
    name: Optional[str]
    path: Optional[str]
    full_name: Optional[str]
    owner: Optional[UserHook]
    private: Optional[bool]
    html_url: Optional[str]
    url: Optional[str]
    description: Optional[str]
    fork: Optional[bool]
    pushed_at: Optional[str]
    created_at: Optional[str]
    updated_at: Optional[str]
    ssh_url: Optional[str]
    git_url: Optional[str]
    clone_url: Optional[str]
    svn_url: Optional[str]
    git_http_url: Optional[str]
    git_ssh_url: Optional[str]
    git_svn_url: Optional[str]
    homepage: Optional[str]
    stargazers_count: Optional[int]
    watchers_count: Optional[int]
    forks_count: Optional[int]
    language: Optional[str]

    has_issues: Optional[bool]
    has_wiki: Optional[bool]
    has_pages: Optional[bool]
    license: Optional[str]

    open_issues_count: Optional[int]
    default_branch: Optional[str]
    namespace: Optional[str]

    name_with_namespace: Optional[str]
    path_with_namespace: Optional[str]


# MilestoneHook : 里程碑信息
@dataclass
class MilestoneHook:
    id: Optional[int]
    html_url: Optional[str]
    number: Optional[int]
    title: Optional[str]
    description: Optional[str]
    open_issues: Optional[int]
    closed_issues: Optional[int]
    state: Optional[str]
    created_at: Optional[str]
    updated_at: Optional[str]
    due_on: Optional[str]


# BranchHook : 分支信息
@dataclass
class BranchHook:
    label: Optional[str]
    ref: Optional[str]
    sha: Optional[str]
    user: Optional[UserHook]
    repo: Optional[ProjectHook]


# PullRequestHook : PR 信息
@dataclass
class PullRequestHook:
    id: Optional[int]
    number: Optional[int]
    state: Optional[str]
    html_url: Optional[str]
    diff_url: Optional[str]
    patch_url: Optional[str]
    title: Optional[str]
    body: Optional[str]
    stale_labels: Optional[List[LabelHook]]
    labels: Optional[List[LabelHook]]
    created_at: Optional[str]
    updated_at: Optional[str]
    closed_at: Optional[str]
    merged_at: Optional[str]
    merge_commit_sha: Optional[str]
    merge_reference_name: Optional[str]
    user: Optional[UserHook]
    assignee: Optional[UserHook]
    assignees: Optional[List[UserHook]]
    tester: Optional[List[UserHook]]
    testers: Optional[List[UserHook]]
    need_test: Optional[bool]
    need_review: Optional[bool]
    milestone: Optional[MilestoneHook]
    head: Optional[BranchHook]
    base: Optional[BranchHook]
    merged: Optional[bool]
    mergeable: Optional[bool]
    merge_status: Optional[str]
    updated_by: Optional[UserHook]
    comments: Optional[int]
    commits: Optional[int]
    additions: Optional[int]
    deletions: Optional[int]
    changed_files: Optional[int]
    issues: Optional[List[Issue]]
    stale_issues: Optional[List[Issue]]


# IssueHook : issue 信息
@dataclass
class IssueHook:
    id: Optional[int]
    html_url: Optional[str]
    number: Optional[str]
    title: Optional[str]
    user: Optional[UserHook]
    labels: Optional[List[LabelHook]]
    state: Optional[str]
    state_name: Optional[str]
    type_name: Optional[str]
    assignee: Optional[UserHook]
    collaborators: Optional[List[UserHook]]
    milestone: Optional[MilestoneHook]
    comments: Optional[int]
    created_at: Optional[str]
    updated_at: Optional[str]
    body: Optional[str]


# CommitHook : git commit 中的信息
@dataclass
class CommitHook:
    id: Optional[str]
    tree_id: Optional[str]
    parent_ids: Optional[List[str]]
    message: Optional[str]
    timestamp: Optional[str]
    url: Optional[str]
    author: Optional[UserHook]
    committer: Optional[UserHook]
    distinct: Optional[bool]
    added: Optional[List[str]]
    removed: Optional[List[str]]
    modified: Optional[List[str]]
