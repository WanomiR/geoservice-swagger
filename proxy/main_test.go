package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

const uriBase = "http://localhost:8080"

func TestMainFunc(t *testing.T) {
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

	go main()

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
