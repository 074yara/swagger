package main

import (
	"context"
	"encoding/json"
	"github.com/ekomobile/dadata/v2"
	"log"
	"net/http"
)

const (
	AddressFromGeoAPIURL = "http://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"
	dadataAPIKey         = "2cb6ce15db3cb2b52b47f9c39d250875b89d0723"
	dadataSecretKey      = "617a4664fbe290c49ea12e5500a53e5e69995246"
)

type Address struct {
	Result string `json:"result,omitempty"`
	GeoLat string `json:"lat,omitempty"`
	GeoLon string `json:"lon,omitempty"`
}

//Маршрут: /api/address/search метод POST

type SearchRequest struct {
	Query string `json:"query"`
}
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}

//Маршрут: /api/address/geocode метод POST

type AddressFromGeoResponse struct {
	Addresses []*struct {
		Value string `json:"value"`
	} `json:"suggestions"`
}

// @Summary		Get coordinates from address
// @Description	Get coordinates from address
// @Tags			search
// @Accept			json
// @Produce		json
// @Param			input	body		SearchRequest	true	"SearchRequest"
// @Success		200		{object}	SearchResponse
// @Failure		503
// @Failure		404
// @Router			/address/search [post]
func geoFromAddressHandler(rw http.ResponseWriter, r *http.Request) {
	var searchRequest SearchRequest
	var searchResponse SearchResponse

	err := json.NewDecoder(r.Body).Decode(&searchRequest)
	checkError(err)

	api := dadata.NewCleanApi()
	addresses, err := api.Address(context.Background(), searchRequest.Query)
	log.Println(searchRequest.Query)
	checkError(err)

	if addresses == nil {
		http.Error(rw, "No answer from Dadata API", http.StatusServiceUnavailable)
		return
	}

	searchResponse.Addresses = []*Address{{GeoLat: addresses[0].GeoLat, GeoLon: addresses[0].GeoLon, Result: addresses[0].Result}}
	err = json.NewEncoder(rw).Encode(searchResponse)
	checkError(err)
}

// @Summary		Get address from coordinates
// @Description	Get address from coordinates
// @Tags			geocode
// @Accept			json
// @Produce		json
// @Success		200		{object}	SearchResponse
// @Param			input	body		Address	true	"Address"
// @Failure		503
// @Failure		403
// @Router			/address/geocode [post]
func addressFromGeoHandler(rw http.ResponseWriter, r *http.Request) {
	var addressFromGeo AddressFromGeoResponse
	var searchResponse SearchResponse
	newReq, err := http.NewRequest("POST", AddressFromGeoAPIURL, r.Body)
	checkError(err)

	newReq.Header.Set("Content-Type", "application/json")
	newReq.Header.Set("Accept", "application/json")
	newReq.Header.Set("Authorization", "Token "+dadataAPIKey)
	newReq.Header.Set("X-Secret", dadataSecretKey)

	response, err := http.DefaultClient.Do(newReq)
	checkError(err)
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&addressFromGeo)
	checkError(err)
	query := addressFromGeo.Addresses[0].Value
	log.Println(query)

	api := dadata.NewCleanApi()
	addresses, err := api.Address(context.Background(), query)
	checkError(err)

	if addresses == nil {
		http.Error(rw, "No answer from Dadata API", http.StatusServiceUnavailable)
		return
	}

	searchResponse.Addresses = []*Address{{Result: addresses[0].Result, GeoLon: addresses[0].GeoLon, GeoLat: addresses[0].GeoLat}}
	err = json.NewEncoder(rw).Encode(searchResponse)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}

}
