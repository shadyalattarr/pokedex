package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shadyalattarr/pokedex/internal/pokecache"
)

func NewClient(cacheInterval time.Duration) PokeApiClient {
	return PokeApiClient{
		cache: pokecache.NewCache(cacheInterval),
	}
}

func (pokeClient PokeApiClient) GetLocationAreas(url string) (LocationAreaResponse, error) {
	responseBytes, err := pokeClient.getResponseBytes(url)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("Error with getting response bytes: %v", err)
	}
	var loc_area_response LocationAreaResponse

	err = json.Unmarshal(responseBytes, &loc_area_response)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("Error unmarshaling: %w", err)
	}
	return loc_area_response, nil
}

func (pokeClient PokeApiClient) GetLocationAreaInfo(url string) (LocationAreaInformationResponse, error) {
	responseBytes, err := pokeClient.getResponseBytes(url)
	if err != nil {
		return LocationAreaInformationResponse{}, fmt.Errorf("Error with getting response bytes: %v", err)
	}
	var loc_area_info_response LocationAreaInformationResponse

	err = json.Unmarshal(responseBytes, &loc_area_info_response)
	if err != nil {
		return LocationAreaInformationResponse{}, fmt.Errorf("Error unmarshaling: %w", err)
	}
	return loc_area_info_response, nil

}

func (pokeClient PokeApiClient) getResponseBytes(url string) ([]byte, error) {
	var responseBytes []byte
	responseBytes, ok := pokeClient.cache.Get(url)
	if !ok { // if no cache, we overwrite responseBytes with the request!
		// piping and receiving the response
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []byte{}, fmt.Errorf("Failed to create request: %w", err)
		}
		req.Header.Set("User-Agent", "PokedexCLI/1.0 - shady")

		client := &http.Client{}
		resp, err := client.Do(req)

		// resp, err := http.Get(*config.Next)
		if err != nil {
			return []byte{}, fmt.Errorf("Failed to GET location area info: %w", err)
		}
		defer resp.Body.Close()

		//resp is RAW TEXT:  --- we need ot convert it to []Bytes and then to go struct --- unmarshal method
		// else we take resp.Body -> as a stream of text and send it to newDecoder
		responseBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, fmt.Errorf("Error reading response body: %w", err)
		}
		// caching result
		pokeClient.cache.Add(url, responseBytes)
	}
	return responseBytes, nil
}
