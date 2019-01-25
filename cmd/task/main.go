package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/stratumn/elk-sample-go-services/task"
)

// Constants.
const (
	DBConnectionString = "postgres://postgres:postgres@localhost:5432/elk?sslmode=disable"
)

func main() {
	log.Info("Task service starting...")
	task.Start(DBConnectionString)
	log.Info("Task service stopping...")
}
