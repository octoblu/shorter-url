package shorterurl

import (
	"crypto/tls"
	"fmt"
	"net"
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
	Run(onListen func()) error
}

// New constructs a new Server that listens for incoming
// HTTP requests
func New(auth, mongoDBURL string, port int, redisNamespace, redisURL, shortProtocol, version string) Server {
	return &HTTPServer{
		auth:           auth,
		mongoDBURL:     mongoDBURL,
		port:           port,
		redisNamespace: redisNamespace,
		redisURL:       redisURL,
		shortProtocol:  shortProtocol,
		version:        version,
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
	version        string
}

// Run starts listing for incoming requests. This
// function will run until the server exits due to
// an error.
func (server *HTTPServer) Run(onListen func()) error {
	mongoDB, err := server.dialMongo()
	if err != nil {
		return err
	}

	redisConn, err := redis.DialURL(server.redisURL)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf(":%v", server.port)
	router := newRouter(server.auth, mongoDB, redisConn, server.redisNamespace, server.shortProtocol, server.version)
	onListen()
	return http.ListenAndServe(addr, router)
}

func (server *HTTPServer) dialMongo() (*mgo.Database, error) {
	if !strings.HasSuffix(server.mongoDBURL, "?ssl=true") {
		mongo, err := mgo.Dial(server.mongoDBURL)
		if err != nil {
			return nil, err
		}
		mongoDB := mongo.DB(mongoDatabaseName(server.mongoDBURL))
		return mongoDB, err
	}

	mongoDBURL := strings.TrimSuffix(server.mongoDBURL, "?ssl=true")

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	dialInfo, err := mgo.ParseURL(mongoDBURL)
	if err != nil {
		return nil, err
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	defer session.Close()

	mongoDB := session.DB(mongoDatabaseName(mongoDBURL))
	return mongoDB, nil
}

func mongoDatabaseName(mongoDBURL string) string {
	parts := strings.Split(mongoDBURL, "/")
	return parts[len(parts)-1]
}
