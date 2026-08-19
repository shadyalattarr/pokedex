package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"math/rand"

	"github.com/shadyalattarr/pokedex/internal/pokeapi"
)

/// panic last and first page and the moduel of netowrking

type cliCommand struct {
	name        string
	description string
	callback    func(config *config, args ...string) error
}

type config struct {
	Next     *string
	Previous *string
}

var supportedCommands map[string]cliCommand
var pokedex map[string]pokeapi.PokemonResponse

var cfg config
var pokeClient pokeapi.PokeApiClient

func init() { // init runs before main
	// assigning memory dynamically
	firstPage := "https://pokeapi.co/api/v2/location-area"
	cfg.Next = &firstPage // or cfg.Next = new(string)
	cfg.Previous = new(string)

	pokeClient = pokeapi.NewClient(5 * time.Second)

	pokedex = map[string]pokeapi.PokemonResponse{} // initialized empty

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
		"explore": {
			name:        "explore",
			description: "Explores the given location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect information about pokemon in your pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List the pokemon ou have in your pokedex",
			callback:    commandPokedex,
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
		// problem is first word is always loc area now
		err := cmd.callback(&cfg, cmdWords[1:]...)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
	}

}

func commandExit(config *config, args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("command exit requires zero arguments.")
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("command help requires zero arguments.")
	}
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range supportedCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(config *config, args ...string) error {
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

func commandMapb(config *config, args ...string) error {
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

func commandExplore(config *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command explore requires only one argument: the location area:\n example: explore <location-area>")
	}
	locArea := args[0]

	url := "https://pokeapi.co/api/v2/location-area/" + locArea
	fmt.Printf("Exploring %s ...\n", locArea)
	locationAreaInfoResp, err := pokeClient.GetLocationAreaInfo(url)
	if err != nil {
		return fmt.Errorf("Error with GetLocationAreaInfo: %w", err)
	}

	// printing
	if len(locationAreaInfoResp.PokemonEncounters) != 0 {
		fmt.Println("Found Pokemon:")
	} else {
		fmt.Println("No Pokemon were found!")
	}

	for _, pokeEncounter := range locationAreaInfoResp.PokemonEncounters {
		fmt.Println("- ", pokeEncounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command catch requires only one argument: the pokemon:\n example: explore <pokemon-name>")
	}
	pokemon := args[0]
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon)

	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon

	pokemonResp, err := pokeClient.GetPokemonInfo(url)
	if err != nil {
		return fmt.Errorf("Error with GetPokemonInfo: %w", err)
	}

	// printing
	baseXP := pokemonResp.BaseExperience
	randomRoll := rand.Intn(baseXP)
	threshold := 40

	if randomRoll < threshold {
		pokedex[pokemon] = pokemonResp
		fmt.Printf("%v was caught!\n", pokemon)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%v escaped!\n", pokemon)
	}
	return nil
}

func commandInspect(config *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command inspect requires only one argument: the pokemon:\n example: inspect <pokemon-name>")
	}
	pokemon := args[0]
	pokemonResp, ok := pokedex[pokemon]
	if !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		printPokemonInfo(pokemonResp)
	}
	return nil
}

func printPokemonInfo(pokemonResp pokeapi.PokemonResponse) {
	fmt.Printf("Name: %v\n", pokemonResp.Name)
	fmt.Printf("Height: %v\n", pokemonResp.Height)
	fmt.Printf("Weight: %v\n", pokemonResp.Weight)
	fmt.Println("Stats:")
	for _, statistic := range pokemonResp.Stats {
		fmt.Printf("\t-%v: %v\n", statistic.Stat.Name, statistic.BaseStat)
	}
	fmt.Println("Types:")
	for _, type_ := range pokemonResp.Types {
		fmt.Printf("\t-%v\n", type_.Type.Name)
	}
}

func commandPokedex(config *config, args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("command pokedex requires zero arguments.")
	}

	if len(pokedex) == 0 {
		fmt.Println("You have no pokemon in your pokedex!")
	} else {
		fmt.Println("Your Pokedex:")
		for k := range pokedex {
			fmt.Printf("\t- %v\n", k)
		}
	}
	return nil
}
