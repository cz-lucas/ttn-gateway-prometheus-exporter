package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	var ttnRequestUrl = os.Getenv("TTN_BASE_URL") + os.Getenv("TTN_GATEWAY_ID") + os.Getenv("TTN_URL_STATS_SUFFIX")
	log.Println(ttnRequestUrl)

	var apiService = NewTTNApiService(ttnRequestUrl, os.Getenv("TTN_API_KEY"))
	log.Println(apiService.Get())
}
