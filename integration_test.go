package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TestHTTPServiceIntegration tests the HTTP service with real handlers
func TestHTTPServiceIntegration(t *testing.T) {
	// Setup
	httpService := NewHttpService(":0") // :0 picks a random available port

	// Register a test route
	httpService.RegisterRoute("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))

	// Start the service
	err := httpService.Start()
	if err != nil {
		t.Fatalf("Failed to start HTTP service: %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// We can't easily test the actual HTTP call without knowing the random port
	// But we verified it starts without panicking
	t.Log("HTTP service started successfully")
}

// TestHealthEndpoint tests the /health endpoint
func TestHealthEndpoint(t *testing.T) {
	// Create the HTTP service
	httpService := NewHttpService()

	// Register health endpoint
	httpService.RegisterRoute("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	// Create a test request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Serve the request
	httpService.mux.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "ok" {
		t.Errorf("Expected 'ok', got '%s'", w.Body.String())
	}
}

// TestPrometheusMetricsEndpoint tests that metrics are exposed correctly
func TestPrometheusMetricsEndpoint(t *testing.T) {
	// Initialize Prometheus with app metrics enabled
	reg := InitPrometheus(false, true) // runtime=false, app=true for simpler output

	// Increment some metrics
	apiCallsTotal.Inc()
	apiCallsTotal.Inc()
	apiCallFailures.Inc()
	lastApiCallDuration.Set(1.234)

	// Create HTTP service and register metrics
	httpService := NewHttpService()
	httpService.RegisterRoute("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	// Create test request
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	// Serve the request
	httpService.mux.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()

	// Verify metrics are present
	expectedMetrics := []string{
		"api_calls_total",
		"api_call_failures_total",
		"last_api_call_duration_seconds",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(body, metric) {
			t.Errorf("Expected metric '%s' not found in response", metric)
		}
	}

	// Verify values
	if !strings.Contains(body, "api_calls_total 2") {
		t.Errorf("Expected 'api_calls_total 2' not found")
	}
	if !strings.Contains(body, "api_call_failures_total 1") {
		t.Errorf("Expected 'api_call_failures_total 1' not found")
	}
}

// TestTTNApiServiceWithMockServer tests the API service with a mock TTN server
func TestTTNApiServiceWithMockServer(t *testing.T) {
	// Create a mock TTN server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-api-key" {
			t.Errorf("Expected 'Bearer test-api-key', got '%s'", authHeader)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Return mock gateway stats
		mockStats := GatewayStats{
			ConnectedAt:          time.Now(),
			Protocol:             "udp",
			LastStatusReceivedAt: time.Now(),
			LastUplinkReceivedAt: time.Now(),
			UplinkCount:          "42",
			RoundTripTimes: RoundTripTimes{
				Min:    "10ms",
				Max:    "100ms",
				Median: "50ms",
				Count:  10,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockStats)
	}))
	defer mockServer.Close()

	// Create API service pointing to mock server
	apiService := NewTTNApiService(mockServer.URL, "test-api-key")

	// Make request
	stats, err := apiService.Get()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify response
	if stats.UplinkCount != "42" {
		t.Errorf("Expected UplinkCount '42', got '%s'", stats.UplinkCount)
	}

	if stats.RoundTripTimes.Count != 10 {
		t.Errorf("Expected RTT Count 10, got %d", stats.RoundTripTimes.Count)
	}

	if stats.Protocol != "udp" {
		t.Errorf("Expected protocol 'udp', got '%s'", stats.Protocol)
	}
}

// TestTTNApiServiceErrorHandling tests various error scenarios
func TestTTNApiServiceErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectError    bool
		errorContains  string
	}{
		{
			name: "successful request",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				mockStats := GatewayStats{
					UplinkCount: "100",
					RoundTripTimes: RoundTripTimes{
						Min:    "5ms",
						Max:    "50ms",
						Median: "25ms",
						Count:  5,
					},
				}
				json.NewEncoder(w).Encode(mockStats)
			},
			expectError: false,
		},
		{
			name: "404 not found",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectError:   true,
			errorContains: "unexpected status code: 404",
		},
		{
			name: "500 server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectError:   true,
			errorContains: "unexpected status code: 500",
		},
		{
			name: "invalid json",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("this is not json"))
			},
			expectError:   true,
			errorContains: "unmarshalling response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer mockServer.Close()

			apiService := NewTTNApiService(mockServer.URL, "test-key")
			_, err := apiService.Get()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if tt.expectError && err != nil && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
			}
		})
	}
}

