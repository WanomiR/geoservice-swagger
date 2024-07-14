package geoservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"net/http"
	"net/url"
	"strings"
)

type GeoService struct {
	api       *suggest.Api
	apiKey    string
	secretKey string
}

type RequestAddressSearch struct {
	Query string `json:"query" example:"Подкопаевский переулок"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat" example:"55.753214"`
	Lng string `json:"lng" example:"37.642589"`
}

func NewGeoService(apiKey, secretKey string) *GeoService {
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

// HandleAddressSearch
// @Summary Search by street name
// @Security ApiKeyAuth
// @Description Return a list of addresses provided street name
// @Tags address
// @Accept json
// @Produce json
// @Param query body RequestAddressSearch true "street name"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal error"
// @Router /api/address/search [post]
func (g *GeoService) HandleAddressSearch(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressSearch
	json.NewDecoder(r.Body).Decode(&req)

	if req.Query == "" {
		http.Error(w, "bad query", http.StatusBadRequest)
		return
	}

	addresses, err := g.AddressSearch(req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseAddress{Addresses: addresses})
}

// HandleAddressGeocode
// @Summary Search by coordinates
// @Security ApiKeyAuth
// @Description Return a list of addresses provided geo coordinates
// @Tags address
// @Accept json
// @Produce json
// @Param query body RequestAddressGeocode true "coordinates"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal error"
// @Router /api/address/geocode [post]
func (g *GeoService) HandleAddressGeocode(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressGeocode
	json.NewDecoder(r.Body).Decode(&req)

	if req.Lng == "" || req.Lat == "" {
		http.Error(w, "bad query", http.StatusBadRequest)
		return
	}

	addresses, err := g.GeoCode(req.Lat, req.Lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseAddress{Addresses: addresses})
}
