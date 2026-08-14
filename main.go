package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

/// panic last and first page and the moduel of netowrking

type cliCommand struct {
	name        string
	description string
	callback    func(config *config) error
}

type config struct {
	Next     *string
	Previous *string
}

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

var supportedCommands map[string]cliCommand

var cfg config

func init() { // init runs before main
	// assigning memory dynamically
	firstPage := "https://pokeapi.co/api/v2/location-area"
	cfg.Next = &firstPage // or cfg.Next = new(string)
	cfg.Previous = new(string)

	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func main() {

	// REPL program
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan() // reads until enter

		line := scanner.Text()
		cmdWords := cleanInput(line)
		if len(cmdWords) == 0 { // no command
			continue
		}

		// fmt.Printf("Your command was: %s\n", sliceStr[0])
		cmd, ok := supportedCommands[cmdWords[0]]
		if !ok {
			fmt.Println("Unknown Command")
			continue
		}

		err := cmd.callback(&cfg)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
	}

}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range supportedCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(config *config) error {

	// piping and receiving the response
	req, err := http.NewRequest("GET", *config.Next, nil)
	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "PokedexCLI/1.0 - shady")

	client := &http.Client{}
	resp, err := client.Do(req)

	// resp, err := http.Get(*config.Next)
	if err != nil {
		return fmt.Errorf("Failed to GET location areas: %w", err)
	}
	defer resp.Body.Close()

	//resp is RAW TEXT:  --- we need ot convert it to []Bytes and then to go struct --- unmarshal method
	// else we take resp.Body -> as a stream of text and send it to newDecoder

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}

	var locationAreaResponse LocationAreaResponse

	err = json.Unmarshal(responseBytes, &locationAreaResponse)
	if err != nil {
		return fmt.Errorf("Error unmarshaling: %w", err)
	}

	// printing
	for _, result := range locationAreaResponse.Results {
		fmt.Println(result.Name)
	}

	// next and prev
	*config.Previous = *config.Next
	if locationAreaResponse.Next != nil {
		// last page, if you try to go next, just stay
		*config.Next = *locationAreaResponse.Next
	}

	return nil
}

func commandMapb(config *config) error {

	// piping and receiving the response
	req, err := http.NewRequest("GET", *config.Previous, nil)
	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "PokedexCLI/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)

	// resp, err := http.Get(*config.Previous)
	if err != nil {
		return fmt.Errorf("Failed to GET location areas: %w", err)
	}
	defer resp.Body.Close()

	//resp is RAW TEXT:  --- we need ot convert it to []Bytes and then to go struct --- unmarshal method
	// else we take resp.Body -> as a stream of text and send it to newDecoder

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}

	var locationAreaResponse LocationAreaResponse

	err = json.Unmarshal(responseBytes, &locationAreaResponse)
	if err != nil {
		return fmt.Errorf("Error unmarshaling: %w", err)
	}

	// printing
	for _, result := range locationAreaResponse.Results {
		fmt.Println(result.Name)
	}

	// next and prev
	*config.Next = *config.Previous
	if locationAreaResponse.Previous != nil {
		*config.Previous = *locationAreaResponse.Previous
	}
	return nil
}
