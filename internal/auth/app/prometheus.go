package auth_app

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "path"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "request_duration",
			Help: "Duration of HTTP request",
		},
		[]string{"service", "path"},
	)

	statusResponse = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "status_response_count",
			Help: "Status of HTTP response",
		},
		[]string{"service", "path", "status"},
	)
)
