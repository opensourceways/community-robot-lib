import base64
from typing import List

import gitee
from robotgitee.client.interface import Client, ListPullRequestOpt


class _Client(Client):
    def __init__(self, token: str):
        cli = gitee.ApiClient(header_name="Authorization", header_value="Bearer " + token)

        self.pr_api = gitee.PullRequestsApi(cli)
        self.repository_api = gitee.RepositoriesApi(cli)
        self.organization_api = gitee.OrganizationsApi(cli)
        self.issue_api = gitee.IssuesApi(cli)
        self.label_api = gitee.LabelsApi(cli)
        self.git_data_api = gitee.GitDataApi(cli)
        self.user_api = gitee.UsersApi(cli)
        self.enterprise_api = gitee.EnterprisesApi(cli)

    def create_pull_request(self, org: str, repo: str, title: str, body, head, base: str):
        body = gitee.CreatePullRequestParam(title=title, body=body, base=base, head=head, prune_source_branch=True)
        return self.pr_api.post_v5_repos_owner_repo_pulls(org, repo, body)

    def get_pull_requests(self, org: str, repo: str, opts: ListPullRequestOpt = None) -> List[gitee.PullRequest]:
        if opts is not None:
            param = vars(opts)
            return self.pr_api.get_v5_repos_owner_repo_pulls(org, repo, **param)
        else:
            return self.pr_api.get_v5_repos_owner_repo_pulls(org, repo)

    def update_pull_request(self, org: str, repo: str, number: int,
                            body: gitee.PullRequestUpdateParam) -> gitee.PullRequest:
        return self.pr_api.patch_v5_repos_owner_repo_pulls_number(org, repo, number, body)

    def list_collaborators(self, org: str, repo: str) -> List[gitee.ProjectMember]:
        return self.repository_api.get_v5_repos_owner_repo_collaborators(org, repo)

    def is_collaborator(self, org: str, repo: str, username: str) -> bool:
        try:
            self.repository_api.get_v5_repos_owner_repo_collaborators_username(org, repo, username)
        except Exception as e:
            return False
        return True

    def is_member(self, org: str, username: str) -> bool:
        try:
            self.organization_api.get_v5_orgs_org_memberships_username(org, username)
        except Exception as e:
            return False
        return True

    def remove_repo_member(self, org: str, repo: str, username: str):
        return self.repository_api.delete_v5_repos_owner_repo_collaborators_username(org, repo, username)

    def add_repo_member(self, org: str, repo: str, username: str, permission: str):
        body = gitee.ProjectMemberPutParam(permission=permission)
        self.repository_api.put_v5_repos_owner_repo_collaborators_username(org, repo, username, body)

    def get_ref(self, org: str, repo: str, ref: str) -> str:
        branch = self.repository_api.get_v5_repos_owner_repo_branches_branch(org, repo, ref)
        return branch.commit.sha

    def get_pull_request_changes(self, org: str, repo: str, number: int) -> List[gitee.PullRequestFiles]:
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_files(org, repo, number)

    def get_pr_labels(self, org: str, repo: str, number: int) -> List[gitee.Label]:
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_labels(org, repo, number)

    def list_pr_comments(self, org: str, repo: str, number: int) -> List[gitee.PullRequestComments]:
        comments = self.pr_api.get_v5_repos_owner_repo_pulls_number_comments(org, repo, number)
        return comments

    def list_pr_issues(self, org: str, repo: str, number: int) -> List[gitee.Issue]:
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_issues(org, repo, number)

    def delete_pr_comment(self, org: str, repo: str, id: int):
        self.pr_api.delete_v5_repos_owner_repo_pulls_comments_id(org, repo, id)

    def create_pr_comment(self, org: str, repo: str, number: int, comment: str):
        body = gitee.PullRequestCommentPostParam(body=comment)
        return self.pr_api.post_v5_repos_owner_repo_pulls_number_comments(org, repo, number, body)

    def update_pr_comment(self, org: str, repo: str, comment_id: int, comment: str):
        body = gitee.PullRequestCommentPatchParam(body=comment)
        return self.pr_api.patch_v5_repos_owner_repo_pulls_comments_id(org, repo, comment_id, body)

    def add_multi_pr_label(self, org: str, repo: str, number: int, labels: List[str]):
        return self.pr_api.post_v5_repos_owner_repo_pulls_number_labels(org, repo, number, labels)

    def add_pr_label(self, org: str, repo: str, number: int, label: str):
        labels = [label]
        return self.add_multi_pr_label(org, repo, number, labels)

    def remove_pr_label(self, org: str, repo: str, number: int, label: str):
        return self.pr_api.delete_v5_repos_owner_repo_pulls_label(org, repo, number, label)

    def remove_pr_labels(self, org: str, repo: str, number: int, labels: List[str]):
        pr_labels = self.get_pr_labels(org, repo, number)
        label_list = []
        for i in pr_labels:
            if isinstance(i, gitee.Label):
                label_list.append(i.name)
        error_labels = list(set(labels) - set(label_list))
        if error_labels:
            labels_str = ", ".join(error_labels)
            print(f'{labels_str} labels does not exist in pr. cancel delete labels.')
            return
        else:
            for i in label_list:
                return self.pr_api.delete_v5_repos_owner_repo_pulls_label(org, repo, number, i)

    def replace_pr_all_labels(self, owner, repo: str, number: int, labels: List[str]):
        return self.pr_api.put_v5_repos_owner_repo_pulls_number_labels(owner, repo, number, labels)

    def list_pr_operation_logs(self, org: str, repo: str, number: int) -> List[gitee.OperateLog]:
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_operate_logs(org, repo, number)

    def close_pr(self, org: str, repo: str, number: int):
        body = gitee.PullRequestUpdateParam(state='closed')
        return self.update_pull_request(org, repo, number, body)

    def assign_pr(self, owner, repo: str, number: int, usernames: List[str]):
        assignees = ','.join(usernames)
        body = gitee.PullRequestAssigneePostParam(assignees=assignees)
        return self.pr_api.post_v5_repos_owner_repo_pulls_number_assignees(owner, repo, number, body)

    def unassign_pr(self, owner, repo: str, number: int, usernames: List[str]):
        assignees = ','.join(usernames)
        return self.pr_api.delete_v5_repos_owner_repo_pulls_number_assignees(owner, repo, number, assignees)

    def get_pr_commits(self, org: str, repo: str, number: int) -> List[gitee.PullRequestCommits]:
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_commits(org, repo, number)

    def get_pull_request(self, org: str, repo: str, number: int) -> gitee.PullRequest:
        return self.pr_api.get_v5_repos_owner_repo_pulls_number(org, repo, number)

    def get_repo_commit(self, org: str, repo: str, sha: str) -> gitee.RepoCommit:
        return self.repository_api.get_v5_repos_owner_repo_commits_sha(org, repo, sha)

    def merge_pr(self, owner, repo: str, number: int, merge_method='merge', title=None, description=None
                 , prune_source_branch=False):
        body = gitee.PullRequestMergePutParam(merge_method=merge_method, prune_source_branch=prune_source_branch,
                                              title=title, description=description)
        self.pr_api.put_v5_repos_owner_repo_pulls_number_merge(owner, repo, number, body)

    def get_org_repos(self, org: str) -> List[gitee.Project]:
        return self.repository_api.get_v5_orgs_org_repos(org)

    def create_org_repo(self, org: str, repo: gitee.RepositoryPostParam):
        return self.repository_api.post_v5_orgs_org_repos(org, repo)

    def update_repo(self, org: str, repo: str, info: gitee.RepoPatchParam):
        return self.repository_api.patch_v5_repos_owner_repo(org, repo, info)

    def get_repo(self, org: str, repo: str) -> gitee.Project:
        return self.repository_api.get_v5_repos_owner_repo(org, repo)

    def get_gitee_repo(self, org: str, repo: str) -> gitee.Project:
        return self.repository_api.get_v5_repos_owner_repo(org, repo)

    def set_repo_reviewer(self, org: str, repo: str, reviewer: gitee.SetRepoReviewer):
        return self.repository_api.put_v5_repos_owner_repo_reviewer(org, repo, reviewer)

    def create_repo_label(self, org: str, repo: str, label, color: str):
        body = gitee.LabelPostParam(name=label, color=color)
        return self.label_api.post_v5_repos_owner_repo_labels(org, repo, body)

    def get_repo_labels(self, owner, repo: str) -> List[gitee.Label]:
        return self.label_api.get_v5_repos_owner_repo_labels(owner, repo)

    def assign_gitee_issue(self, org: str, repo: str, number: str, username: str):
        body = gitee.IssueUpdateParam(repo=repo, assignee=username)
        return self.issue_api.patch_v5_repos_owner_issues_number(org, number, body)

    def unassign_gitee_issue(self, org: str, repo: str, number: str):
        return self.assign_gitee_issue(org, repo, number, " ")

    def remove_issue_assignee(self, org: str, repo: str, number: str):
        return self.assign_gitee_issue(org, repo, number, " ")

    def create_issue_comment(self, org: str, repo: str, number: str, comment: str):
        body = gitee.IssueCommentPostParam(body=comment)
        return self.issue_api.post_v5_repos_owner_repo_issues_number_comments(org, repo, number, body)

    def update_issue_comment(self, org: str, repo: str, comment_id: int, comment: str):
        body = gitee.IssueCommentPatchParam(body=comment)
        return self.issue_api.patch_v5_repos_owner_repo_issues_comments_id(org, repo, comment_id, body)

    def list_issue_comments(self, org: str, repo: str, number: str) -> List[gitee.Note]:
        return self.issue_api.get_v5_repos_owner_repo_issues_number_comments(org, repo, number)

    def get_issue_labels(self, org: str, repo: str, number: str) -> List[gitee.Label]:
        return self.label_api.get_v5_repos_owner_repo_issues_number_labels(org, repo, number)

    def remove_issue_label(self, org: str, repo: str, number: str, label: str):
        return self.label_api.delete_v5_repos_owner_repo_issues_number_labels_name(org, repo, number, label)

    def remove_issue_labels(self, org: str, repo: str, number: str, labels: List[str]):
        issue_labels = self.get_issue_labels(org, repo, number)
        label_list = []
        for i in issue_labels:
            if isinstance(i, gitee.Label):
                label_list.append(i.name)
        error_labels = list(set(labels) - set(label_list))
        if error_labels:
            labels_str = ", ".join(error_labels)
            print(f'{labels_str} labels does not exist in issue. cancel delete labels.')
            return
        else:
            for i in label_list:
                return self.label_api.delete_v5_repos_owner_repo_issues_number_labels_name(org, repo, number, i)

    def add_issue_label(self, org: str, repo: str, number, label: str):
        labels = [label]
        return self.label_api.post_v5_repos_owner_repo_issues_number_labels(org, repo, number, labels)

    def add_multi_issue_label(self, org: str, repo: str, number: str, labels: List[str]):
        return self.label_api.post_v5_repos_owner_repo_issues_number_labels(org, repo, number, labels)

    def update_issue(self, owner, number: str, param: gitee.IssueUpdateParam) -> gitee.Issue:
        return self.issue_api.patch_v5_repos_owner_issues_number(owner, number, param)

    def close_issue(self, owner, repo: str, number: str):
        param = gitee.IssueUpdateParam(repo=repo, state='closed')
        return self.update_issue(owner, number, param)

    def reopen_issue(self, owner, repo: str, number: str):
        param = gitee.IssueUpdateParam(repo=repo, state='open')
        return self.update_issue(owner, number, param)

    def get_issue(self, org: str, repo: str, number: str) -> gitee.Issue:
        return self.issue_api.get_v5_repos_owner_repo_issues_number(org, repo, number)

    def add_project_labels(self, org: str, repo: str, labels: List[str]):
        return self.label_api.post_v5_repos_owner_repo_project_labels(org, repo, labels)

    def update_project_labels(self, org: str, repo: str, labels: List[str]):
        return self.label_api.put_v5_repos_owner_repo_project_labels(org, repo, labels)

    def create_branch(self, org: str, repo: str, branch: str, parent_branch: str):
        body = gitee.CreateBranchParam(branch_name=branch, refs=parent_branch)
        return self.repository_api.post_v5_repos_owner_repo_branches(org, repo, body)

    def get_repo_all_branch(self, org: str, repo: str) -> List[gitee.Branch]:
        return self.repository_api.get_v5_repos_owner_repo_branches(org, repo)

    def set_protection_branch(self, org: str, repo: str, branch: str):
        body = gitee.BranchProtectionPutParam()
        return self.repository_api.put_v5_repos_owner_repo_branches_branch_protection(org, repo, branch, body)

    def cancel_protection_branch(self, org: str, repo: str, branch: str):
        return self.repository_api.delete_v5_repos_owner_repo_branches_branch_protection(org, repo, branch)

    def create_file(self, org: str, repo: str, branch, path, content: str, commit_msg: str) -> gitee.CommitContent:
        content = base64.b64encode(content.encode()).decode()
        body = gitee.NewFileParam(branch=branch, content=content, message=commit_msg)
        return self.repository_api.post_v5_repos_owner_repo_contents_path(org, repo, path, body)

    def get_path_content(self, org: str, repo: str, path, ref: str) -> gitee.Content:
        return self.repository_api.get_v5_repos_owner_repo_contents_path(org, repo, path, ref=ref)

    def get_directory_tree(self, org: str, repo: str, sha: str, recursive: int = 1) -> gitee.Tree:
        return self.git_data_api.get_v5_repos_owner_repo_git_trees_sha(org, repo, sha, recursive=recursive)

    def get_bot(self) -> gitee.User:
        return self.user_api.get_v5_user()

    def get_user_permission_of_repo(self, org: str, repo: str, username: str) -> gitee.ProjectMemberPermission:
        return self.repository_api.get_v5_repos_owner_repo_collaborators_username_permission(org, repo, username)

    def get_enterprise_member(self, enterprise, username: str) -> gitee.EnterpriseMember:
        return self.enterprise_api.get_v5_enterprises_enterprise_members_username(enterprise, username)

    def create_issue(self, org: str, repo: str, title, body: str) -> gitee.Issue:
        param = gitee.IssueCreateParam(repo=repo, title=title, body=body)
        return self.issue_api.post_v5_repos_owner_issues(org, body=param)


def new_client(token: str) -> Client:
    return _Client(token)
