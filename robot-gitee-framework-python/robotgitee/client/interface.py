from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import List, Optional

import gitee


@dataclass
class ListPullRequestOpt:
    state: Optional[str] = None
    head: Optional[str] = None
    base: Optional[str] = None
    sort: Optional[str] = None
    direction: Optional[str] = None
    milestone_number: Optional[int] = None
    labels: Optional[List[str]] = field(default_factory=list)


class Client(ABC):
    @abstractmethod
    def create_pull_request(self, org: str, repo: str, title: str, body, head, base: str) -> gitee.PullRequest:
        pass

    @abstractmethod
    def get_pull_requests(self, org: str, repo: str, opts: ListPullRequestOpt = None) -> List[gitee.PullRequest]:
        pass

    @abstractmethod
    def update_pull_request(self, org: str, repo: str, number: int,
                            body: gitee.PullRequestUpdateParam) -> gitee.PullRequest:
        pass

    @abstractmethod
    def list_collaborators(self, org: str, repo: str) -> List[gitee.ProjectMember]:
        pass

    @abstractmethod
    def is_collaborator(self, org: str, repo: str, username: str) -> bool:
        pass

    @abstractmethod
    def is_member(self, org: str, username: str) -> bool:
        pass

    @abstractmethod
    def remove_repo_member(self, org: str, repo: str, username: str):
        pass

    @abstractmethod
    def add_repo_member(self, org: str, repo: str, username: str, permission: str):
        pass

    @abstractmethod
    def get_ref(self, org: str, repo: str, ref: str) -> str:
        pass

    @abstractmethod
    def get_pull_request_changes(self, org: str, repo: str, number: int) -> List[gitee.PullRequestFiles]:
        pass

    @abstractmethod
    def get_pr_labels(self, org: str, repo: str, number: int) -> List[gitee.Label]:
        pass

    @abstractmethod
    def list_pr_comments(self, org: str, repo: str, number: int) -> List[gitee.PullRequestComments]:
        pass

    @abstractmethod
    def list_pr_issues(self, org: str, repo: str, number: int) -> List[gitee.Issue]:
        pass

    @abstractmethod
    def delete_pr_comment(self, org: str, repo: str, id: int):
        pass

    @abstractmethod
    def create_pr_comment(self, org: str, repo: str, number: int, comment: str):
        pass

    @abstractmethod
    def update_pr_comment(self, org: str, repo: str, comment_id: int, comment: str):
        pass

    @abstractmethod
    def add_pr_label(self, org: str, repo: str, number: int, label: str):
        pass

    @abstractmethod
    def add_multi_pr_label(self, org: str, repo: str, number: int, labels: List[str]):
        pass

    @abstractmethod
    def remove_pr_label(self, org: str, repo: str, number: int, label: str):
        pass

    @abstractmethod
    def remove_pr_labels(self, org: str, repo: str, number: int, labels: List[str]):
        pass

    @abstractmethod
    def replace_pr_all_labels(self, owner, repo: str, number: int, labels: List[str]):
        pass

    @abstractmethod
    def list_pr_operation_logs(self, org: str, repo: str, number: int) -> List[gitee.OperateLog]:
        pass

    @abstractmethod
    def close_pr(self, org: str, repo: str, number: int):
        pass

    @abstractmethod
    def assign_pr(self, owner, repo: str, number: int, usernames: List[str]):
        pass

    @abstractmethod
    def unassign_pr(self, owner, repo: str, number: int, usernames: List[str]):
        pass

    @abstractmethod
    def get_pr_commits(self, org: str, repo: str, number: int) -> List[gitee.PullRequestCommits]:
        pass

    @abstractmethod
    def get_pull_request(self, org: str, repo: str, number: int) -> gitee.PullRequest:
        pass

    @abstractmethod
    def get_repo_commit(self, org: str, repo: str, sha: str) -> gitee.RepoCommit:
        pass

    @abstractmethod
    def merge_pr(self, owner, repo: str, number: int, merge_method='merge', title=None, description=None,
                 prune_source_branch=False):
        pass

    @abstractmethod
    def get_org_repos(self, org: str) -> List[gitee.Project]:
        pass

    @abstractmethod
    def create_org_repo(self, org: str, repo: gitee.RepositoryPostParam):
        pass

    @abstractmethod
    def update_repo(self, org: str, repo: str, info: gitee.RepoPatchParam):
        pass

    @abstractmethod
    def get_repo(self, org: str, repo: str) -> gitee.Project:
        pass

    @abstractmethod
    def get_gitee_repo(self, org: str, repo: str) -> gitee.Project:
        pass

    @abstractmethod
    def set_repo_reviewer(self, org: str, repo: str, reviewer: gitee.SetRepoReviewer):
        pass

    @abstractmethod
    def create_repo_label(self, org: str, repo: str, label, color: str):
        pass

    @abstractmethod
    def get_repo_labels(self, owner, repo: str) -> List[gitee.Label]:
        pass

    @abstractmethod
    def assign_gitee_issue(self, org: str, repo: str, number: str, username: str):
        pass

    @abstractmethod
    def unassign_gitee_issue(self, org: str, repo: str, number: str):
        pass

    @abstractmethod
    def remove_issue_assignee(self, org: str, repo: str, number: str):
        pass

    @abstractmethod
    def create_issue_comment(self, org: str, repo: str, number: str, comment: str):
        pass

    @abstractmethod
    def update_issue_comment(self, org: str, repo: str, comment_id: int, comment: str):
        pass

    @abstractmethod
    def list_issue_comments(self, org: str, repo: str, number: str) -> List[gitee.Note]:
        pass

    @abstractmethod
    def get_issue_labels(self, org: str, repo: str, number: str) -> List[gitee.Label]:
        pass

    @abstractmethod
    def remove_issue_label(self, org: str, repo: str, number: str, label: str):
        pass

    @abstractmethod
    def remove_issue_labels(self, org: str, repo: str, number: str, labels: List[str]):
        pass

    @abstractmethod
    def add_issue_label(self, org: str, repo: str, number, label: str):
        pass

    @abstractmethod
    def add_multi_issue_label(self, org: str, repo: str, number: str, labels: List[str]):
        pass

    @abstractmethod
    def close_issue(self, owner, repo: str, number: str):
        pass

    @abstractmethod
    def reopen_issue(self, owner, repo: str, number: str):
        pass

    @abstractmethod
    def update_issue(self, owner, number: str, param: gitee.IssueUpdateParam) -> gitee.Issue:
        pass

    @abstractmethod
    def get_issue(self, org: str, repo: str, number: str) -> gitee.Issue:
        pass

    @abstractmethod
    def add_project_labels(self, org: str, repo: str, labels: List[str]):
        pass

    @abstractmethod
    def update_project_labels(self, org: str, repo: str, labels: List[str]):
        pass

    @abstractmethod
    def create_branch(self, org: str, repo: str, branch, parent_branch: str):
        pass

    @abstractmethod
    def get_repo_all_branch(self, org: str, repo: str) -> List[gitee.Branch]:
        pass

    @abstractmethod
    def set_protection_branch(self, org: str, repo: str, branch: str):
        pass

    @abstractmethod
    def cancel_protection_branch(self, org: str, repo: str, branch: str):
        pass

    @abstractmethod
    def create_file(self, org: str, repo: str, branch, path, content, commit_msg: str) -> gitee.CommitContent:
        pass

    @abstractmethod
    def get_path_content(self, org: str, repo: str, path, ref: str) -> gitee.Content:
        pass

    @abstractmethod
    def get_directory_tree(self, org: str, repo: str, sha: str, recursive: int) -> gitee.Tree:
        pass

    @abstractmethod
    def get_bot(self) -> gitee.User:
        pass

    @abstractmethod
    def get_user_permission_of_repo(self, org: str, repo: str, username: str) -> gitee.ProjectMemberPermission:
        pass

    @abstractmethod
    def get_enterprise_member(self, enterprise, username: str) -> gitee.EnterpriseMember:
        pass

    @abstractmethod
    def create_issue(self, org: str, repo: str, title, body: str) -> gitee.Issue:
        pass
