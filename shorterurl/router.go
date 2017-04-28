package shorterurl

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/octoblu/shorter-url/cachedurls"
	"github.com/octoblu/shorter-url/urls"
)

func newRouter(mongoDB *mgo.Database, shortProtocol string) http.Handler {
	urlsController := urls.NewController(mongoDB, shortProtocol)
	cachedUrlsController := cachedurls.NewController(mongoDB, shortProtocol)

	router := mux.NewRouter()
	router.Methods("POST").Path("/api/urls").HandlerFunc(urlsController.Create)
	router.Methods("GET").Path("/{token}").HandlerFunc(cachedUrlsController.Get)
	return router
}
