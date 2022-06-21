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
        """ 创建 pr """
        body = gitee.CreatePullRequestParam(title=title, body=body, base=base, head=head, prune_source_branch=True)
        return self.pr_api.post_v5_repos_owner_repo_pulls(org, repo, body)

    def get_pull_requests(self, org: str, repo: str, opts: ListPullRequestOpt = None) -> List[gitee.PullRequest]:
        """ 获取Pull Request列表 """
        if opts is not None:
            param = vars(opts)
            return self.pr_api.get_v5_repos_owner_repo_pulls(org, repo, **param)
        else:
            return self.pr_api.get_v5_repos_owner_repo_pulls(org, repo)

    def update_pull_request(self, org: str, repo: str, number: int,
                            body: gitee.PullRequestUpdateParam) -> gitee.PullRequest:
        return self.pr_api.patch_v5_repos_owner_repo_pulls_number(org, repo, number, body)

    def list_collaborators(self, org: str, repo: str) -> List[gitee.ProjectMember]:
        """ 获取仓库的所有成员 """
        return self.repository_api.get_v5_repos_owner_repo_collaborators(org, repo)

    def is_collaborator(self, org: str, repo: str, username: str) -> bool:
        """判断用户是否为仓库成员"""
        try:
            self.repository_api.get_v5_repos_owner_repo_collaborators_username(org, repo, username)
        except Exception as e:
            return False
        return True

    def is_member(self, org: str, username: str) -> bool:
        """判断用户是否为组织成员"""
        try:
            self.organization_api.get_v5_orgs_org_memberships_username(org, username)
        except Exception as e:
            return False
        return True

    def remove_repo_member(self, org: str, repo: str, username: str):
        """移除仓库成员"""
        return self.repository_api.delete_v5_repos_owner_repo_collaborators_username(org, repo, username)

    def add_repo_member(self, org: str, repo: str, username: str, permission: str):
        """添加仓库成员"""
        body = gitee.ProjectMemberPutParam(permission=permission)
        self.repository_api.put_v5_repos_owner_repo_collaborators_username(org, repo, username, body)

    def get_ref(self, org: str, repo: str, ref: str) -> str:
        """ 获取分支 sha"""
        branch = self.repository_api.get_v5_repos_owner_repo_branches_branch(org, repo, ref)
        return branch.commit.sha

    def get_pull_request_changes(self, org: str, repo: str, number: int) -> List[gitee.PullRequestFiles]:
        """ Pull Request Commit文件列表。最多显示300条diff """
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_files(org, repo, number)

    def get_pr_labels(self, org: str, repo: str, number: int) -> List[gitee.Label]:
        """ 获取某个 Pull Request 的所有标签 """
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_labels(org, repo, number)

    def list_pr_comments(self, org: str, repo: str, number: int) -> List[gitee.PullRequestComments]:
        """ 获取某个Pull Request的所有评论 """
        comments = self.pr_api.get_v5_repos_owner_repo_pulls_number_comments(org, repo, number)
        return comments

    def list_pr_issues(self, org: str, repo: str, number: int) -> List[gitee.Issue]:
        """ 获取 Pull Request 关联的 issues """
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_issues(org, repo, number)

    def delete_pr_comment(self, org: str, repo: str, id: int):
        """删除Pull Request评论"""
        self.pr_api.delete_v5_repos_owner_repo_pulls_comments_id(org, repo, id)

    def create_pr_comment(self, org: str, repo: str, number: int, comment: str):
        """提交Pull Request评论"""
        body = gitee.PullRequestCommentPostParam(body=comment)
        return self.pr_api.post_v5_repos_owner_repo_pulls_number_comments(org, repo, number, body)

    def update_pr_comment(self, org: str, repo: str, comment_id: int, comment: str):
        """ 编辑评论 """
        body = gitee.PullRequestCommentPatchParam(body=comment)
        return self.pr_api.patch_v5_repos_owner_repo_pulls_comments_id(org, repo, comment_id, body)

    def add_multi_pr_label(self, org: str, repo: str, number: int, labels: List[str]):
        return self.pr_api.post_v5_repos_owner_repo_pulls_number_labels(org, repo, number, labels)

    def add_pr_label(self, org: str, repo: str, number: int, label: str):
        labels = [label]
        return self.add_multi_pr_label(org, repo, number, labels)

    def remove_pr_label(self, org: str, repo: str, number: int, label: str):
        """ 删除 Pull Request 标签 """
        return self.pr_api.delete_v5_repos_owner_repo_pulls_label(org, repo, number, label)

    def remove_pr_labels(self, org: str, repo: str, number: int, labels: List[str]):
        """ 删除 Pull Request 标签"""
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
        """ 替换 Pull Request 所有标签 """
        return self.pr_api.put_v5_repos_owner_repo_pulls_number_labels(owner, repo, number, labels)

    def list_pr_operation_logs(self, org: str, repo: str, number: int) -> List[gitee.OperateLog]:
        """ 获取某个Pull Request的操作日志 """
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_operate_logs(org, repo, number)

    def close_pr(self, org: str, repo: str, number: int):
        """ 关闭pr """
        body = gitee.PullRequestUpdateParam(state='closed')
        return self.update_pull_request(org, repo, number, body)

    def assign_pr(self, owner, repo: str, number: int, usernames: List[str]):
        """ 指派用户审查 Pull Request """
        assignees = ','.join(usernames)
        body = gitee.PullRequestAssigneePostParam(assignees=assignees)
        return self.pr_api.post_v5_repos_owner_repo_pulls_number_assignees(owner, repo, number, body)

    def unassign_pr(self, owner, repo: str, number: int, usernames: List[str]):
        assignees = ','.join(usernames)
        return self.pr_api.delete_v5_repos_owner_repo_pulls_number_assignees(org, repo, number, assignees)

    def get_pr_commits(self, org: str, repo: str, number: int) -> List[gitee.PullRequestCommits]:
        """ 获取某Pull Request的所有Commit信息。最多显示250条Commit """
        return self.pr_api.get_v5_repos_owner_repo_pulls_number_commits(org, repo, number)

    def get_pull_request(self, org: str, repo: str, number: int) -> gitee.PullRequest:
        """获取单个Pull Request"""
        return self.pr_api.get_v5_repos_owner_repo_pulls_number(org, repo, number)

    def get_repo_commit(self, org: str, repo: str, sha: str) -> gitee.RepoCommit:
        """ 仓库的某个提交 """
        return self.repository_api.get_v5_repos_owner_repo_commits_sha(org, repo, sha)

    def merge_pr(self, owner, repo: str, number: int, merge_method='merge', title=None, description=None
                 , prune_source_branch=False):
        body = gitee.PullRequestMergePutParam(merge_method=merge_method, prune_source_branch=prune_source_branch,
                                              title=title, description=description)
        self.pr_api.put_v5_repos_owner_repo_pulls_number_merge(owner, repo, number, body)

    # "message":"Group"
    def get_org_repos(self, org: str) -> List[gitee.Project]:
        """获取组织的仓库列表"""
        return self.repository_api.get_v5_orgs_org_repos(org)

    # "message":"Group"
    def create_org_repo(self, org: str, repo: gitee.RepositoryPostParam):
        """ 创建组织仓库 """
        return self.repository_api.post_v5_orgs_org_repos(org, repo)

    def update_repo(self, org: str, repo: str, info: gitee.RepoPatchParam):
        """ 修改仓库 """
        return self.repository_api.patch_v5_repos_owner_repo(org, repo, info)

    def get_repo(self, org: str, repo: str) -> gitee.Project:
        """ 获取用户的某个仓库  """
        return self.repository_api.get_v5_repos_owner_repo(org, repo)

    def get_gitee_repo(self, org: str, repo: str) -> gitee.Project:
        """ 获取用户的某个仓库  """
        return self.repository_api.get_v5_repos_owner_repo(org, repo)

    def set_repo_reviewer(self, org: str, repo: str, reviewer: gitee.SetRepoReviewer):
        """ 修改代码审查设置 """
        return self.repository_api.put_v5_repos_owner_repo_reviewer(org, repo, reviewer)

    def create_repo_label(self, org: str, repo: str, label, color: str):
        """ 创建仓库任务标签 """
        body = gitee.LabelPostParam(name=label, color=color)
        return self.label_api.post_v5_repos_owner_repo_labels(org, repo, body)

    def get_repo_labels(self, owner, repo: str) -> List[gitee.Label]:
        """ 获取仓库所有任务标签 """
        return self.label_api.get_v5_repos_owner_repo_labels(owner, repo)

    def assign_gitee_issue(self, org: str, repo: str, number: str, username: str):
        """ 指派Issue给用户 """
        body = gitee.IssueUpdateParam(repo=repo, assignee=username)
        return self.issue_api.patch_v5_repos_owner_issues_number(org, number, body)

    def unassign_gitee_issue(self, org: str, repo: str, number: str):
        """ 取消指派 """
        return self.assign_gitee_issue(org, repo, number, " ")

    def remove_issue_assignee(self, org: str, repo: str, number: str):
        """ 取消指派 """
        return self.assign_gitee_issue(org, repo, number, " ")

    def create_issue_comment(self, org: str, repo: str, number: str, comment: str):
        """ 创建某个Issue评论 """
        body = gitee.IssueCommentPostParam(body=comment)
        return self.issue_api.post_v5_repos_owner_repo_issues_number_comments(org, repo, number, body)

    def update_issue_comment(self, org: str, repo: str, comment_id: int, comment: str):
        """ 更新Issue某条评论 """
        body = gitee.IssueCommentPatchParam(body=comment)
        return self.issue_api.patch_v5_repos_owner_repo_issues_comments_id(org, repo, comment_id, body)

    def list_issue_comments(self, org: str, repo: str, number: str) -> List[gitee.Note]:
        """ 获取仓库某个Issue所有的评论 """
        return self.issue_api.get_v5_repos_owner_repo_issues_number_comments(org, repo, number)

    def get_issue_labels(self, org: str, repo: str, number: str) -> List[gitee.Label]:
        """ 获取仓库任务的所有标签 """
        return self.label_api.get_v5_repos_owner_repo_issues_number_labels(org, repo, number)

    def remove_issue_label(self, org: str, repo: str, number: str, label: str):
        """ 删除Issue标签 """
        return self.label_api.delete_v5_repos_owner_repo_issues_number_labels_name(org, repo, number, label)

    def remove_issue_labels(self, org: str, repo: str, number: str, labels: List[str]):
        """ 删除Issue标签  """
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
        """ 创建Issue标签 """
        labels = [label]
        return self.label_api.post_v5_repos_owner_repo_issues_number_labels(org, repo, number, labels)

    def add_multi_issue_label(self, org: str, repo: str, number: str, labels: List[str]):
        """ 创建Issue标签 """
        return self.label_api.post_v5_repos_owner_repo_issues_number_labels(org, repo, number, labels)

    def update_issue(self, owner, number: str, param: gitee.IssueUpdateParam) -> gitee.Issue:
        """ 更新 Issue """
        return self.issue_api.patch_v5_repos_owner_issues_number(owner, number, param)

    def close_issue(self, owner, repo: str, number: str):
        """ 关闭 Issue """
        param = gitee.IssueUpdateParam(repo=repo, state='closed')
        return self.update_issue(owner, number, param)

    def reopen_issue(self, owner, repo: str, number: str):
        """ 打开 Issue"""
        param = gitee.IssueUpdateParam(repo=repo, state='open')
        return self.update_issue(owner, number, param)

    def get_issue(self, org: str, repo: str, number: str) -> gitee.Issue:
        """ 获取某个 Issue """
        return self.issue_api.get_v5_repos_owner_repo_issues_number(org, repo, number)

    # 没找到对应的api
    def add_project_labels(self, org: str, repo: str, labels: List[str]):
        return self.label_api.post_v5_repos_owner_repo_project_labels(org, repo, labels)

    def update_project_labels(self, org: str, repo: str, labels: List[str]):
        return self.label_api.put_v5_repos_owner_repo_project_labels(org, repo, labels)

    def create_branch(self, org: str, repo: str, branch: str, parent_branch: str):
        """ 创建分支 """
        body = gitee.CreateBranchParam(branch_name=branch, refs=parent_branch)
        return self.repository_api.post_v5_repos_owner_repo_branches(org, repo, body)

    def get_repo_all_branch(self, org: str, repo: str) -> List[gitee.Branch]:
        """ 获取所有分支 """
        return self.repository_api.get_v5_repos_owner_repo_branches(org, repo)

    def set_protection_branch(self, org: str, repo: str, branch: str):
        """ 设置分支保护 """
        body = gitee.BranchProtectionPutParam()
        return self.repository_api.put_v5_repos_owner_repo_branches_branch_protection(org, repo, branch, body)

    # 有问题 "message":"Operation is not allowed"
    def cancel_protection_branch(self, org: str, repo: str, branch: str):
        """ 取消保护分支的设置 """
        return self.repository_api.delete_v5_repos_owner_repo_branches_branch_protection(org, repo, branch)

    def create_file(self, org: str, repo: str, branch, path, content: str, commit_msg: str) -> gitee.CommitContent:
        """ 创建文件 """
        content = base64.b64encode(content.encode()).decode()
        body = gitee.NewFileParam(branch=branch, content=content, message=commit_msg)
        return self.repository_api.post_v5_repos_owner_repo_contents_path(org, repo, path, body)

    def get_path_content(self, org: str, repo: str, path, ref: str) -> gitee.Content:
        """ 获取仓库具体路径下的内容 """
        return self.repository_api.get_v5_repos_owner_repo_contents_path(org, repo, path, ref=ref)

    def get_directory_tree(self, org: str, repo: str, sha: str, recursive: int = 1) -> gitee.Tree:
        """ 获取目录Tree """
        return self.git_data_api.get_v5_repos_owner_repo_git_trees_sha(org, repo, sha, recursive=recursive)

    def get_bot(self) -> gitee.User:
        """ 获取机器人信息 """
        return self.user_api.get_v5_user()

    def get_user_permission_of_repo(self, org: str, repo: str, username: str) -> gitee.ProjectMemberPermission:
        """ 查看仓库成员的权限 """
        return self.repository_api.get_v5_repos_owner_repo_collaborators_username_permission(org, repo, username)

    # 没有企业，无法测试
    def get_enterprise_member(self, enterprise, username: str) -> gitee.EnterpriseMember:
        """ 获取企业的一个成员 """
        return self.enterprise_api.get_v5_enterprises_enterprise_members_username(enterprise, username)

    def create_issue(self, org: str, repo: str, title, body: str) -> gitee.Issue:
        """创建 Issue"""
        param = gitee.IssueCreateParam(repo=repo, title=title, body=body)
        return self.issue_api.post_v5_repos_owner_issues(org, body=param)


