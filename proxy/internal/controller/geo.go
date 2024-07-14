package controller

import (
	"errors"
	"net/http"
	"proxy/utils/readresponder"
)

type AddressSearchRequestBody struct {
	Query string `json:"query" binding:"required" example:"Подкопаевский переулок"`
}

type AddressGeocodeRequestBody struct {
	Lat string `json:"lat" example:"55.753214" binding:"required"`
	Lng string `json:"lng" example:"37.642589" binding:"required"`
}

// AddressSearch
// @Summary Search by street name
// @Security ApiKeyAuth
// @Description Return a list of addresses provided street name
// @Tags address
// @Accept json
// @Produce json
// @Param query body AddressSearchRequestBody true "street name"
// @Success 200 {object} readresponder.JSONResponse
// @Failure 400,500 {object} readresponder.JSONResponse
// @Router /api/address/search [post]
func (c *AppController) AddressSearch(w http.ResponseWriter, r *http.Request) {
	var req AddressSearchRequestBody

	if err := c.readResponder.ReadJSON(w, r, &req); err != nil {
		c.readResponder.WriteJSONError(w, err)
		return
	} else if req.Query == "" {
		c.readResponder.WriteJSONError(w, errors.New("query is required"))
		return
	}

	addresses, err := c.geoService.AddressSearch(req.Query)
	if err != nil {
		c.readResponder.WriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	resp := readresponder.JSONResponse{
		Error:   false,
		Message: "search completed",
		Data:    addresses,
	}

	c.readResponder.WriteJSON(w, http.StatusOK, resp)
}

// AddressGeocode
// @Summary Search by coordinates
// @Security ApiKeyAuth
// @Description Return a list of addresses provided geo coordinates
// @Tags address
// @Accept json
// @Produce json
// @Param query body AddressGeocodeRequestBody true "coordinates"
// @Success 200 {object} readresponder.JSONResponse
// @Failure 400,500 {object} readresponder.JSONResponse
// @Router /api/address/geocode [post]
func (c *AppController) AddressGeocode(w http.ResponseWriter, r *http.Request) {
	var req AddressGeocodeRequestBody
	if err := c.readResponder.ReadJSON(w, r, &req); err != nil {
		c.readResponder.WriteJSONError(w, err)
		return
	} else if req.Lat == "" || req.Lng == "" {
		c.readResponder.WriteJSONError(w, errors.New("both lat and lng are required"))
		return
	}

	addresses, err := c.geoService.GeoCode(req.Lat, req.Lng)
	if err != nil {
		c.readResponder.WriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	resp := readresponder.JSONResponse{
		Error:   false,
		Message: "search completed",
		Data:    addresses,
	}

	c.readResponder.WriteJSON(w, http.StatusOK, resp)
}
