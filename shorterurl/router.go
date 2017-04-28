package shorterurl

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/octoblu/shorter-url/cachedurls"
	"github.com/octoblu/shorter-url/urls"
)

func newRouter(auth string, cache redis.Conn, mongoDB *mgo.Database, redisNamespace, shortProtocol, version string) http.Handler {
	urlsController := urls.NewController(auth, cache, mongoDB, redisNamespace, shortProtocol)
	cachedUrlsController := cachedurls.NewController(cache, mongoDB, redisNamespace, shortProtocol)

	router := mux.NewRouter()

	router.Methods("GET").Path("/healthcheck").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("{\"online\":true}"))
	})
	router.Methods("GET").Path("/version").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("{\"version\":\"%v\"}", version)))
	})

	router.Methods("POST").Path("/").HandlerFunc(urlsController.Create)
	router.Methods("DELETE").Path("/{token}").HandlerFunc(urlsController.Delete)
	router.Methods("GET").Path("/{token}").HandlerFunc(cachedUrlsController.Get)

	return router
}
