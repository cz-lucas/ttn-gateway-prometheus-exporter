package main

import (
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
	LastUplinkReceivedAt time.Time `json:"last_uplink_received_at"`
	UplinkCount          string    `json:"uplink_count"`
	RoundTripTimes       struct {
		Min    string `json:"min"`
		Max    string `json:"max"`
		Median string `json:"median"`
		Count  int    `json:"count"`
	} `json:"round_trip_times"`
	GatewayRemoteAddress struct {
		IP string `json:"ip"`
	} `json:"gateway_remote_address"`
}
