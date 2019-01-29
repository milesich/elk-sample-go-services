package store

import (
	"context"
	"database/sql"
	"time"

	// PQ is needed to interact with postgres.
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"go.elastic.co/apm/module/apmsql"
	// PQ drivers should be loaded for automatic instrumentation.
	_ "go.elastic.co/apm/module/apmsql/pq"
)

// Store provides access to the data store.
type Store struct {
	db *sql.DB

	addUserStmt    *sql.Stmt
	addTaskStmt    *sql.Stmt
	updateTaskStmt *sql.Stmt

	getUserStmt      *sql.Stmt
	getUsersStmt     *sql.Stmt
	getUserTasksStmt *sql.Stmt
}

// New connects to the database and makes sure it's setup.
func New(connStr string) *Store {
	var err error
	var db *sql.DB

	log.Infof("Setting up database (%s)...", connStr)

	for i := 0; i < 12; i++ {
		<-time.After(5 * time.Second)

		db, err = apmsql.Open("postgres", connStr)
		if err != nil {
			log.Errorf("Can't connect to DB: %s", err.Error())
			continue
		}

		_, err = db.Exec(SQLCreateSchema)
		if err != nil {
			log.Errorf("Can't create DB schema: %s", err.Error())
			continue
		}

		_, err = db.Exec(SQLCreateUserTable)
		if err != nil {
			log.Errorf("Can't create user table: %s", err.Error())
			continue
		}

		_, err = db.Exec(SQLCreateTaskTable)
		if err != nil {
			log.Errorf("Can't create task table: %s", err.Error())
			continue
		}
	}

	if err != nil {
		log.Fatal("Database setup failed.")
	}

	s := &Store{db: db}
	s.prepare()

	return s
}

func (s *Store) prepare() {
	prep := func(stmt string) *sql.Stmt {
		prepared, err := s.db.Prepare(stmt)
		if err != nil {
			log.WithError(err).WithField("SQL", stmt).Fatal("Could not prepare statement")
		}

		return prepared
	}

	s.addUserStmt = prep(SQLAddUser)
	s.addTaskStmt = prep(SQLAddTask)
	s.updateTaskStmt = prep(SQLUpdateTask)
	s.getUserStmt = prep(SQLGetUser)
	s.getUsersStmt = prep(SQLGetUsers)
	s.getUserTasksStmt = prep(SQLGetUserTasks)
}

// AddUser to the DB.
func (s *Store) AddUser(ctx context.Context, name string) (*User, error) {
	var userID int

	err := s.addUserStmt.QueryRow(name).Scan(&userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &User{
		ID:   userID,
		Name: name,
	}, nil
}

// GetUser from the DB.
func (s *Store) GetUser(ctx context.Context, userID int) (*User, error) {
	var ID int
	var name string

	err := s.getUserStmt.QueryRow(userID).Scan(&ID, &name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &User{
		ID:   ID,
		Name: name,
	}, nil
}

// GetUsers from the DB.
func (s *Store) GetUsers(ctx context.Context) ([]*User, error) {
	rows, err := s.getUsersStmt.Query()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var ID int
		var name string

		err = rows.Scan(&ID, &name)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		users = append(users, &User{ID: ID, Name: name})
	}

	return users, nil
}

// AddTask to the DB.
func (s *Store) AddTask(ctx context.Context, userID int, name string) (*Task, error) {
	var taskID int

	err := s.addTaskStmt.QueryRow(userID, name).Scan(&taskID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Task{
		ID:   taskID,
		Name: name,
		Done: false,
	}, nil
}

// UpdateTask updates a task in the DB.
func (s *Store) UpdateTask(ctx context.Context, taskID int, name string, done bool) (*Task, error) {
	_, err := s.updateTaskStmt.Exec(taskID, name, done)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Task{
		ID:   taskID,
		Name: name,
		Done: done,
	}, nil
}

// GetUserTasks gets all the tasks for the given user.
func (s *Store) GetUserTasks(ctx context.Context, userID int) ([]*Task, error) {
	rows, err := s.getUserTasksStmt.Query(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var ID int
		var userID int
		var name string
		var done bool

		err = rows.Scan(&ID, &userID, &name, &done)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		tasks = append(tasks, &Task{ID: ID, Name: name, Done: done})
	}

	return tasks, nil
}
