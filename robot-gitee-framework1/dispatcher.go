package framework

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type dispatcher struct {
	h map[string]func([]byte, *logrus.Entry)

	// Tracks running handlers for graceful shutdown
	wg sync.WaitGroup
}

func (d *dispatcher) wait() {
	d.wg.Wait() // Handle remaining requests
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventType, eventGUID, payload, ok := parseRequest(w, r)
	if !ok {
		return
	}

	handle, ok := d.h[eventType]
	if !ok {
		return
	}

	l := logrus.WithFields(
		logrus.Fields{
			"event-type": eventType,
			"event_id":   eventGUID,
		},
	)

	d.wg.Add(1)

	go func() {
		handle(payload, l)

		d.wg.Done()
	}()
}

func parseRequest(w http.ResponseWriter, r *http.Request) (
	eventType string,
	uuid string,
	payload []byte,
	ok bool,
) {
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
