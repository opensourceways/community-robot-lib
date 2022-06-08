# func ConvertToPREvent(payload []byte) (e PullRequestEvent, err error) {
# 	if err = json.Unmarshal(payload, &e); err != nil {
# 		return
# 	}
#
# 	err = checkPullRequestEvent(&e)
# 	return
# }

def convert_to_pr_event(payload):
    pass
