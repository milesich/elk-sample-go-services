package task

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	taskCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_backend_task_count",
			Help: "the total number of tasks created",
		},
		[]string{"userId"},
	)
)
