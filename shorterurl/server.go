package shorterurl

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
func New(auth, mongoDBURL string, port int, redisNamespace, redisURL string) Server {
	return &HTTPServer{
		auth:           auth,
		mongoDBURL:     mongoDBURL,
		port:           port,
		redisNamespace: redisNamespace,
		redisURL:       redisURL,
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
}

// Run starts listing for incoming requests. This
// function will run until the server exits due to
// an error.
func (server *HTTPServer) Run() error {
	return nil
}
