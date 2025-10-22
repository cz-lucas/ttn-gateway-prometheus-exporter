package main

import (
	"fmt"
	"strconv"
	"time"
)

type GatewayStats struct {
	ConnectedAt          time.Time `json:"connected_at"`
	Protocol             string    `json:"protocol"`
	LastStatusReceivedAt time.Time `json:"last_status_received_at"`
	LastStatus           struct {
		Versions struct {
			Package  string `json:"package"`
			Platform string `json:"platform"`
			Station  string `json:"station"`
			Firmware string `json:"firmware"`
		} `json:"versions"`
		Advanced struct {
			Features string `json:"features"`
			Model    string `json:"model"`
		} `json:"advanced"`
	} `json:"last_status"`
	LastUplinkReceivedAt time.Time      `json:"last_uplink_received_at"`
	UplinkCount          string         `json:"uplink_count"`
	RoundTripTimes       RoundTripTimes `json:"round_trip_times"`
	GatewayRemoteAddress struct {
		IP string `json:"ip"`
	} `json:"gateway_remote_address"`
}

type RoundTripTimes struct {
	Min    string `json:"min"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Count  int    `json:"count"`
}

// Method to convert RTT values (min, max, median) to seconds
func (rtt *RoundTripTimes) ConvertToSeconds() (float64, float64, float64, error) {
	// Convert each time string to seconds and store in the struct
	var err error
	min, err := convertDurationToSeconds(rtt.Min)
	if err != nil {
		return -1.0, -1.0, -1.0, fmt.Errorf("error parsing min duration: %v", err)
	}
	max, err := convertDurationToSeconds(rtt.Max)
	if err != nil {
		return -1.0, -1.0, -1.0, fmt.Errorf("error parsing max duration: %v", err)
	}
	median, err := convertDurationToSeconds(rtt.Median)
	if err != nil {
		return -1.0, -1.0, -1.0, fmt.Errorf("error parsing median duration: %v", err)
	}
	return min, median, max, nil
}

// Helper function to convert time.Duration strings to seconds as float64
func convertDurationToSeconds(durationStr string) (float64, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return -1, err
	}
	// Convert to seconds and return as float64
	return duration.Seconds(), nil
}

// Helper function to convert string to float64
func stringsToFloat64(input string) (float64, error) {
	// Convert the string to float64 using strconv.ParseFloat
	result, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return -1, err // If there's an error, return -1 and the error
	}
	return result, nil // Return the parsed float64 value
}

// Method to convert RTT values (min, max, median) to seconds
func (GatewayStats *GatewayStats) GetUplinkCount() (float64, error) {
	// Convert each time string to seconds and store in the struct
	var err error
	uplinkCount, err := stringsToFloat64(GatewayStats.UplinkCount)
	if err != nil {
		return -1, fmt.Errorf("error parsing uplinkCount: %v", err)
	}

	return uplinkCount, nil
}
