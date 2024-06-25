package geoservice

import (
	"context"
	"encoding/json"
	"fmt"
	"geoservice-swagger/proxy/reverse"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/url"
	"strings"
)

type GeoProvider interface {
	AddressSearch(input string) ([]*Address, error)
	GeoCode(lat, lng string) ([]*Address, error)
}

type GeoService struct {
	addr      string
	api       *suggest.Api
	apiKey    string
	secretKey string
}

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

func NewGeoService(apiKey, secretKey, addr string) *GeoService {
	var err error
	endpointUrl, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	if err != nil {
		return nil
	}

	creds := client.Credentials{
		ApiKeyValue:    apiKey,
		SecretKeyValue: secretKey,
	}

	api := suggest.Api{
		Client: client.NewClient(endpointUrl, client.WithCredentialProvider(&creds)),
	}

	return &GeoService{
		addr:      addr,
		api:       &api,
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
	House  string `json:"house"`
	Lat    string `json:"lat"`
	Lon    string `json:"lon"`
}

func (g *GeoService) AddressSearch(input string) ([]*Address, error) {
	var res []*Address
	rawRes, err := g.api.Address(context.Background(), &suggest.RequestParams{Query: input})
	if err != nil {
		return nil, err
	}

	for _, r := range rawRes {
		if r.Data.City == "" || r.Data.Street == "" {
			continue
		}
		res = append(res, &Address{City: r.Data.City, Street: r.Data.Street, House: r.Data.House, Lat: r.Data.GeoLat, Lon: r.Data.GeoLon})
	}

	return res, nil
}

func (g *GeoService) GeoCode(lat, lng string) ([]*Address, error) {
	httpClient := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"lat": %s, "lon": %s}`, lat, lng))
	req, err := http.NewRequest("POST", "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", g.apiKey))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var geoCode GeoCode

	err = json.NewDecoder(resp.Body).Decode(&geoCode)
	if err != nil {
		return nil, err
	}
	var res []*Address
	for _, r := range geoCode.Suggestions {
		var address Address
		address.City = string(r.Data.City)
		address.Street = string(r.Data.Street)
		address.House = r.Data.House
		address.Lat = r.Data.GeoLat
		address.Lon = r.Data.GeoLon

		res = append(res, &address)
	}

	return res, nil
}

func (g *GeoService) HandleAddressSearch(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressSearch
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Query == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addresses, err := g.AddressSearch(req.Query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseAddress{Addresses: addresses})
}

func (g *GeoService) HandleAddressGeocode(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressGeocode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Lat == "" || req.Lng == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addresses, err := g.GeoCode(req.Lat, req.Lng)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseAddress{Addresses: addresses})
}

func (g *GeoService) ListenAndServe() error {

	r := chi.NewRouter()

	proxy := reverse.NewReverseProxy("hugo", "1313")
	r.Use(proxy.ReverseProxy)

	r.Route("/api/address", func(r chi.Router) {
		r.Post("/search", g.HandleAddressSearch)
		r.Post("/geocode", g.HandleAddressGeocode)
	})

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello from API"))
	})

	return http.ListenAndServe(g.addr, r)
}
