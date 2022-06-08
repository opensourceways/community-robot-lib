from typing import List

from model_user_basic import UserBasic
from model_label import Label
from model_milestone import Milestone
from model_program_basic import ProgramBasic


class Issue:
    id: int
    url: str
    repository_url: str
    labels_url: str
    comments_url: str
    html_url: str
    parent_url: str
    number: str
    state: str
    title: str
    body: str
    body_html: str
    user: UserBasic
    labels: List[Label]
    assignee: UserBasic
    collaborators: List[UserBasic]
    repository: str
    milestone: Milestone
    created_at: str
    updated_at: str
    plan_started_at: str
    deadline: str
    finished_at: str
    scheduled_time: str
    comments: int
    issue_type: str
    program: ProgramBasic
