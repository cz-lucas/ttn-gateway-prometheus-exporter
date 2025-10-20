package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	godotenv.Load(".env")

	if !keyExistsInConfig("TTN_GATEWAY_ID") {
		log.Fatalln("The TTN_GATEWAY_ID is not configured")
		os.Exit(-1)
	}

	if !keyExistsInConfig("TTN_API_KEY") {
		log.Fatalln("The TTN_API_KEY is not configured")
		os.Exit(-1)
	}

	// Get URL
	var ttnRequestUrl = getEnvString("TTN_BASE_URL", "https://eu1.cloud.thethings.network/api/v3/gs/gateways/")
	ttnRequestUrl += os.Getenv("TTN_GATEWAY_ID")
	ttnRequestUrl += getEnvString("TTN_URL_STATS_SUFFIX", "/connection/stats")
	var apiService = NewTTNApiService(ttnRequestUrl, os.Getenv("TTN_API_KEY"))

	// HTTP Server
	var addr = getEnvString("ADDRESS", ":9000")
	httpService := NewHttpService(addr)

	// Register the /metrics endpoint
	var enableRuntimeMetrics, _ = getEnvBool("ENABLE_RUNTIME_METRICS", true)
	var enableAppMetrics, _ = getEnvBool("ENABLE_APP_METRICS", true)
	reg := InitPrometheus(enableRuntimeMetrics, enableAppMetrics)
	httpService.RegisterRoute("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	// You can register more routes here, e.g. health checks
	httpService.RegisterRoute("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	// Start the HTTP service
	httpService.Start()

	// Read interval
	intervalInSeconds, _ := getEnvInt("READ_INTERVAL", 600)

	ticker := time.NewTicker(time.Duration(intervalInSeconds) * time.Second)
	defer ticker.Stop()

	// Main loop
	for range ticker.C {
		//start := time.Now()

		response, err := apiService.Get()
		apiCallsTotal.Inc()

		if err != nil {
			apiCallFailures.Inc()
		} else {
			log.Println(response)
		}

		//duration := time.Since(start).Seconds()
		//lastApiCallDuration.Set(duration)
	}
}
