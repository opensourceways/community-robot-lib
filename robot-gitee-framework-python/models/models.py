from typing import List


class Owner(object):
    def __init__(self, id=None, login=None, avatar_url=None, html_url=None, type=None, site_admin=None, name=None,
                 email=None, username=None, user_name=None, url=None, remark=None):
        self.id = id
        self.login = login
        self.avatar_url = avatar_url
        self.html_url = html_url
        self.type = type
        self.site_admin = site_admin
        self.name = name
        self.email = email
        self.username = username
        self.user_name = user_name
        self.url = url
        self.remark = remark


class Project(object):
    def __init__(self, clone_url=None, created_at=None, default_branch=None, description=None, fork=None,
                 forks_count=None, full_name=None, git_http_url=None, git_ssh_url=None, git_svn_url=None, git_url=None,
                 has_issues=None, has_pages=None, has_wiki=None, homepage=None, html_url=None, id=None, language=None,
                 name=None, name_with_namespace=None, namespace=None, open_issues_count=None, owner: Owner = Owner(),
                 path=None, path_with_namespace=None, private=None, pushed_at=None, ssh_url=None, stargazers_count=None,
                 svn_url=None, updated_at=None, url=None, watchers_count=None):
        self.clone_url = clone_url
        self.created_at = created_at
        self.default_branch = default_branch
        self.description = description
        self.fork = fork
        self.forks_count = forks_count
        self.full_name = full_name
        self.git_http_url = git_http_url
        self.git_ssh_url = git_ssh_url
        self.git_svn_url = git_svn_url
        self.git_url = git_url
        self.has_issues = has_issues
        self.has_pages = has_pages
        self.has_wiki = has_wiki
        self.homepage = homepage
        self.html_url = html_url
        self.id = id
        self.language = language
        self.name = name
        self.name_with_namespace = name_with_namespace
        self.namespace = namespace
        self.open_issues_count = open_issues_count
        self.owner = owner
        self.path = path
        self.path_with_namespace = path_with_namespace
        self.private = private
        self.pushed_at = pushed_at
        self.ssh_url = ssh_url
        self.stargazers_count = stargazers_count
        self.svn_url = svn_url
        self.updated_at = updated_at
        self.url = url
        self.watchers_count = watchers_count


class Sender(object):
    def __int__(self, avatar_url=None, email=None, html_url=None, id=None, login=None, name=None, remark=None,
                site_admin=None, type=None, url=None, user_name=None, username=None):
        self.avatar_url = avatar_url
        self.email = email
        self.html_url = html_url
        self.id = id
        self.login = login
        self.name = name
        self.remark = remark
        self.site_admin = site_admin
        self.type = type
        self.url = url
        self.user_name = user_name
        self.username = username


class User(object):
    def __init__(self, id=None, login=None, avatar_url=None, html_url=None, type=None, site_admin=None, name=None,
                 email=None, user_name=None, username=None, url=None, remark=None):
        self.id = id
        self.login = login
        self.avatar_url = avatar_url
        self.html_url = html_url
        self.type = type
        self.site_admin = site_admin
        self.email = email
        self.name = name
        self.url = url
        self.user_name = user_name
        self.username = username
        self.remark = remark


class Labels(object):
    def __init__(self, id=None, name=None, color=None):
        self.id = id
        self.name = name
        self.color = color


class Assignee(object):
    def __init__(self, id=None, login=None, avatar_url=None, html_url=None, type=None, site_admin=None, name=None,
                 email=None, username=None, user_name=None, url=None, remark=None):
        self.id = id
        self.login = login
        self.avatar_url = avatar_url
        self.html_url = html_url
        self.type = type
        self.site_admin = site_admin
        self.name = name
        self.email = email
        self.username = username
        self.user_name = user_name
        self.url = url
        self.remark = remark


class CommitsAuthor(object):
    def __init__(self, email=None, id=None, name=None, remark=None, time=None, url=None, user=None, user_name=None,
                 username=None):
        self.email = email
        self.id = id
        self.name = name
        self.remark = remark
        self.url = url
        self.user_name = user_name
        self.username = username
        self.time = time
        self.user = user


