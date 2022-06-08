from gitee_api.gitee_api import GiteeApi
from gitee_robot import GiteeRobot
from handler.event_handler import NoteHandler
from models.models import NotePayload


def get_issue(note_payload: NotePayload):
    owner = note_payload.repository.owner.user_name
    repo = note_payload.repository.name
    number = note_payload.issue.number
    access_token = note_payload.access_token
    gitee_api = GiteeApi()
    result = gitee_api.get_v5_repos_owner_repo_issues_number_comments(owner=owner, repo=repo,
                                                                      number=number, access_token=access_token)
    print(result)


if __name__ == '__main__':
    client = GiteeRobot('./conf/conf.yaml')
    note_handler = NoteHandler(get_issue)
    client.add_handler(note_handler)
    client.run()
