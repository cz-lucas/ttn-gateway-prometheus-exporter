package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
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

	averageApiCallDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "api_call_duration_seconds",
			Help: "The average duration of the API calls to the TTN in seconds",
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
		reg.MustRegister(averageApiCallDuration)
	}

	if enableRuntimeMetrics {
		reg.MustRegister(collectors.NewGoCollector())
		reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	return reg
}
