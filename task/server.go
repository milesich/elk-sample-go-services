package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/stratumn/elk-sample-go-services/store"
)

// Server is the http server for the Task service.
type Server struct {
	router *httprouter.Router
	db     *store.Store
}

// Start the task service.
func Start(connStr string) {
	db := store.New(connStr)
	router := httprouter.New()

	server := &Server{router: router, db: db}

	router.GET("/user/:userId/tasks", server.Tasks)
	router.POST("/user/:userId/tasks", server.AddTask)
	router.POST("/user/:userId/task/:taskId", server.UpdateTask)

	log.Infof("Starting HTTP server on port %d", 4002)
	log.Fatal(http.ListenAndServe(":4002", router))
}

// Tasks returns the user's tasks.
func (s *Server) Tasks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Infof("GET user %s tasks", ps.ByName("userId"))

	userID, err := strconv.ParseInt(ps.ByName("userId"), 0, 0)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks, err := s.db.GetUserTasks(r.Context(), int(userID))
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(tasks)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// AddTask adds a new task.
func (s *Server) AddTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Infof("POST add task to user %s", ps.ByName("userId"))

	userID, err := strconv.ParseInt(ps.ByName("userId"), 0, 0)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input store.Task
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&input)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := s.db.AddTask(r.Context(), int(userID), input.Name)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(task)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// UpdateTask updates a given task.
func (s *Server) UpdateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Infof("POST update task %s for user %s", ps.ByName("taskId"), ps.ByName("userId"))

	_, err := strconv.ParseInt(ps.ByName("userId"), 0, 0)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskID, err := strconv.ParseInt(ps.ByName("taskId"), 0, 0)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input store.Task
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&input)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := s.db.UpdateTask(r.Context(), int(taskID), input.Name, input.Done)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(task)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}