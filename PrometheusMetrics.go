package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	enableRuntimeMetrics := os.Getenv("ENABLE_RUNTIME_METRICS") == "true"

	// Create a new custom registry
	reg := prometheus.NewRegistry()

	// Register your app's custom metrics
	reg.MustRegister(apiCallsTotal)
	reg.MustRegister(apiCallFailures)
	reg.MustRegister(lastApiCallDuration)

	if enableRuntimeMetrics {
		// Conditionally register Go runtime and process metrics
		reg.MustRegister(collectors.NewGoCollector())
		reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	// Use the custom registry in the handler
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	go func() {
		log.Println("Metrics server started at :2112/metrics")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()
}
