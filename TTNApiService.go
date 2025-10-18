package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TTNApiService struct {
	url      string
	apiToken string
	client   *http.Client
}

func NewTTNApiService(url string, apiToken string) *TTNApiService {
	return &TTNApiService{
		url:      url,
		apiToken: apiToken,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (ttn *TTNApiService) Get() (GatewayStats, error) {
	req, err := http.NewRequest("GET", ttn.url, nil)
	if err != nil {
		return GatewayStats{}, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+ttn.apiToken)

	resp, err := ttn.client.Do(req)
	if err != nil {
		return GatewayStats{}, fmt.Errorf("making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GatewayStats{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GatewayStats{}, fmt.Errorf("reading response body: %w", err)
	}

	var stats GatewayStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return GatewayStats{}, fmt.Errorf("unmarshalling response: %w", err)
	}

	return stats, nil
}
