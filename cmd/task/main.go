package main

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/stratumn/elk-sample-go-services/task"
)

// Constants.
const (
	DefaultDBURL = "postgres://postgres:postgres@localhost:5432/elk?sslmode=disable"
	DefaultPort  = 4002
)

func main() {
	port := DefaultPort
	portStr, ok := os.LookupEnv("PORT")
	if ok {
		parsedPort, err := strconv.ParseInt(portStr, 0, 0)
		if err != nil {
			log.WithError(err).Fatal("Invalid port")
		}

		port = int(parsedPort)
	}

	dbURL, ok := os.LookupEnv("DB_URL")
	if !ok {
		dbURL = DefaultDBURL
	}

	log.Info("Task service starting...")
	task.Start(port, dbURL)
	log.Info("Task service stopping...")
}
