package urls

import (
	"fmt"
	"net/http"

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
func NewController(mongoDB *mgo.Database, shortProtocol string) Controller {
	service := newService(mongoDB, shortProtocol)
	return &_Controller{service: service}
}

type _Controller struct {
	service *_Service
}

func (controller *_Controller) Create(rw http.ResponseWriter, r *http.Request) {
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
