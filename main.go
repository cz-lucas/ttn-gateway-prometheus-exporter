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
	}

	if !keyExistsInConfig("TTN_API_KEY") {
		log.Fatalln("The TTN_API_KEY is not configured")
	}

	// Read interval
	intervalInSeconds, err := getEnvInt("READ_INTERVAL", 600)
	if err != nil {
		log.Fatalln("READ_INTERVAL is not a number")
	}

	var gatewayId string = os.Getenv("TTN_GATEWAY_ID")

	log.Printf("Starting TTN-Gateway-Prometheus-exporter\n")
	log.Printf("GatewayID: %s \n", gatewayId)

	// Get URL
	var ttnRequestUrl = getEnvString("TTN_BASE_URL", "https://eu1.cloud.thethings.network/api/v3/gs/gateways/")
	ttnRequestUrl += gatewayId
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

	ticker := time.NewTicker(time.Duration(intervalInSeconds) * time.Second)
	defer ticker.Stop()

	// Main loop
	for range ticker.C {
		start := time.Now()

		response, err := apiService.Get()
		apiCallsTotal.Inc()
		log.Println("Getting gateway-statistics")

		if err != nil {
			apiCallFailures.Inc()
			log.Printf("ERROR: Request to the TTN failed: %v", err)
			continue
		} else {
			log.Println(response)

			// Set values in prometheus
			numberOfDownlinkMessages.WithLabelValues(gatewayId).Set(float64(response.RoundTripTimes.Count))
			uplinkMessages, err := response.GetUplinkCount()

			if err != nil {
				log.Printf("WARNING: Failed to parse uplink count: %v", err)
				continue
			} else {
				numberOfUplinkMessages.WithLabelValues(gatewayId).Set(uplinkMessages)
			}

			min, max, median, err := response.RoundTripTimes.ConvertToSeconds()

			if err != nil {
				log.Printf("WARNING: Failed to parse rtt values: %v", err)
				continue
			} else {
				numberOfUplinkMessages.WithLabelValues(gatewayId).Set(uplinkMessages)
			}

			rtt_min.WithLabelValues(gatewayId).Set(min)
			rtt_median.WithLabelValues(gatewayId).Set(median)
			rtt_max.WithLabelValues(gatewayId).Set(float64(max))
		}

		duration := time.Since(start).Seconds()
		log.Printf("Done (Last request duration: %.5fs) \n", duration)
		lastApiCallDuration.Set(duration)
	}
}
