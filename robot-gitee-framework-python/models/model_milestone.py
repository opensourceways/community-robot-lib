from dataclasses import dataclass


@dataclass
class Milestone:
    url: str
    html_url: str
    number: int
    repository_id: int
    state: str
    title: str
    description: str
    updated_at: str
    created_at: str
    open_issues: int
    closed_issues: int
    due_on: str
