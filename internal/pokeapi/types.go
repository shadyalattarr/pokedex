package pokeapi

import (
	"github.com/shadyalattarr/pokedex/internal/pokecache"
)

type Resource struct {
	Name string  `json:"name"`
	Url  *string `json:"url"`
}

type LocationAreaResponse struct {
	Count    int        `json:"count"` // json tags
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Resource `json:"results"`
}

type PokeApiClient struct {
	cache pokecache.Cache
}

type LocationAreaInformationResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod Resource `json:"encounter_method"`
		VersionDetails  []struct {
			Rate    int      `json:"rate"`
			Version Resource `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language Resource `json:"language"`
		Name     string   `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon        Resource `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int      `json:"chance"`
				ConditionValues []any    `json:"condition_values"`
				MaxLevel        int      `json:"max_level"`
				Method          Resource `json:"method"`
				MinLevel        int      `json:"min_level"`
				PokemonDetails  any      `json:"pokemon_details"`
			} `json:"encounter_details"`
			MaxChance int      `json:"max_chance"`
			Version   Resource `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
