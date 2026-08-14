package pokeapi

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
