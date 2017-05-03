package shorterurl

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/octoblu/shorter-url/serverinfo"
	"github.com/octoblu/shorter-url/urls"
)

func newRouter(auth string, redisPool *redis.Pool, mongoSession *mgo.Session, redisNamespace, shortProtocol, version string) http.Handler {
	serverInfoController := serverinfo.New(version)
	urlsController := urls.NewController(auth, redisPool, mongoSession, redisNamespace, shortProtocol)

	router := mux.NewRouter()

	router.Methods("GET").Path("/healthcheck").HandlerFunc(serverInfoController.Healthcheck)
	router.Methods("GET").Path("/version").HandlerFunc(serverInfoController.Version)

	router.Methods("POST").Path("/").HandlerFunc(urlsController.Create)
	router.Methods("DELETE").Path("/{token}").HandlerFunc(urlsController.Delete)
	router.Methods("GET").Path("/{token}").HandlerFunc(urlsController.Get)

	return router
}
