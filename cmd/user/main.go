package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/stratumn/elk-sample-go-services/user"
)

// Constants.
const (
	DBConnectionString = "postgres://postgres:postgres@localhost:5432/elk?sslmode=disable"
)

func main() {
	log.Info("User service starting...")
	user.Start(DBConnectionString)
	log.Info("User service stopping...")
}
