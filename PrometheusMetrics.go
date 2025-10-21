package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	// App metrics
	apiCallsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "api_calls_total",
			Help: "Total number of API calls made",
		},
	)

	apiCallFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "api_call_failures_total",
			Help: "Total number of failed API calls",
		},
	)

	lastApiCallDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "last_api_call_duration_seconds",
			Help: "The duration of the last API call to the TTN in seconds",
		},
	)
)

// InitPrometheus returns a custom registry
func InitPrometheus(enableRuntimeMetrics bool, enableAppMetrics bool) *prometheus.Registry {
	// Create a new custom registry
	reg := prometheus.NewRegistry()

	if enableAppMetrics {
		// Register your app's custom metrics
		reg.MustRegister(apiCallsTotal)
		reg.MustRegister(apiCallFailures)
		reg.MustRegister(lastApiCallDuration)
	}

	if enableRuntimeMetrics {
		reg.MustRegister(collectors.NewGoCollector())
		reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	return reg
}