class Author(object):
    def __init__(self, avatar_url=None, email=None, html_url=None, id=None, login=None, name=None, remark=None,
                 site_admin=None, type=None, url=None, user_name=None, username=None):
        self.avatar_url = avatar_url
        self.email = email
        self.html_url = html_url
        self.id = id
        self.login = login
        self.name = name
        self.remark = remark
        self.site_admin = site_admin
        self.type = type
        self.url = url
        self.user_name = user_name
        self.username = username


class Committer(object):
    def __init__(self, id=None, name=None, email=None, username=None, user_name=None, url=None, remark=None, time=None):
        self.id = id
        self.name = name
        self.email = email
        self.username = username
        self.user_name = user_name
        self.url = url
        self.remark = remark
        self.time = time


class Enterprise(object):
    def __init__(self, name=None, url=None):
        self.name = name
        self.url = url


class Repository(object):
    def __init__(self, id=None, name=None, full_name=None, owner: Owner = Owner(), private=None, html_url=None,
                 url=None, description=None, fork=None, created_at=None, updated_at=None, pushed_at=None, git_url=None,
                 ssh_url=None, clone_url=None, svn_url=None, git_http_url=None, git_ssh_url=None, git_svn_url=None,
                 homepage=None, stargazers_count=None, watchers_count=None, forks_count=None, language=None,
                 has_issues=None, has_wiki=None, has_pages=None, license=None, open_issues_count=None,
                 default_branch=None, namespace=None, name_with_namespace=None, path_with_namespace=None, path=None):
        self.id = id
        self.name = name
        self.full_name = full_name
        self.owner = owner
        self.private = private
        self.html_url = html_url
        self.url = url
        self.description = description
        self.fork = fork
        self.created_at = created_at
        self.updated_at = updated_at
        self.pushed_at = pushed_at
        self.git_url = git_url
        self.ssh_url = ssh_url
        self.clone_url = clone_url
        self.svn_url = svn_url
        self.git_http_url = git_http_url
        self.git_ssh_url = git_ssh_url
        self.git_svn_url = git_svn_url
        self.homepage = homepage
        self.stargazers_count = stargazers_count
        self.watchers_count = watchers_count
        self.forks_count = forks_count
        self.language = language
        self.has_issues = has_issues
        self.has_wiki = has_wiki
        self.has_pages = has_pages
        self.license = license
        self.open_issues_count = open_issues_count
        self.default_branch = default_branch
        self.namespace = namespace
        self.name_with_namespace = name_with_namespace
        self.path_with_namespace = path_with_namespace
        self.path = path


class Commits(object):
    def __init__(self, id=None, tree_id=None, distinct=None, message=None, timestamp=None, url=None,
                 author: CommitsAuthor = CommitsAuthor(), committer: Committer = Committer(), added=[], removed=[],
                 modified=[], parent_ids=[]):
        self.id = id
        self.tree_id = tree_id
        self.distinct = distinct
        self.message = message
        self.timestamp = timestamp
        self.url = url
        self.author = author
        self.committer = committer
        self.added = added
        self.removed = removed
        self.modified = modified
        self.parent_ids = parent_ids


class HeadCommit(object):
    def __init__(self, id=None, tree_id=None, distinct=None, message=None, timestamp=None, url=None,
                 author: Author = Author(), committer: Committer = Committer(), added=None, removed=None,
                 modified=[]):
        self.id = id
        self.tree_id = tree_id
        self.distinct = distinct
        self.message = message
        self.timestamp = timestamp
        self.url = url
        self.author = author
        self.committer = committer
        self.added = added
        self.removed = removed
        self.modified = modified


