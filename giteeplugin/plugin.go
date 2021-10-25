package giteeplugin

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/community-robot-lib/config"
	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/opensourceways/community-robot-lib/options"
)

type HandlerRegitster interface {
	RegisterIssueHandler(IssueHandler)
	RegisterPullRequestHandler(PullRequestHandler)
	RegisterPushEventHandler(PushEventHandler)
	RegisterNoteEventHandler(NoteEventHandler)
}

type Plugin interface {
	NewPluginConfig() config.PluginConfig
	RegisterEventHandler(HandlerRegitster)
}

func Run(p Plugin, o options.PluginOptions) {
	agent := config.NewConfigAgent(p.NewPluginConfig)
	if err := agent.Start(o.PluginConfig); err != nil {
		return
	}

	h := handlers{}
	p.RegisterEventHandler(&h)

	d := &dispatcher{agent: &agent, h: h}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.OnInterrupt(func() {
		agent.Stop()
		d.Wait()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	http.Handle("/gitee-hook", d)

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(o.Port)}

	interrupts.ListenAndServe(httpServer, o.GracePeriod)
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventType, eventGUID, payload, ok := parseRequest(w, r)

	if !ok {
		return
	}

	l := logrus.WithFields(
		logrus.Fields{
			"event-type": eventType,
			"event_id":   eventGUID,
		},
	)

	if err := d.Dispatch(eventType, payload, l); err != nil {
		l.WithError(err).Error()
	}
}

func parseRequest(w http.ResponseWriter, r *http.Request) (eventType string, uuid string, payload []byte, ok bool) {
	defer r.Body.Close()

	resp := func(code int, msg string) {
		http.Error(w, msg, code)
	}

	if r.Header.Get("User-Agent") != "Robot-Gitee-Access" {
		resp(http.StatusBadRequest, "400 Bad Request: unknown User-Agent Header")
		return
	}

	if eventType = r.Header.Get("X-Gitee-Event"); eventType == "" {
		resp(http.StatusBadRequest, "400 Bad Request: Missing X-Gitee-Event Header")
		return
	}

	if uuid = r.Header.Get("X-Gitee-Timestamp"); uuid == "" {
		resp(http.StatusBadRequest, "400 Bad Request: Missing X-Gitee-Timestamp Header")
		return
	}

	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp(http.StatusInternalServerError, "500 Internal Server Error: Failed to read request body")
		return
	}
	payload = v
	ok = true

	return
}
