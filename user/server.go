package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/stratumn/elk-sample-go-services/store"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/module/apmhttprouter"
)

// Server is the http server for the User service.
type Server struct {
	db *store.Store
}

// Start the user service.
func Start(port int, dbURL string) {
	db := store.New(dbURL)
	router := apmhttprouter.New()

	server := &Server{db: db}

	router.GET("/user/:id", server.User)
	router.GET("/users", server.Users)
	router.POST("/users", server.AddUser)

	log.Infof("Starting HTTP server on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), apmhttp.Wrap(router)))
}

// User returns the user with the given id.
func (s *Server) User(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	span, ctx := apm.StartSpan(r.Context(), "server/user", "app.request")
	defer span.End()

	apm.TransactionFromContext(ctx).Context.SetUserID(ps.ByName("id"))

	userID, err := strconv.ParseInt(ps.ByName("id"), 0, 0)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUser(ctx, int(userID))
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apm.TransactionFromContext(ctx).Context.SetUsername(user.Name)

	b, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// Users returns the list of all users.
func (s *Server) Users(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	span, ctx := apm.StartSpan(r.Context(), "server/users", "app.request")
	defer span.End()

	users, err := s.db.GetUsers(ctx)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apm.TransactionFromContext(ctx).Context.SetTag("userCount", fmt.Sprintf("%d", len(users)))

	b, _ := json.Marshal(users)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// AddUser adds a new user.
func (s *Server) AddUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	span, ctx := apm.StartSpan(r.Context(), "server/addUser", "app.request")
	defer span.End()

	var input store.User
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.db.AddUser(ctx, input.Name)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userCount.Inc()
	apm.TransactionFromContext(ctx).Context.SetUserID(fmt.Sprintf("%d", user.ID))
	apm.TransactionFromContext(ctx).Context.SetUsername(user.Name)

	b, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