class Milestone(object):
    def __init__(self, html_url=None, id=None, number=None, title=None, description=None, open_issues=None,
                 started_issues=None, closed_issues=None, approved_issues=None, state=None, created_at=None,
                 updated_at=None, due_on=None):
        self.html_url = html_url
        self.id = id
        self.number = number
        self.title = title
        self.description = description
        self.open_issues = open_issues
        self.started_issues = started_issues
        self.closed_issue = closed_issues
        self.approved_issues = approved_issues
        self.state = state
        self.created_at = created_at
        self.updated_at = updated_at
        self.due_on = due_on


class Pusher(object):
    def __init__(self, email=None, id=None, name=None, url=None, user_name=None, username=None):
        self.email = email
        self.id = id
        self.name = name
        self.url = url
        self.user_name = user_name
        self.username = username


class PushPayloadUser(object):
    def __init__(self, email=None, id=None, name=None, url=None, user_name=None, username=None):
        self.email = email
        self.id = id
        self.name = name
        self.url = url
        self.user_name = user_name
        self.username = username


class PushPayload(object):

    def __init__(self, after=None, before=None, commits: List[Commits] = [], commits_more_than_ten=None, compare=None,
                 created=None, deleted=None, enterprise: Enterprise = Enterprise(), head_commit=None, hook_id=None,
                 hook_name=None, hook_url=None, password=None, project=None, push_data=None, pusher: Pusher = Pusher(),
                 ref=None, repository: Repository = Repository(), sender: Sender = Sender(), sign=None, timestamp=None,
                 total_commits_count=None, user_id=None, user_name=None, user: PushPayloadUser = PushPayloadUser,
                 access_token=None):
        self.after = after
        self.before = before
        self.commits = commits
        self.commits_more_than_ten = commits_more_than_ten
        self.compare = compare
        self.created = created
        self.deleted = deleted
        self.enterprise = enterprise
        self.head_commit = head_commit
        self.hook_id = hook_id
        self.hook_name = hook_name
        self.hook_url = hook_url
        self.password = password
        self.project = project
        self.push_data = push_data
        self.pusher = pusher
        self.ref = ref
        self.repository = repository
        self.sender = sender
        self.sign = sign
        self.timestamp = timestamp
        self.total_commits_count = total_commits_count
        self.user = user
        self.user_id = user_id
        self.user_name = user_name
        self.access_token = access_token


class IssueMilestone(object):
    def __init__(self, closed_issues=None, created_at=None, description=None, due_on=None, html_url=None, id=None,
                 number=None, open_issues=None, state=None, title=None, updated_at=None):
        self.closed_issues = closed_issues
        self.created_at = created_at
        self.description = description
        self.due_on = due_on
        self.html_url = html_url
        self.id = id
        self.number = number
        self.open_issues = open_issues
        self.state = state
        self.title = title
        self.updated_at = updated_at


class Issue(object):
    def __init__(self, html_url=None, id=None, number=None, title=None, body=None, state=None, comments=None,
                 create_at=None, update_at=None, user: User = User(), labels=List[Labels],
                 assignee: Assignee = Assignee(), state_name=None, type_name=None,
                 milestone: IssueMilestone = IssueMilestone(),
                 collaborators: List[User] = []):
        self.html_url = html_url
        self.id = id
        self.number = number
        self.title = title
        self.body = body
        self.state = state
        self.comments = comments
        self.create_at = create_at
        self.update_at = update_at
        self.user = user
        self.labels = labels
        self.assignee = assignee
        self.milestone = milestone
        self.state_name = state_name
        self.type_name = type_name
        self.collaborators = collaborators


