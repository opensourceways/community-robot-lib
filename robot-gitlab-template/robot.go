package main

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"

	"github.com/opensourceways/community-robot-lib/config"
	"github.com/opensourceways/community-robot-lib/robot-gitlab-framework"
	"github.com/sirupsen/logrus"
)

// TODO: set botName
const botName = "test"

func newRobot(cli gitlab.Client) *robot {
	return &robot{cli: cli}
}

type robot struct {
	cli gitlab.Client
}

func (bot *robot) NewConfig() config.Config {
	return &configuration{}
}

func (bot *robot) RobotName() string {
	return botName
}

func (bot *robot) getConfig(cfg config.Config) (*configuration, error) {
	if c, ok := cfg.(*configuration); ok {
		return c, nil
	}
	return nil, errors.New("can't convert to configuration")
}

func (bot *robot) RegisterEventHandler(f framework.HandlerRegister) {
	f.RegisterMergeEventHandler(bot.handleMergeEvent)
}

func (bot *robot) handleMergeEvent(e *gitlab.MergeEvent, c config.Config, log *logrus.Entry) error {
	// TODO: if it doesn't needed to hand PR event, delete this function.
	ac := e.ObjectAttributes.Action
	fmt.Println(ac)
	project := e.Project.ID
	mrNumber := e.ObjectAttributes.ID
	mr, _, err := bot.cli.MergeRequests.GetMergeRequest(project, mrNumber, &gitlab.GetMergeRequestsOptions{})
	if err != nil {
		log.Error(err)
	}
	fmt.Println(mr)
	return nil
}
