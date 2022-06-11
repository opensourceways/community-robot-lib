from __future__ import print_function

from gitee import UsersApi, RepositoriesApi, PullRequestsApi, IssuesApi, EmailsApi, EnterprisesApi, GistsApi, \
    GitDataApi, ActivityApi, LabelsApi, MilestonesApi, MiscellaneousApi, OrganizationsApi, SearchApi, WebhooksApi


class GiteeApi(ActivityApi, EmailsApi, EnterprisesApi, GistsApi, GitDataApi, IssuesApi, LabelsApi, MilestonesApi,
               MiscellaneousApi, OrganizationsApi, PullRequestsApi, RepositoriesApi, SearchApi, UsersApi,
               WebhooksApi):
    def __init__(self, api_client=None):
        super().__init__(api_client)
