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

	// Gateway stats
	numberOfDownlinkMessages = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gw_number_of_downlink_messages",
			Help: "The total number of downlink messages",
		},
		[]string{"gateway_id"}, // Add gateway_id as a label
	)

	numberOfUplinkMessages = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gw_number_of_uplink_messages",
			Help: "The total number of uplink messages",
		},
		[]string{"gateway_id"}, // Add gateway_id as a label
	)

	rtt_min = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gw_rtt_min",
			Help: "The minimal round trip time in ms",
		},
		[]string{"gateway_id"}, // Add gateway_id as a label
	)

	rtt_median = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gw_rtt_median",
			Help: "The median round trip time in ms",
		},
		[]string{"gateway_id"}, // Add gateway_id as a label
	)

	rtt_max = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gw_rtt_may",
			Help: "The maximal round trip time in ms",
		},
		[]string{"gateway_id"}, // Add gateway_id as a label
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

	// Register gateway metrics
	reg.MustRegister(numberOfDownlinkMessages)
	reg.MustRegister(numberOfUplinkMessages)
	reg.MustRegister(rtt_min)
	reg.MustRegister(rtt_median)
	reg.MustRegister(rtt_max)

	return reg
}
