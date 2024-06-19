package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleAddressSearch(t *testing.T) {
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

			handleAddressSearch(w, req)

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

func TestHandleAddressGeocode(t *testing.T) {
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

			handleAddressGeocode(w, req)

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
