from dataclasses import dataclass


@dataclass
class Label:
    id: int
    name: str
    color: str
    repository_id: int
    url: str
