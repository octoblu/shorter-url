package serverinfo

import (
	"fmt"
	"net/http"
)

// Controller provides a few server information
// endpoints, including healthcheck & version
type Controller interface {
	// Healthcheck always responds with:
	// < HTTP/1.1 200 OK
	// < Content-Type application/json
	// < Content-Length 15
	// {"online":true}
	Healthcheck(rw http.ResponseWriter, r *http.Request)

	// Version always (adjusting for the current version) responds with:
	// < HTTP/1.1 200 OK
	// < Content-Type application/json
	// < Content-Length 15
	// {"version":"2.1.0"}
	Version(rw http.ResponseWriter, r *http.Request)
}

// New constructs a new controller instance
func New(version string) Controller {
	return &_Controller{version: version}
}

type _Controller struct {
	version string
}

func (controller *_Controller) Healthcheck(rw http.ResponseWriter, r *http.Request) {
	responseBody, err := formatHealthcheck()
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error formatting response: %v", err.Error()), 500)
	}

	rw.WriteHeader(200)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(responseBody)
}

func (controller *_Controller) Version(rw http.ResponseWriter, r *http.Request) {
	responseBody, err := formatVersion(controller.version)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error formatting response: %v", err.Error()), 500)
	}

	rw.WriteHeader(200)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(responseBody)
}
