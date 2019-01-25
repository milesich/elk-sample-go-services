package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// Server is the http server for the User service.
type Server struct {
	router *httprouter.Router
}

// Start the user service.
func Start() {
	router := httprouter.New()
	server := &Server{router: router}

	router.GET("/user/:id", server.User)
	router.GET("/users", server.Users)
	router.POST("/users", server.AddUser)

	log.Infof("Starting HTTP server on port %d", 4001)
	log.Fatal(http.ListenAndServe(":4001", router))
}

// User returns the user with the given id.
func (s *Server) User(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	log.Infof("GET user %s", ps.ByName("id"))

	userID, err := strconv.ParseInt(ps.ByName("id"), 0, 0)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, _ := json.Marshal(&User{
		ID:   int(userID),
		Name: "alice",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(user)
}

// Users returns the list of all users.
func (s *Server) Users(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	log.Info("GET users")

	users := []User{
		User{
			ID:   1,
			Name: "alice",
		},
		User{
			ID:   2,
			Name: "bob",
		},
	}

	b, _ := json.Marshal(users)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// AddUser adds a new user.
func (s *Server) AddUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("POST users")

	var user User
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&user)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = 42
	b, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
