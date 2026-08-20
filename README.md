# Pokedex CLI

A command-line REPL (Read-Eval-Print Loop) application written in Go that interacts with the [PokéAPI](https://pokeapi.co/). This project was built as a hands-on dive into backend engineering concepts and the Go programming language. This was a guided project, as part of the [boot.dev](https://www.boot.dev/dashboard) Backend Developer path.

## Learning Goals Achieved
* Built a custom CLI tool to seamlessly interact with a back-end server.
* Managed HTTP network requests and handled API responses in Go.
* Parsed complex JSON data structures into native Go structs.
* Gained practical experience with local Go development environments and tooling.

## Features & Commands

Once the REPL is running, you can use the following commands to navigate the Pokemon world[cite: 1]:

* `help`: Displays a help message with available commands.
* `map` / `mapb`: Paginates through the next or previous 20 location areas in the Pokémon world.
* `explore <location-area>`: Lists all the Pokémon that can be encountered in a specific location.
* `catch <pokemon-name>`: Attempts to catch a Pokémon. The catch rate is randomized and scales based on the Pokémon's base experience.
* `inspect <pokemon-name>`: Displays the height, weight, stats, and types of a Pokémon you have successfully caught.
* `pokedex`: Lists all the Pokémon currently stored in your local Pokedex.
* `exit`: Closes the application.

## Installation & Usage

1. Ensure you have Go installed (version 1.26 or higher recommended).
2. Clone the repository:
   ```bash
   git clone https://github.com/shadyalattarr/pokedex.git
   ```
3. Navigate into the project directory and run the application:
    ```bash
    go run .
    ``` 
    or do this instead:
    ```bash 
    ./pokedex
    ```
