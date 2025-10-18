package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	var ttnRequestUrl = os.Getenv("TTN_BASE_URL") + os.Getenv("TTN_GATEWAY_ID") + os.Getenv("TTN_URL_STATS_SUFFIX")
	var apiService = NewTTNApiService(ttnRequestUrl, os.Getenv("TTN_API_KEY"))

	intervalInSeconds, err := strconv.Atoi(os.Getenv("READ_INTERVAL"))

	if err != nil {
		log.Fatal("Invalid READ_INTERVAL:", err)
	}

	ticker := time.NewTicker(time.Duration(intervalInSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Doing work at", time.Now())
			log.Println(apiService.Get())
		}
	}
}
