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
const (
	botName           = "test"
	statusOpen        = "opened"
	statusUpdate      = "update"
	commentClearLabel = `New code changes of pr are detected and remove these labels ***%s***. :flushed: `
)

type iClient interface {
	UpdateMergeRequest(pid interface{}, mrID int, options gitlab.UpdateMergeRequestOptions) (gitlab.MergeRequest, error)
	GetMergeRequest(pid interface{}, mrID int) (gitlab.MergeRequest, error)
}

func newRobot(cli iClient) *robot {
	return &robot{cli: cli}
}

type robot struct {
	cli iClient
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
	f.RegisterMergeCommentEventHandler(bot.handleMergeCommentEvent)
	f.RegisterIssueCommentHandler(bot.handleIssueCommentEvent)
	f.RegisterIssueHandler(bot.handleIssueEvent)
	f.RegisterPushEventHandler(bot.handlePushEvent)
}

func (bot *robot) handleIssueCommentEvent(e *gitlab.IssueCommentEvent, c config.Config, log *logrus.Entry) error {
	fmt.Println(e.ObjectAttributes.Note)
	return nil
}

func (bot *robot) handleIssueEvent(e *gitlab.IssueEvent, c config.Config, log *logrus.Entry) error {
	fmt.Println(e.ObjectAttributes.Action)
	return nil
}

func (bot *robot) handlePushEvent(e *gitlab.PushEvent, c config.Config, log *logrus.Entry) error {
	fmt.Println(e.Project.Name)
	return nil
}

func (bot *robot) handleMergeEvent(e *gitlab.MergeEvent, c config.Config, log *logrus.Entry) error {
	// TODO: if it doesn't needed to hand PR event, delete this function.
	state := e.ObjectAttributes.State
	if state != statusOpen {
		return nil
	}
	//p := "sig"
	//ref := "master"
	//recursive := true
	//tree, _, err := bot.cli.Repositories.ListTree(e.Project.ID, &gitlab.ListTreeOptions{Path: &p, Ref: &ref, Recursive: &recursive})
	//fmt.Println(tree)
	//fmt.Println(len(tree))
	err := bot.clearLabels(e)
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (bot *robot) handleMergeCommentEvent(e *gitlab.MergeCommentEvent, c config.Config, log *logrus.Entry) error {

	if e.ObjectAttributes.Note == "/lgtm" {
		pid := e.ProjectID
		mrNumber := e.MergeRequest.IID
		_, err := bot.cli.UpdateMergeRequest(pid, mrNumber, gitlab.UpdateMergeRequestOptions{AddLabels: &gitlab.Labels{"lgtm"}})
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (bot *robot) clearLabels(e *gitlab.MergeEvent) error {
	if e.ObjectAttributes.Action != statusUpdate {
		return nil
	}
	if e.ObjectAttributes.OldRev != "" {
		fmt.Println("source branch has been changed")
		mrID := e.ObjectAttributes.IID
		projectID := e.Project.ID
		labels := gitlab.Labels{}
		for _, i := range e.Labels {
			if i.Name == "cla/yes" || i.Name == "sig/TC" {
				continue
			}
			labels = append(labels, i.Name)
		}

		_, err := bot.cli.UpdateMergeRequest(projectID, mrID, gitlab.UpdateMergeRequestOptions{
			RemoveLabels: &labels,
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
