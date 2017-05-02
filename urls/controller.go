package urls

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"

	mgo "gopkg.in/mgo.v2"
)

// Controller defines a controller that can create
// short urls
type Controller interface {
	// Create stores a new short url in the database
	Create(rw http.ResponseWriter, r *http.Request)

	// Delete removes a short url from the database
	// and cache
	Delete(rw http.ResponseWriter, r *http.Request)

	// Get retrieves a short url and redirects the user
	// This is the only endpoint that does not require
	// authentication
	Get(rw http.ResponseWriter, r *http.Request)
}

// NewController returns a new controller instance
// for managing urls
func NewController(auth string, cache redis.Conn, mongoSession *mgo.Session, redisNamespace, shortProtocol string) Controller {
	service := newService(cache, mongoSession, redisNamespace, shortProtocol)
	return &_Controller{auth: auth, service: service}
}

type _Controller struct {
	auth    string
	service *_Service
}

func (controller *_Controller) Create(rw http.ResponseWriter, r *http.Request) {
	if !controller.authenticated(r.BasicAuth()) {
		http.Error(rw, "Unauthorized", 401)
		return
	}

	createBody, err := parseCreateBody(r.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Could not parse request body: %v", err.Error()), 422)
		return
	}

	shorterURL, err := controller.service.Create(createBody.LongURL, r.Host)
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

func (controller *_Controller) Delete(rw http.ResponseWriter, r *http.Request) {
	if !controller.authenticated(r.BasicAuth()) {
		http.Error(rw, "Unauthorized", 401)
		return
	}

	token := mux.Vars(r)["token"]
	err := controller.service.Delete(r.Host, token)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to delete shorterUrl: %v", err.Error()), 500)
		return
	}

	rw.WriteHeader(204)
}

func (controller *_Controller) Get(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	longURL, err := controller.service.GetLongURL(r.Host, token)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to retrieve longUrl: %v", err.Error()), 500)
		return
	}
	http.Redirect(rw, r, longURL, http.StatusPermanentRedirect)
}

func (controller *_Controller) authenticated(username, password string, authPresent bool) bool {
	if !authPresent {
		return false
	}

	parts := strings.Split(controller.auth, ":")
	return username == parts[0] && password == parts[1]
}
