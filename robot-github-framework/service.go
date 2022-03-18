package framework

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/community-robot-lib/config"
	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/opensourceways/community-robot-lib/options"
)

type HandlerRegister interface {
	RegisterIssueHandler(IssueHandler)
	RegisterPullRequestHandler(PullRequestHandler)
	RegisterPushEventHandler(PushEventHandler)
	RegisterIssueCommentHandler(IssueCommentHandler)
	RegisterStatusEventHandler(StatusEventHandler)
	RegisterReviewEventHandler(ReviewEventHandler)
	RegisterReviewCommentEventHandler(ReviewCommentEventHandler)
	RegisterCommitCommentEventHandler(CommitCommentEventHandler)
}

type Robot interface {
	NewConfig() config.Config
	RegisterEventHandler(HandlerRegister)
	RobotName() string
}

func Run(bot Robot, o options.ServiceOptions) {
	agent := config.NewConfigAgent(bot.NewConfig)
	if err := agent.Start(o.ConfigFile); err != nil {
		logrus.WithError(err).Errorf("start config:%s", o.ConfigFile)
		return
	}

	h := handlers{}
	bot.RegisterEventHandler(&h)

	d := &dispatcher{agent: &agent, h: h}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.OnInterrupt(func() {
		agent.Stop()
		d.Wait()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	http.Handle("/github-hook", d)

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(o.Port)}

	interrupts.ListenAndServe(httpServer, o.GracePeriod)
}
