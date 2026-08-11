package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var supportedCommands map[string]cliCommand

func init() { // init runs before main
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

		err := cmd.callback()
		if err != nil {
			fmt.Println("Error with the command!")
			continue
		}
	}

}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range supportedCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}
