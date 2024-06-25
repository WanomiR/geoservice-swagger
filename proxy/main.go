package main

import (
	"geoservice-swagger/proxy/geoservice"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	apiKey := os.Getenv("DADATA_API_KEY")
	secretKey := os.Getenv("DADATA_SECRET_KEY")
	addr := os.Getenv("ADDR")

	geoService := geoservice.NewGeoService(apiKey, secretKey, addr)

	log.Fatal(geoService.ListenAndServe())
}
