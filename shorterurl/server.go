package shorterurl

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/garyburd/redigo/redis"

	mgo "gopkg.in/mgo.v2"
)

// Server serves up an API for creating and resolving
// short URLs
type Server interface {
	// Run starts listing for incoming requests. This
	// function will run until the server exits due to
	// an error.
	Run() error
}

// New constructs a new Server that listens for incoming
// HTTP requests
func New(auth, mongoDBURL string, port int, redisNamespace, redisURL, shortProtocol string) Server {
	return &HTTPServer{
		auth:           auth,
		mongoDBURL:     mongoDBURL,
		port:           port,
		redisNamespace: redisNamespace,
		redisURL:       redisURL,
		shortProtocol:  shortProtocol,
	}
}

// HTTPServer implements Server and listens for incoming
// HTTP requests
type HTTPServer struct {
	auth           string
	mongoDBURL     string
	port           int
	redisNamespace string
	redisURL       string
	shortProtocol  string
}

// Run starts listing for incoming requests. This
// function will run until the server exits due to
// an error.
func (server *HTTPServer) Run() error {
	mongo, err := mgo.Dial(server.mongoDBURL)
	if err != nil {
		return err
	}
	mongoDB := mongo.DB(mongoDatabaseName(server.mongoDBURL))

	redisConn, err := redis.DialURL(server.redisURL)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf(":%v", server.port)
	router := newRouter(mongoDB, redisConn, server.redisNamespace, server.shortProtocol)
	return http.ListenAndServe(addr, router)
}

func mongoDatabaseName(mongoDBURL string) string {
	parts := strings.Split(mongoDBURL, "/")
	return parts[len(parts)-1]
}