def new_client(token: str) -> Client:
    return _Client(token)


if __name__ == '__main__':
    c = new_client('')
    org = 'william_wong'
    repo = 'robot'

    # create_pull_request
    # print(c.get_pull_requests(org, repo))

    # print(c.list_collaborators(org, repo))
    # print(c.is_collaborator(org, repo, 'william_wong'))
    # print(c.is_member(org, 'william_wong')) # 没有组织无法识别
    # print(c.remove_repo_member(org, repo, 'lxxhope'))
    # print(c.add_repo_member(org, repo, 'lxxhope', 'admin'))
    # print(c.get_ref(org, repo, 'master'))
    # print(c.get_pull_request_changes(org, repo, 1))
    # print(c.get_pr_labels(org, repo, 1))
    # print(c.list_pr_comments(org, repo, 1))
    # print(c.list_pr_issues(org, repo, 2))
    # print(c.delete_pr_comment(org, repo, 11006054))
    # print(c.create_pr_comment(org, repo, 2, 'it is good'))
    # print(c.update_pr_comment(org, repo, 11053733, 'it is not good'))
    # print(c.add_pr_label(org, repo, 2, 'balabala'))
    # print(c.add_multi_pr_label(org, repo, 2, labels=['bug', 'lalala']))
    # print(c.remove_pr_label(org, repo, 2, 'bug'))
    # print(c.get_repos(org))
    # reviewer = gitee.SetRepoReviewer(assignees='lxxhope', testers='lxxhope', assignees_number=1, testers_number=1)
    # print(c.set_repo_reviewer(org, repo,reviewer))
    # print(c.create_repo_label(org, repo, 'llll111', '000000'))
    # print(c.get_repo_labels(org,repo))
    # print(c.assign_gitee_issue(org, repo, 'I5D4O2', 'lxxhope'))
    # print(c.unassign_gitee_issue(org, repo, 'I5D4O2'))
    # print(c.create_issue_comment(org, repo, 'I5A4ZC', 'aczxc12321321'))
    # print(c.update_issue_comment(org, repo, 11056464, 'aaaaaaa11'))
    # print(c.list_issue_comments(org, repo, 'I5A4ZC'))
    # print(c.get_issue_labels(org, repo, 'I5A4ZC'))
    # print(c.remove_issue_label(org, repo, 'I5A4ZC', 'hello1'))
    # print(c.add_multi_issue_label(org, repo, 'I5A4ZC', ["aaa", "bbb"]))
    # param = gitee.IssueUpdateParam(repo='robot')
    # print(c.update_issue(org, 'I5A4ZC', param))
    # print(c.close_issue(org, repo, 'I5A4ZC'))
    # print(c.reopen_issue(org, repo, 'I5A4ZC'))
    # print(c.get_issue(org, repo, 'I5A4ZC'))
    # print(c.create_branch(org, repo, 'test111', 'master'))
    # print(c.set_protection_branch(org, repo, 'master'))
    # print(c.cancel_protection_branch(org, repo, 'master'))
    # print(c.create_file(org, repo, 'master', './conf/2.txt', 'aaaaaaaaabbbbbb', 'testststst'))
    # print(c.get_path_content(org, repo, '/conf/1.txt', 'master'))
    # print(c.get_directory_tree(org, repo, 'master', 1))
    # print(c.create_pull_request(org, repo, '111122', '2222', 'master', 'test111'))
    # print(c.get_pull_request_changes(org, repo, 1))
    # c.create_pr_comment(org, repo, 1, 'memeda')
    # print(c.add_multi_pr_label(org, repo, 1, ['zzz', 'xxx', 'yyy']))
    # print(c.get_pull_request(org, repo, 1))
    # print(c.create_issue(org, repo, title='hello issue', body='tttttttttttttttttt'))
    # c.update_pull_request(org, repo,)
    # opt = ListPullRequestOpt()
    # print(c.get_pull_requests(org, repo))
    # body = gitee.PullRequestUpdateParam()
    # # body.state = 'open'
    # body.title = 'test pr 111'
    # print(c.update_pull_request(org, repo, 1, body))
    # print(c.assign_pr(org, repo, 1, ['lxxhope']))
    # print(c.unassign_pr(org, repo, 1, ['lxxhope']))
    # print(c.get_ref(org, repo, 'master'))
    # print(c.create_issue(org, repo, title='hahahahaahahhahaha123456', body='hello issue'))
    # print(c.remove_pr_label(org, repo, 2, 'bug'))
    # print(c.remove_pr_labels(org, repo, 2, ['aaaa']))
    # print(c.replace_pr_all_labels(org, repo, 2, ["aaaa"]))
    # print(c.list_pr_operation_logs(org, repo, 2))
    # c.close_pr(org,repo,2)
    # print(c.assign_pr(org, repo, 2, ['lxxhope']))
    # print(c.unassign_pr(org, repo, 2, ['lxxhope']))
    # print(c.get_pr_commits(org, repo, 2))
    # print(c.get_pull_request(org, repo, 2))
    # print(c.get_repo_commit(org, repo, '6df68cb007c5db87aff889e41eb566d6b948a6b1'))
    # print(c.merge_pr(org, repo, 2)) # "message":"此 Pull Request 未通过设置的审查"
    # print(c.get_org_repos(org))
    # print(c.get_repo(org, repo))
    # info = gitee.RepoPatchParam(description='1234567890',name='robot')
    # print(c.update_repo(org, repo, info))
    # print(c.get_enterprise_member(org, 'william_wong'))
    # print(c.get_issue_labels(org, repo, "I5A4ZC"))
    # print(c.add_issue_label(org, repo, 'I5A4ZC', "hello1"))
    # print(c.add_multi_issue_label(org, repo, 'I5A4ZC', ["a1", "a2", "a3"]))
    # c.remove_issue_labels(org, repo, 'I5A4ZC', ["a3", "a6", "a711"])
    # print(c.remove_issue_label(org, repo, 'I5A4ZC', "a6"))
    # print(c.set_protection_branch(org, repo, 'test111'))
    # print(c.cancel_protection_branch(org, repo, 'test111'))
    # print(c.create_issue(org, repo, "123456123", "hahahaahahaha"))
    # print(c.update_project_labels(org,repo,["z","x","y"]))
    # print(c.add_project_labels(org, repo, ["z1", "x1", "y1"]))