class IssuePayload(object):
    def __init__(self, action=None, action_desc=None, hook_name=None, password=None, hoos_id=None, hook_url=None,
                 timestamp=None, sign=None, issue: Issue = Issue(), repository: Repository = Repository(),
                 sender: Sender = Sender(), enterprise: Enterprise = Enterprise(), assignee: Assignee = Assignee(),
                 description=None, iid=None, milestone=None, project: Project = Project(), push_data=None,
                 state=None, title=None, target_user=None, updated_by: User = User(), user: User = User(), url=None,
                 access_token=None):
        self.action = action
        self.action_desc = action_desc
        self.assignee = assignee
        self.description = description
        self.iid = iid
        self.milestone = milestone
        self.project = project
        self.push_data = push_data
        self.hook_name = hook_name
        self.password = password
        self.hoos_id = hoos_id
        self.hook_url = hook_url
        self.timestamp = timestamp
        self.sign = sign
        self.issue = issue
        self.repository = repository
        self.sender = sender
        self.enterprise = enterprise
        self.state = state
        self.title = title
        self.target_user = target_user
        self.updated_by = updated_by
        self.url = url
        self.user = user
        self.access_token = access_token


class Repo(object):
    def __init__(self, id=None, name=None, full_name=None, owner: Owner = Owner(), private=None, html_url=None,
                 url=None, description=None, fork=None, created_at=None, updated_at=None, pushed_at=None, git_url=None,
                 ssh_url=None, clone_url=None, svn_url=None, git_http_url=None, git_ssh_url=None, git_svn_url=None,
                 homepage=None, stargazers_count=None, watchers_count=None, forks_count=None, language=None,
                 has_issues=None, has_wiki=None, has_pages=None, license=None, open_issues_count=None,
                 default_branch=None, namespace=None, name_with_namespace=None, path_with_namespace=None):
        self.id = id
        self.name = name
        self.full_name = full_name
        self.owner = owner
        self.private = private
        self.html_url = html_url
        self.url = url
        self.description = description
        self.fork = fork
        self.created_at = created_at
        self.updated_at = updated_at
        self.pushed_at = pushed_at
        self.git_url = git_url
        self.ssh_url = ssh_url
        self.clone_url = clone_url
        self.svn_url = svn_url
        self.git_http_url = git_http_url
        self.git_ssh_url = git_ssh_url
        self.git_svn_url = git_svn_url
        self.homepage = homepage
        self.stargazers_count = stargazers_count
        self.watchers_count = watchers_count
        self.forks_count = forks_count
        self.language = language
        self.has_issues = has_issues
        self.has_wiki = has_wiki
        self.has_pages = has_pages
        self.license = license
        self.open_issues_count = open_issues_count
        self.default_branch = default_branch
        self.namespace = namespace
        self.name_with_namespace = name_with_namespace
        self.path_with_namespace = path_with_namespace


class Head(object):
    def __init__(self, label=None, ref=None, sha=None, user: User = User(), repo: Repo = Repo()):
        self.label = label
        self.ref = ref
        self.sha = sha
        self.user = user
        self.repo = repo
        self.label = label


class Base(object):
    def __init__(self, label=None, ref=None, sha=None, user: User = User(), repo: Repo = Repo()):
        self.label = label
        self.ref = ref
        self.sha = sha
        self.user = user
        self.repo = repo
        self.label = label


class PullRequest(object):
    def __init__(self, id=None, number=None, state=None, html_url=None, diff_url=None, patch_url=None, title=None,
                 body=None, created_at=None, updated_at=None, closed_at=None, merged_at=None, merge_commit_sha=None,
                 user: User = User(), assignee: Assignee = Assignee(), tester=None, milestone: Milestone = Milestone(),
                 languages: List[str] = [], merge_reference_name=None, merge_status=None,
                 head: Head = Head(), base: Base = Base(), merged=None, mergeable=None, comments=None,
                 need_review=None, need_test=None, stale_issues=[], stale_labels=[],
                 commits: List[Commits] = [], issues=[], labels=[], additions=None, deletions=None, changed_files=None,
                 assignees: List[Assignee] = None, testers: List[User] = [], updated_by: List[User] = []):
        self.id = id
        self.number = number
        self.state = state
        self.languages = languages
        self.issues = issues
        self.labels = labels
        self.merge_reference_name = merge_reference_name
        self.merge_status = merge_status
        self.need_review = need_review
        self.need_test = need_test
        self.stale_issues = stale_issues
        self.stale_labels = stale_labels
        self.testers = testers
        self.html_url = html_url
        self.diff_url = diff_url
        self.patch_url = patch_url
        self.updated_by = updated_by
        self.title = title
        self.body = body
        self.created_at = created_at
        self.updated_at = updated_at
        self.closed_at = closed_at
        self.merged_at = merged_at
        self.merge_commit_sha = merge_commit_sha
        self.user = user
        self.assignee = assignee
        self.tester = tester
        self.milestone = milestone
        self.head = head
        self.base = base
        self.merged = merged
        self.mergeable = mergeable
        self.comments = comments
        self.commits = commits
        self.additions = additions
        self.deletions = deletions
        self.changed_files = changed_files
        self.assignees = assignees


