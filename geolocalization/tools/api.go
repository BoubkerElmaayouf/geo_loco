package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func FetchData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

func FetchArtistData(baseURL string) ([]Card, error) {
	var api APIindex
	if err := FetchData(baseURL, &api); err != nil {
		return nil, err
	}
	var artists []Artist
	if err := FetchData(api.Artists, &artists); err != nil {
		return nil, err
	}
	var location APIlocations
	if err := FetchData(api.Locations, &location); err != nil {
		return nil, err
	}
	var dates APIdates
	if err := FetchData(api.Dates, &dates); err != nil {
		return nil, err
	}
	var relation APIrelations
	if err := FetchData(api.Relations, &relation); err != nil {
		return nil, err
	}

	var cards []Card
	for i, artist := range artists {
		cards = append(cards, Card{
			Id:           artist.Id,
			Image:        artist.Image,
			Name:         artist.Name,
			Members:      artist.Members,
			CreationDate: artist.CreationDate,
			FirstAlbum:   artist.FirstAlbum,
			Locations:    location.Index[i].Location,
			ConcertDates: dates.Index[i].Dates,
			Relation:     relation.Index[i].DatesLocations,
		})
	}
	return cards, nil
}

func MakerMaps(ArrMaps [][]string) string {
	baseURL := "https://maps.googleapis.com/maps/api/staticmap"
	params := url.Values{}
	params.Set("key", "AIzaSyC_wqzY6A0mcDLKKXJds45WfDmVkqBFUSQ")
	params.Set("size", "600x600")
	var markers []string
	for _, coords := range ArrMaps {
		if len(coords) != 2 {
			continue // Skip if coordinates are not valid
		}
		lat := coords[0]
		lng := coords[1]
		markers = append(markers, fmt.Sprintf("%s,%s", lat, lng))
	}

	// Use the first coordinate set for the center of the map
	if len(ArrMaps) > 0 {
		center := fmt.Sprintf("%s,%s", ArrMaps[0][0], ArrMaps[0][1])
		params.Set("center", center)
	}

	// Add all markers to the map
	if len(markers) > 0 {
		params.Set("markers", strings.Join(markers, "|"))
	}

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func GetStaticMapURL(address []string) [][]string {
	var data [][]string
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Set("key", "AIzaSyC_wqzY6A0mcDLKKXJds45WfDmVkqBFUSQ")

	for _, add := range address {
		params.Set("address", add)
		resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
		if err != nil {
			log.Printf("Error fetching geocoding data: %v", err)
			continue // Skip to the next address if there's an error
		}
		defer resp.Body.Close()

		var result GeocodingResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("Error decoding geocoding response: %v", err)
			continue
		}

		if len(result.Results) > 0 {
			lat := strconv.FormatFloat(result.Results[0].Geometry.Location.Lat, 'f', 6, 64)
			lng := strconv.FormatFloat(result.Results[0].Geometry.Location.Lng, 'f', 6, 64)
			data = append(data, []string{lat, lng})
		}
	}
	return data
}