// TestGatewayStatsConversions tests the data conversion methods
func TestGatewayStatsConversions(t *testing.T) {
	t.Run("RTT conversion success", func(t *testing.T) {
		rtt := RoundTripTimes{
			Min:    "10ms",
			Max:    "100ms",
			Median: "50ms",
			Count:  5,
		}

		min, median, max, err := rtt.ConvertToSeconds()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if min != 0.01 {
			t.Errorf("Expected min 0.01, got %f", min)
		}
		if median != 0.05 {
			t.Errorf("Expected median 0.05, got %f", median)
		}
		if max != 0.1 {
			t.Errorf("Expected max 0.1, got %f", max)
		}
	})

	t.Run("RTT conversion with invalid data", func(t *testing.T) {
		rtt := RoundTripTimes{
			Min:    "invalid",
			Max:    "100ms",
			Median: "50ms",
			Count:  5,
		}

		_, _, _, err := rtt.ConvertToSeconds()
		if err == nil {
			t.Error("Expected error for invalid duration, got none")
		}
	})

	t.Run("UplinkCount conversion success", func(t *testing.T) {
		stats := GatewayStats{
			UplinkCount: "12345",
		}

		count, err := stats.GetUplinkCount()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if count != 12345.0 {
			t.Errorf("Expected 12345.0, got %f", count)
		}
	})

	t.Run("UplinkCount conversion with invalid data", func(t *testing.T) {
		stats := GatewayStats{
			UplinkCount: "not-a-number",
		}

		_, err := stats.GetUplinkCount()
		if err == nil {
			t.Error("Expected error for invalid number, got none")
		}
	})
}

// TestUtilityFunctions tests the utility helper functions
func TestUtilityFunctions(t *testing.T) {
	t.Run("getEnvBool with valid values", func(t *testing.T) {
		os.Setenv("TEST_BOOL_TRUE", "true")
		os.Setenv("TEST_BOOL_FALSE", "false")
		defer func() {
			os.Unsetenv("TEST_BOOL_TRUE")
			os.Unsetenv("TEST_BOOL_FALSE")
		}()

		val, err := getEnvBool("TEST_BOOL_TRUE", false)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if !val {
			t.Error("Expected true, got false")
		}

		val, err = getEnvBool("TEST_BOOL_FALSE", true)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if val {
			t.Error("Expected false, got true")
		}
	})

	t.Run("getEnvBool with default", func(t *testing.T) {
		val, err := getEnvBool("NONEXISTENT_BOOL", true)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if !val {
			t.Error("Expected default value true, got false")
		}
	})

	t.Run("getEnvInt with valid value", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		val, err := getEnvInt("TEST_INT", 0)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if val != 42 {
			t.Errorf("Expected 42, got %d", val)
		}
	})

	t.Run("getEnvInt with invalid value", func(t *testing.T) {
		os.Setenv("TEST_INT_INVALID", "not-a-number")
		defer os.Unsetenv("TEST_INT_INVALID")

		_, err := getEnvInt("TEST_INT_INVALID", 100)
		if err == nil {
			t.Error("Expected error for invalid int, got none")
		}
	})

	t.Run("getEnvString", func(t *testing.T) {
		os.Setenv("TEST_STRING", "hello")
		defer os.Unsetenv("TEST_STRING")

		val := getEnvString("TEST_STRING", "default")
		if val != "hello" {
			t.Errorf("Expected 'hello', got '%s'", val)
		}

		val = getEnvString("NONEXISTENT_STRING", "default")
		if val != "default" {
			t.Errorf("Expected 'default', got '%s'", val)
		}
	})

	t.Run("keyExistsInConfig", func(t *testing.T) {
		os.Setenv("TEST_KEY", "value")
		defer os.Unsetenv("TEST_KEY")

		if !keyExistsInConfig("TEST_KEY") {
			t.Error("Expected key to exist")
		}

		if keyExistsInConfig("NONEXISTENT_KEY") {
			t.Error("Expected key to not exist")
		}
	})
}

