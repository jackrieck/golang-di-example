package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jrieck1991/golang-di-example/internal/storage"
)

type Server interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type server struct {
	router      http.Handler
	storeClient storage.Client
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// New initializes a new Server, support for dependency injection
func New(opts ...func(*server)) Server {

	s := &server{}

	router := mux.NewRouter()
	router.HandleFunc("/", s.index).Methods(http.MethodGet)
	router.HandleFunc("/get", s.get).Methods(http.MethodGet)

	s.router = router

	for _, opt := range opts {
		opt(s)
	}

	if s.storeClient == nil {
		client := storage.New()
		s.storeClient = client
	}

	return s
}

// WithStorageClient used for dependency injection
func WithStorageClient(client storage.Client) func(*server) {
	return func(s *server) {
		s.storeClient = client
	}
}

// serves the index page
func (s *server) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index page\n"))
}

// get returns a value matching a key or a 400
func (s *server) get(w http.ResponseWriter, r *http.Request) {

	// get query parameters from url
	if err := r.ParseForm(); err != nil {
		log.Println("error parsing query parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get key parameter
	searchKey := r.FormValue("key")
	if searchKey == "" {
		log.Println("query parameter 'key' not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ask storage layer for value
	value, err := s.storeClient.Get(searchKey)
	if err != nil {
		log.Printf("value not found for key: %s", searchKey)
		w.WriteHeader(http.StatusBadRequest)
	}

	// return value to client
	w.Write([]byte(fmt.Sprintf("key: %s, item: %v\n", searchKey, value)))
}
