package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/stratumn/elk-sample-go-services/user"
)

func main() {
	log.Info("User service starting...")
	user.Start()
	log.Info("User service stopping...")
}
