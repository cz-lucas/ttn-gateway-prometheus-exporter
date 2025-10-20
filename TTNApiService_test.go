package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTTNApiService(t *testing.T) {
	url := "https://example.com/api"
	apiToken := "test-token"

	service := NewTTNApiService(url, apiToken)

	assert.NotNil(t, service)
	assert.Equal(t, url, service.url)
	assert.Equal(t, apiToken, service.apiToken)

	// Check if client is initialized
	assert.NotNil(t, service.client)
	assert.IsType(t, &http.Client{}, service.client)

	// Check client timeout
	expectedTimeout := 10 * time.Second
	assert.Equal(t, service.client.Timeout, expectedTimeout, "Timeout is not as expected")
}
