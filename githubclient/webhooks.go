package githubclient

import (
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ValidateWebhook ensures that the provided request conforms to the
// format of a GitHub webhook and the payload can be validated with
// the provided hmac secret. It returns the event type, the event guid,
// the payload of the request, whether the webhook is valid or not,
// and finally the resultant HTTP status code
func ValidateWebhook(
	w http.ResponseWriter,
	r *http.Request,
	tokenGenerator func() []byte,
) (eType string, guid string, payload []byte, ok bool, status int) {
	defer r.Body.Close()
	// Header checks: It must be a POST with an event type and a signature.
	if r.Method != http.MethodPost {
		status = http.StatusMethodNotAllowed
		responseHTTPError(w, status, "405 Method not allowed")

		return
	}

	if eType = r.Header.Get("X-GitHub-Event"); eType == "" {
		status = http.StatusBadRequest
		responseHTTPError(w, status, "400 Bad Request: Missing X-GitHub-Event Header")

		return
	}

	if guid = r.Header.Get("X-GitHub-Delivery"); guid == "" {
		status = http.StatusBadRequest
		responseHTTPError(w, status, "400 Bad Request: Missing X-GitHub-Delivery Header")

		return
	}

	sig := r.Header.Get("X-Hub-Signature")
	if sig == "" {
		status = http.StatusForbidden
		responseHTTPError(w, status, "403 Forbidden: Missing X-Hub-Signature")
		return
	}

	if contentType := r.Header.Get("content-type"); contentType != "application/json" {
		status = http.StatusBadRequest
		responseHTTPError(
			w, status,
			"400 Bad Request: Hook only accepts content-type: application/json - please reconfigure this hook on GitHub",
		)

		return
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status = http.StatusInternalServerError
		responseHTTPError(w, status, "500 Internal Server Error: Failed to read request body")
		return
	}

	// Validate the payload with our HMAC secret.
	if !ValidatePayload(payload, sig, tokenGenerator) {
		status = http.StatusForbidden
		responseHTTPError(w, status, "403 Forbidden: Invalid X-Hub-Signature")

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
