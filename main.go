package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/shadyalattarr/pokedex/internal/pokeapi"
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

var supportedCommands map[string]cliCommand

var cfg config
var pokeClient pokeapi.PokeApiClient

func init() { // init runs before main
	// assigning memory dynamically
	firstPage := "https://pokeapi.co/api/v2/location-area"
	cfg.Next = &firstPage // or cfg.Next = new(string)
	cfg.Previous = new(string)

	pokeClient = pokeapi.NewClient(5 * time.Second)

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
	loc_area_response, err := pokeClient.GetLocationAreas(*config.Next)
	if err != nil {
		return fmt.Errorf("Error with GetLocationAreas: %w", err)
	}
	// printing
	for _, result := range loc_area_response.Results {
		fmt.Println(result.Name)
	}

	// next and prev
	*config.Previous = *config.Next
	if loc_area_response.Next != nil {
		// last page, if you try to go next, just stay
		*config.Next = *loc_area_response.Next
	}

	return nil
}

func commandMapb(config *config) error {
	loc_area_response, err := pokeClient.GetLocationAreas(*config.Previous)
	if err != nil {
		return fmt.Errorf("Error with GetLocationAreas: %w", err)
	}

	// printing
	for _, result := range loc_area_response.Results {
		fmt.Println(result.Name)
	}

	// next and prev
	*config.Next = *config.Previous
	if loc_area_response.Previous != nil {
		*config.Previous = *loc_area_response.Previous
	}
	return nil
}
