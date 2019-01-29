package user

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	userCount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_backend_user_count",
			Help: "the total number of users created",
		},
	)
)
