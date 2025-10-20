package main

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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

// InitPrometheus returns a custom registry
func InitPrometheus() *prometheus.Registry {
	enableRuntimeMetrics := os.Getenv("ENABLE_RUNTIME_METRICS") == "true"

	// Create a new custom registry
	reg := prometheus.NewRegistry()

	// Register your app's custom metrics
	reg.MustRegister(apiCallsTotal)
	reg.MustRegister(apiCallFailures)
	reg.MustRegister(lastApiCallDuration)

	if enableRuntimeMetrics {
		reg.MustRegister(collectors.NewGoCollector())
		reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	return reg
}
