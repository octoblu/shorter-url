package cachedurls

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"

	mgo "gopkg.in/mgo.v2"
)

// Controller defines a controller that can retrieve
// short urls
type Controller interface {
	// Get retrieves a short url and redirects the user
	Get(rw http.ResponseWriter, r *http.Request)
}

// NewController returns a new controller instance
// for retrieving URLs
func NewController(cache redis.Conn, mongoDB *mgo.Database, redisNamespace, shortProtocol string) Controller {
	service := newService(cache, mongoDB, redisNamespace, shortProtocol)
	return &_Controller{service: service}
}

type _Controller struct {
	service *_Service
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
