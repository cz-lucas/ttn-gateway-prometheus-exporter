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

	// Prometheus
	InitPrometheus()

	// Serve Prometheus metrics at :2112/metrics
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Prometheus metrics available at :2112/metrics")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	// Read interval
	intervalInSeconds, err := strconv.Atoi(os.Getenv("READ_INTERVAL"))

	if err != nil {
		log.Fatal("Invalid READ_INTERVAL:", err)
	}

	ticker := time.NewTicker(time.Duration(intervalInSeconds) * time.Second)
	defer ticker.Stop()

	// Main loop
	for {
		select {
		case <-ticker.C:
			start := time.Now()

			// Simulate API call
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
}
