package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetLocationAreas(url string) (LocationAreaResponse, error) {
	// piping and receiving the response
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "PokedexCLI/1.0 - shady")

	client := &http.Client{}
	resp, err := client.Do(req)

	// resp, err := http.Get(*config.Next)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("Failed to GET location areas: %w", err)
	}
	defer resp.Body.Close()

	//resp is RAW TEXT:  --- we need ot convert it to []Bytes and then to go struct --- unmarshal method
	// else we take resp.Body -> as a stream of text and send it to newDecoder

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("Error reading response body: %w", err)
	}

	var loc_area_response LocationAreaResponse

	err = json.Unmarshal(responseBytes, &loc_area_response)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("Error unmarshaling: %w", err)
	}
	return loc_area_response, nil
}
