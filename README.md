# TTN Gateway Prometheus Exporter
A Go application that fetches gateway statistics from The Things Network (TTN) API and exposes them as Prometheus metrics for monitoring and alerting.

## Overview
This exporter periodically queries the TTN API for gateway connection statistics and converts them into Prometheus metrics. It's designed to help monitor LoRaWAN gateway health and performance through standard observability tools.

## Features
- Fetches gateway statistics from TTN API
- Exposes metrics via /metrics endpoint
- Configurable polling intervals
- Health check endpoint
- Runtime and application metrics
- Docker support with health checks

## Configuration
The application is configured via environment variables

### Required Environment Variables
| Variable               | Description                                                               | Optional | Default value                                           |
|------------------------|---------------------------------------------------------------------------|----------|---------------------------------------------------------|
| TTN_GATEWAY_ID         | The ID of the gateway                                                     | ❌        | -                                                       |
| TTN_API_KEY            | The TTN API-Key with read permissions                                     | ❌        | -                                                       |
| READ_INTERVAL          | The interval in seconds how often the data should be fetched from the TTN | ✅        | 600s                                                    |
| ADDRESS                | The bind address                                                          | ✅        | :9000                                                   |
| TTN_BASE_URL           | The TTN base url (need of you want to use another region                  | ✅        | https://eu1.cloud.thethings.network/api/v3/gs/gateways/ |
| TTN_URL_SUFFIX         | The suffix in the url (normally there is no need to change it)            | ✅        | /connection/stats                                       |
| ENABLE_RUNTIME_METRICS | Enable the go runtime metrics                                             | ✅        | true                                                    |
| ENABLE_APP_METRICS     | Enable the metrics from this tool                                         | ✅        | true                                                    |
## Metrics
### Gateway Metrics
| Metric                         | Type  | Description                        |
|--------------------------------|-------|------------------------------------|
| gw_number_of_uplink_messages   | Gauge | Total number of uplink messages    |
| gw_number_of_downlink_messages | Gauge | Total number of downlink messages  |
| gw_rtt_min                     | Gauge | Minimum round trip time in seconds |
| gw_rtt_median                  | Gauge | Median round trip time in seconds  |
| gw_rtt_max                     | Gauge | Maximum round trip time in seconds |

### Application Metrics
| Metric                         | Type    | Description                      |
|--------------------------------|---------|----------------------------------|
| api_calls_total                | Counter | Total number of API calls made   |
| api_call_failures_total        | Counter | Total number of failed API calls |
| last_api_call_duration_seconds | Gauge   | Duration of the last API call    |

## Installation
### Using Docker
``` bash
docker run --rm --env TTN_GATEWAY_ID=<gw-id> --env TTN_API_KEY="<api-token>" --env READ_INTERVAL=600 -p 9000:9000 czlucas/ttn-gateway-prometheus-exporter:latest
```

### Build and run using Go
1. Clone the repository:
``` bash
git clone <repository-url>
cd ttn-gateway-prometheus-exporter
```
2. Install dependencies:
``` bash
go mod download
```
3.Set environment variables and run:
``` bash
TTN_GATEWAY_ID=your-gateway-id TTN_API_KEY=your-api-key go run .
```
or create a .env-file

## API Endpoints
| Endpoint | Description                 |
|----------|-----------------------------|
| /metrics | Prometheus metrics endpoint |
| /health  | Health check endpoint       |

## Code Structure
### Core Components
- main.go - Application entry point and main loop
- TTNApiService.go - TTN API client implementation
- GatewayStats.go - Data structures and conversion methods
- PrometheusMetrics.go - Prometheus metrics definitions
- HttpService.go - HTTP server implementation
- utils.go - Utility functions for environment variable handling


### Key Functions
#### Gateway Data Processing
- `GatewayStats.GetUplinkCount()` - Converts uplink count to float64
- `RoundTripTimes.ConvertToSeconds()` - Converts RTT duration strings to seconds
- `convertDurationToSeconds()` - Helper for duration conversion
- `stringsToFloat64()` - Helper for string to float conversion
#### Utility Functions
- `getEnvBool()` - Get boolean environment variable with default
- `getEnvInt()` - Get integer environment variable with default
- `getEnvString()` - Get string environment variable with default
- `keyExistsInConfig()` - Check if environment variable exists
#### Services
- `NewTTNApiService()` - Create TTN API client
- `NewHttpService()` - Create HTTP server
- `InitPrometheus()` - Initialize Prometheus registry

## Testing
The project includes tests.
You can run them with:
``` bash
go test
```

## Development
### Dev Container Support
The project includes VS Code dev container configuration in .devcontainer:

### Building

## Docker Health Checks
The Docker image includes health checks that verify the /health endpoint:

## Error Handling
The application includes robust error handling:

- API connection failures are logged and retried on next interval
- Invalid data parsing is logged as warnings
- Failed metric updates don't crash the application
- HTTP server errors are logged appropriately

# License
This project is licensed under the Apache License 2.0. See LICENSE for details.

# Contributing
This is the author's first Go project, so contributions and suggestions are welcome! Please feel free to submit issues and pull requests.