package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	godotenv.Load(".env")

	// Get URL
	var ttnRequestUrl = os.Getenv("TTN_BASE_URL") + os.Getenv("TTN_GATEWAY_ID") + os.Getenv("TTN_URL_STATS_SUFFIX")
	var apiService = NewTTNApiService(ttnRequestUrl, os.Getenv("TTN_API_KEY"))

	// HTTP Server
	httpService := NewHttpService(":2112")

	// Register the /metrics endpoint
	reg := InitPrometheus()
	httpService.RegisterRoute("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	// You can register more routes here, e.g. health checks
	httpService.RegisterRoute("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	// Start the HTTP service
	httpService.Start()

	// Prometheus
	InitPrometheus()

	// Read interval
	intervalInSeconds, err := strconv.Atoi(os.Getenv("READ_INTERVAL"))
	if err != nil {
		log.Fatal("Invalid READ_INTERVAL:", err)
	}

	ticker := time.NewTicker(time.Duration(intervalInSeconds) * time.Second)
	defer ticker.Stop()

	// Main loop
	for range ticker.C {
		start := time.Now()

		response, err := apiService.Get()
		apiCallsTotal.Inc()

		if err != nil {
			apiCallFailures.Inc()
		} else {
			log.Println(response)
		}

		duration := time.Since(start).Seconds()
		lastApiCallDuration.Set(duration)
	}
}