class SourceRepo(object):
    def __init__(self, project: Project = Project(), repository: Repository = Repository()):
        self.project = project
        self.repository = repository


class PullRequestPayload(object):
    def __init__(self, hook_name=None, title=None, target_user: User = User(), state=None, target_branch=None,
                 password=None, url=None, action=None, action_desc=None, body=None,
                 hook_id=None, hook_url=None, updated_by: User = User(), iid=None,
                 number=None, project: Project = Project(), timestamp=None, merge_status=None, merge_commit_sha=None,
                 push_data=None, repository: Repository = Repository(), languages: List[str] = [], sign=None,
                 target_repo: SourceRepo = SourceRepo(), source_repo: SourceRepo = SourceRepo(),
                 pull_request: PullRequest = PullRequest(), author: Author = Author(), sender: Sender = Sender(),
                 enterprise: Enterprise = Enterprise(), source_branch=None,
                 access_token=None):
        self.action = action
        self.state = state
        self.url = url
        self.updated_by = updated_by
        self.target_user = target_user
        self.title = title
        self.target_repo = target_repo
        self.source_repo = source_repo
        self.target_branch = target_branch
        self.action_desc = action_desc
        self.body = body
        self.iid = iid
        self.languages = languages
        self.merge_commit_sha = merge_commit_sha
        self.merge_status = merge_status
        self.number = number
        self.project = project
        self.push_data = push_data
        self.repository = repository
        self.hook_name = hook_name
        self.password = password
        self.hook_id = hook_id
        self.hook_url = hook_url
        self.timestamp = timestamp
        self.sign = sign
        self.pull_request = pull_request
        self.author = author
        self.sender = sender
        self.enterprise = enterprise
        self.access_token = access_token
        self.source_branch = source_branch


class Comment(object):
    def __init__(self, html_url=None, id=None, body=None, user: User = User(), created_at=None, updated_at=None):
        self.html_url = html_url
        self.id = id
        self.body = body
        self.user = user
        self.created_at = created_at
        self.updated_at = updated_at


class NotePayload(object):
    def __init__(self, action=None, action_desc=None, author: Author = Author(), comment: Comment = Comment(),
                 enterprise: Enterprise = Enterprise(), hook_id=None, hook_name=None, hook_url=None,
                 issue: Issue = Issue(), note=None, noteable_id=None, noteable_type=None, password=None, per_iid=None,
                 project: Project = Project(), pull_request: PullRequest = PullRequest(), push_data=None,
                 repository: Repository = Repository(), sender: Sender = Sender(), short_commit_id=None, sign=None,
                 timestamp=None, title=None, url=None, access_token=None):
        self.action = action
        self.action_desc = action_desc
        self.author = author
        self.comment = comment
        self.enterprise = enterprise
        self.hook_id = hook_id
        self.hook_name = hook_name
        self.hook_url = hook_url
        self.issue = issue
        self.note = note
        self.noteable_id = noteable_id
        self.noteable_type = noteable_type
        self.password = password
        self.per_iid = per_iid
        self.project = project
        self.pull_request = pull_request
        self.push_data = push_data
        self.repository = repository
        self.sender = sender
        self.short_commit_id = short_commit_id
        self.sign = sign
        self.timestamp = timestamp
        self.title = title
        self.url = url
        self.access_token = access_token
