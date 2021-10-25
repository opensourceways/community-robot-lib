package main

import (
	"errors"

	sdk "gitee.com/openeuler/go-gitee/gitee"
	libconfig "github.com/opensourceways/community-robot-lib/config"
	libplugin "github.com/opensourceways/community-robot-lib/giteeplugin"
	"github.com/sirupsen/logrus"
)

// TODO: set botName
const botName = ""

type iClient interface {
}

type demonBot struct {
	cli iClient
}

func (bot *demonBot) NewPluginConfig() libconfig.PluginConfig {
	return &configuration{}
}

func (bot *demonBot) botConfig(cfg libconfig.PluginConfig) (*configuration, error) {
	if c, ok := cfg.(*configuration); ok {
		return c, nil
	}
	return nil, errors.New("can't convert to configuration")
}

func (bot *demonBot) RegisterEventHandler(p libplugin.HandlerRegitster) {
	p.RegisterIssueHandler(bot.handleIssueEvent)
	p.RegisterPullRequestHandler(bot.handlePREvent)
	p.RegisterNoteEventHandler(bot.handleNoteEvent)
	p.RegisterPushEventHandler(bot.handlePushEvent)
}

func (bot *demonBot) handlePREvent(e *sdk.PullRequestEvent, cfg libconfig.PluginConfig, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand PR event, delete this function.
	return nil
}

func (bot *demonBot) handleIssueEvent(e *sdk.IssueEvent, cfg libconfig.PluginConfig, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand Issue event, delete this function.
	return nil
}

func (bot *demonBot) handlePushEvent(e *sdk.PushEvent, cfg libconfig.PluginConfig, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand Push event, delete this function.
	return nil
}

func (bot *demonBot) handleNoteEvent(e *sdk.NoteEvent, cfg libconfig.PluginConfig, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand Note event, delete this function.
	return nil
}