// TestEndToEndFlow simulates a complete workflow
func TestEndToEndFlow(t *testing.T) {
	// Create mock TTN server
	requestCount := 0
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		mockStats := GatewayStats{
			ConnectedAt:          time.Now(),
			Protocol:             "udp",
			LastStatusReceivedAt: time.Now(),
			LastUplinkReceivedAt: time.Now(),
			UplinkCount:          fmt.Sprintf("%d", requestCount*100),
			RoundTripTimes: RoundTripTimes{
				Min:    "10ms",
				Max:    "100ms",
				Median: "50ms",
				Count:  requestCount * 10,
			},
		}

		json.NewEncoder(w).Encode(mockStats)
	}))
	defer mockServer.Close()

	// Setup services
	apiService := NewTTNApiService(mockServer.URL, "test-key")
	httpService := NewHttpService(":0")

	// Initialize Prometheus
	reg := InitPrometheus(false, true)
	httpService.RegisterRoute("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	httpService.Start()

	// Simulate fetching metrics multiple times
	gatewayId := "test-gateway"

	for i := 0; i < 3; i++ {
		stats, err := apiService.Get()
		if err != nil {
			t.Fatalf("API call %d failed: %v", i+1, err)
		}

		// Update Prometheus metrics
		apiCallsTotal.Inc()
		numberOfDownlinkMessages.WithLabelValues(gatewayId).Set(float64(stats.RoundTripTimes.Count))

		uplinkCount, err := stats.GetUplinkCount()
		if err != nil {
			t.Errorf("Failed to get uplink count on iteration %d: %v", i+1, err)
			continue
		}
		numberOfUplinkMessages.WithLabelValues(gatewayId).Set(uplinkCount)

		min, median, max, err := stats.RoundTripTimes.ConvertToSeconds()
		if err != nil {
			t.Errorf("Failed to convert RTT on iteration %d: %v", i+1, err)
			continue
		}
		rtt_min.WithLabelValues(gatewayId).Set(min)
		rtt_median.WithLabelValues(gatewayId).Set(median)
		rtt_max.WithLabelValues(gatewayId).Set(max)
	}

	// Verify metrics endpoint returns data
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	httpService.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()

	// Verify gateway metrics are present
	if !strings.Contains(body, "gw_number_of_uplink_messages") {
		t.Error("Gateway uplink messages metric not found")
	}

	if !strings.Contains(body, "gw_number_of_downlink_messages") {
		t.Error("Gateway downlink messages metric not found")
	}

	if !strings.Contains(body, "gw_rtt_min") {
		t.Error("Gateway RTT min metric not found")
	}

	t.Logf("End-to-end test passed! Made %d requests to mock TTN server", requestCount)
}

// TestGracefulShutdown tests that the app can shutdown cleanly (Currently not implemented)
/*func TestGracefulShutdown(t *testing.T) {
	httpService := NewHttpService(":0")
	httpService.RegisterRoute("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	httpService.Start()
	time.Sleep(100 * time.Millisecond)

	// Shutdown with context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := httpService.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	t.Log("Graceful shutdown completed successfully")
}
*/
