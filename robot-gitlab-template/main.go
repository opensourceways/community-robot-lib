package main

import (
	"flag"
	"github.com/opensourceways/community-robot-lib/gitlabclient"
	"os"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/opensourceways/community-robot-lib/robot-gitlab-framework"
	"github.com/opensourceways/community-robot-lib/secret"
	"github.com/sirupsen/logrus"
)

type options struct {
	service       liboptions.ServiceOptions
	gitlab        liboptions.GitLabOptions
	cacheEndpoint string
	maxRetries    int
}

func (o *options) Validate() error {
	if err := o.service.Validate(); err != nil {
		return err
	}

	return o.gitlab.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.gitlab.AddFlags(fs)
	o.service.AddFlags(fs)

	fs.StringVar(&o.cacheEndpoint, "cache-endpoint", "", "The endpoint of repo file cache")
	fs.IntVar(&o.maxRetries, "max-retries", 3, "The number of failed retry attempts to call the cache api")
	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit(botName)

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	secretAgent := new(secret.Agent)
	if err := secretAgent.Start([]string{o.gitlab.TokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error starting secret agent.")
	}

	defer secretAgent.Stop()

	//c, err := gitlabclient.NewGitlabClient(string(secretAgent.GetTokenGenerator(o.gitlab.TokenPath)()), "https://source.openeuler.sh/api/v4")
	// c, err := gitlabclient.NewGitlabClient(secretAgent.GetTokenGenerator(o.gitlab.TokenPath), "https://source.openeuler.sh/api/v4")
	c := gitlabclient.NewGitlabClient(secretAgent.GetTokenGenerator(o.gitlab.TokenPath), "https://source.openeuler.sh/api/v4")
	//if err != nil {
	//	logrus.WithError(err).Fatal("Invalid client")
	//}

	r := newRobot(c)

	framework.Run(r, o.service)
}
