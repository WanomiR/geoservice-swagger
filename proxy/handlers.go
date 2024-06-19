package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

var geoService *GeoService

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	apiKey := os.Getenv("DADATA_API_KEY")
	secretKey := os.Getenv("DADATA_SECRET_KEY")

	geoService = NewGeoService(apiKey, secretKey)
}

func handleAddressSearch(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressSearch
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Query == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addresses, err := geoService.AddressSearch(req.Query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseAddress{Addresses: addresses})
}

func handleAddressGeocode(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressGeocode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Lat == "" || req.Lng == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addresses, err := geoService.GeoCode(req.Lat, req.Lng)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseAddress{Addresses: addresses})
}
