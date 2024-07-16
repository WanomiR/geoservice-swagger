package main

import (
	"context"
	"fmt"
	"log"
	_ "proxy/docs"
	"proxy/internal/app"
	"time"
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
		log.Fatalf("Failed to init app: %v", err.Error())
	}

	go a.Serve()

	// waiting for a stop signal
	<-a.Signal()

	// create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// countdown to make graceful shutdown explicit
	ticker, secs := time.NewTicker(1*time.Second), 5
	for {
		select {
		// countdown
		case <-ticker.C:
			fmt.Printf("%d...\n", secs)
			secs -= 1
		// stop server gracefully
		case <-ctx.Done():
			fmt.Println("Server stopped gracefully")
			return
		}
	}
}
