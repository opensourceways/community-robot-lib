from dataclasses import dataclass


@dataclass
class ProgramBasic:
    id: str
    name: str
    description: str
    assignee: str
    author: str
