package shorterurl

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
)

func newRouter(mongoDB *mgo.Database) http.Handler {
	router := mux.NewRouter()
	return router
}
