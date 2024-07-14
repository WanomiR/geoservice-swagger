package main

import (
	"log"
	_ "proxy/docs"
	"proxy/internal/app"
)

// @title Geoservice API
// @version 2.0.0
// @description Geoservice with swagger docs and authentication

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to init app: %v", err.Error())
	}

	log.Fatal(a.Run())
}
