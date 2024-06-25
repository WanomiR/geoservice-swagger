package geoservice

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var geoservice *GeoService

const uriBase = "http://localhost:8080"

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	apiKey := os.Getenv("DADATA_API_KEY")
	secretKey := os.Getenv("DADATA_SECRET_KEY")
	addr := os.Getenv("ADDR")

	geoservice = NewGeoService(apiKey, secretKey, addr)
}

func TestGeoService_ListenAndServe(t *testing.T) {
	type reqParams struct {
		method string
		url    string
		body   io.Reader
	}
	testCases := []struct {
		name           string
		params         reqParams
		wantStatusCode int
	}{
		{
			"api",
			reqParams{"GET", uriBase + "/api", nil},
			200,
		},
		{
			"address search",
			reqParams{"POST", uriBase + "/api/address/search", strings.NewReader(`{"query":"Ленина"}`)},
			200 | 500,
		},
		{
			"address geocode",
			reqParams{"POST", uriBase + "/api/address/geocode", strings.NewReader(`"lat":"55.753203", "lng":"37.642560"`)},
			200 | 500,
		},
		{
			"redirect to hugo",
			reqParams{"GET", uriBase + "/", nil},
			200,
		},
	}

	go func() {
		log.Fatal(geoservice.ListenAndServe())
	}()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.params.method, tc.params.url, tc.params.body)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("request failed: %v", err)
			}

			if resp.StatusCode&tc.wantStatusCode == 0 {
				t.Errorf("status code should be %d, got %d", tc.wantStatusCode, resp.StatusCode)
			}
		})
	}
}

func TestGeoService_HandleAddressSearch(t *testing.T) {
	testCases := []struct {
		name          string
		query         string
		wantStatus    int
		wantAddresses bool
	}{
		{"good request", "Ленина", 200 | 500, true},
		{"wrong address", "aljgfag", 200 | 500, false},
		{"bad request", "", 400 | 500, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.NewReader(fmt.Sprintf(`{"query":"%s"}`, tc.query))
			req := httptest.NewRequest("POST", "/api/address/search", body)
			w := httptest.NewRecorder()

			geoservice.HandleAddressSearch(w, req)

			r := w.Result()
			defer r.Body.Close()

			var responseAddresses ResponseAddress
			json.NewDecoder(r.Body).Decode(&responseAddresses)

			if r.StatusCode&tc.wantStatus == 0 {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}

			if len(responseAddresses.Addresses) != 0 != tc.wantAddresses {
				t.Errorf("got addresses %v, want %v", responseAddresses.Addresses, tc.wantAddresses)
			}
		})
	}

}

func TestGeoService_handleAddressGeocode(t *testing.T) {
	tests := []struct {
		name          string
		query         RequestAddressGeocode
		wantStatus    int
		wantAddresses bool
	}{
		{"good request", RequestAddressGeocode{Lat: "55.753203", Lng: "37.642560"}, 200 | 500, true},
		{"wrong coordinates", RequestAddressGeocode{Lat: "99.753203", Lng: "0.642560"}, 200 | 500, false},
		{"bad request", RequestAddressGeocode{Lat: "", Lng: ""}, 400 | 500, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.NewReader(fmt.Sprintf(`{"lat":"%s", "lng":"%s"}`, tc.query.Lat, tc.query.Lng))
			req := httptest.NewRequest("POST", "/api/address/geocode", body)
			w := httptest.NewRecorder()

			geoservice.HandleAddressGeocode(w, req)

			r := w.Result()
			defer r.Body.Close()

			var responseAddresses ResponseAddress
			json.NewDecoder(r.Body).Decode(&responseAddresses)

			if r.StatusCode&tc.wantStatus == 0 {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}

			if len(responseAddresses.Addresses) != 0 != tc.wantAddresses {
				t.Errorf("got addresses %v, want %v", responseAddresses.Addresses, tc.wantAddresses)
			}
		})
	}

}
