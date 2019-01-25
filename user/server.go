package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/stratumn/elk-sample-go-services/store"
)

// Server is the http server for the User service.
type Server struct {
	router *httprouter.Router
	db     *store.Store
}

// Start the user service.
func Start(connStr string) {
	db := store.New(connStr)
	router := httprouter.New()

	server := &Server{router: router, db: db}

	router.GET("/user/:id", server.User)
	router.GET("/users", server.Users)
	router.POST("/users", server.AddUser)

	log.Infof("Starting HTTP server on port %d", 4001)
	log.Fatal(http.ListenAndServe(":4001", router))
}

// User returns the user with the given id.
func (s *Server) User(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Infof("GET user %s", ps.ByName("id"))

	userID, err := strconv.ParseInt(ps.ByName("id"), 0, 0)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUser(r.Context(), int(userID))
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// Users returns the list of all users.
func (s *Server) Users(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("GET users")

	users, err := s.db.GetUsers(r.Context())
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(users)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// AddUser adds a new user.
func (s *Server) AddUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("POST users")

	var input store.User
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.db.AddUser(r.Context(), input.Name)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
