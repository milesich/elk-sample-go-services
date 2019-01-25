package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/stratumn/elk-sample-go-services/task"
)

func main() {
	log.Info("Task service starting...")
	task.Start()
	log.Info("Task service stopping...")
}
