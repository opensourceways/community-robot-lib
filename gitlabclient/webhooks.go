package gitlabclient

import (
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ValidateWebhook ensures that the provided request conforms to the
// format of a GitLab webhook and the payload can be validated with
// the provided hmac secret. It returns the event type, the event guid,
// the payload of the request, whether the webhook is valid or not,
// and finally the resultant HTTP status code
func ValidateWebhook(
	w http.ResponseWriter,
	r *http.Request,
	tokenGenerator func() string,
) (eType string, eventGUID string, ua string, payload []byte, ok bool, status int) {
	defer r.Body.Close()
	// Header checks: It must be a POST with an event type and a signature.
	if r.Method != http.MethodPost {
		status = http.StatusMethodNotAllowed
		responseHTTPError(w, status, "405 Method not allowed")

		return
	}

	if eType = r.Header.Get("X-Gitlab-Event"); eType == "" {
		status = http.StatusBadRequest
		responseHTTPError(w, status, "400 Bad Request: Missing X-Gitlab-Event Header")

		return
	}

	if ua = r.Header.Get("User-Agent"); ua == "" {
		status = http.StatusBadRequest
		responseHTTPError(w, status, "400 Bad Request: Missing User-Agent Header")

		return
	}

	if eventGUID = r.Header.Get("X-Gitlab-Event-UUID"); eventGUID == "" {
		status = http.StatusBadRequest
		responseHTTPError(w, status, "400 Bad Request: Missing X-Gitlab-Event-UUID Header")

		return
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		status = http.StatusBadRequest
		responseHTTPError(
			w, status,
			"400 Bad Request: Hook only accepts Content-Type: application/json - please reconfigure this hook on Gitlab",
		)

		return
	}

	sig := r.Header.Get("X-Gitlab-Token")
	if sig == "" {
		status = http.StatusForbidden
		responseHTTPError(w, status, "403 Forbidden: Missing X-Gitlab-Token")
		return
	}

	// Validate the payload with our HMAC secret.
	if sig != tokenGenerator() {
		status = http.StatusForbidden
		responseHTTPError(w, status, "403 Forbidden: Invalid X-Gitlab-Token")
		return
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status = http.StatusInternalServerError
		responseHTTPError(w, status, "500 Internal Server Error: Failed to read request body")
		return
	}

	status = http.StatusOK
	ok = true

	return
}

func responseHTTPError(w http.ResponseWriter, statusCode int, response string) {
	logrus.WithFields(
		logrus.Fields{
			"response":    response,
			"status-code": statusCode,
		},
	).Debug(response)

	http.Error(w, response, statusCode)
}
