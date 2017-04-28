package urls

import (
	"fmt"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"
)

// Controller defines a controller that can create
// short urls
type Controller interface {
	// Create stores a new short url in the database
	Create(rw http.ResponseWriter, r *http.Request)
}

// NewController returns a new controller instance
// for managing urls
func NewController(auth string, mongoDB *mgo.Database, shortProtocol string) Controller {
	service := newService(mongoDB, shortProtocol)
	return &_Controller{auth: auth, service: service}
}

type _Controller struct {
	auth    string
	service *_Service
}

func (controller *_Controller) Create(rw http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	parts := strings.Split(controller.auth, ":")
	if !ok || username != parts[0] || password != parts[1] {
		http.Error(rw, "Unauthorized", 401)
		return
	}

	createBody, err := parseCreateBody(r.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Could not parse request body: %v", err.Error()), 422)
		return
	}

	shorterURL, err := controller.service.Create(createBody.LongURL, r.URL.Scheme, r.Host)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to create shorterUrl: %v", err.Error()), 500)
		return
	}

	response, err := formatCreateResponse(shorterURL)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to generate response: %v", err.Error()), 500)
		return
	}

	rw.WriteHeader(201)
	rw.Write(response)
}
