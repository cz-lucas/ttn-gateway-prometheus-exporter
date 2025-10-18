package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	apiCallsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ttn_api_calls_total",
			Help: "Total number of TTN API calls made",
		},
	)

	apiCallFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ttn_api_call_failures_total",
			Help: "Total number of failed TTN API calls",
		},
	)

	lastApiCallDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ttn_api_call_duration_seconds",
			Help: "Duration of last TTN API call in seconds",
		},
	)
)

func InitPrometheus() {
	prometheus.MustRegister(apiCallsTotal)
	prometheus.MustRegister(apiCallFailures)
	prometheus.MustRegister(lastApiCallDuration)
}
