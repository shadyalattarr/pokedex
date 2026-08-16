package pokeapi

import (
	"github.com/shadyalattarr/pokedex/internal/pokecache"
)

type LocationArea struct {
	Name string  `json:"name"`
	Url  *string `json:"url"`
}

type LocationAreaResponse struct {
	Count    int            `json:"count"` // json tags
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type PokeApiClient struct {
	cache pokecache.Cache
}
