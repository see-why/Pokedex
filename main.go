package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/see-why/Pokedex/internal/pokecache"
)

func main() {
	config := &config{
		pokeapiClient:       pokecache.NewCache(5 * time.Minute),
		nextLocationURL:     "https://pokeapi.co/api/v2/location-area",
		previousLocationURL: nil,
		caughtPokemon:       make(map[string]Pokemon),
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			// No more input available (EOF or error)
			break
		}
		input := scanner.Text()

		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		// Look up command in registry
		commands := getCommands()
		if command, exists := commands[commandName]; exists {
			err := command.callback(config, args...)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

type config struct {
	pokeapiClient       pokecache.Cache
	nextLocationURL     string
	previousLocationURL *string
	caughtPokemon       map[string]Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
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
			description: "Displays the names of 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon",
			callback:    commandInspect,
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()

	return nil
}

func commandMap(cfg *config, args ...string) error {
	locationAreas, err := getLocationAreas(cfg, cfg.nextLocationURL)
	if err != nil {
		return err
	}

	// Update config with new URLs
	cfg.nextLocationURL = locationAreas.Next
	cfg.previousLocationURL = locationAreas.Previous

	// Print all location area names
	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.previousLocationURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreas, err := getLocationAreas(cfg, *cfg.previousLocationURL)
	if err != nil {
		return err
	}

	// Update config with new URLs
	cfg.nextLocationURL = locationAreas.Next
	cfg.previousLocationURL = locationAreas.Previous

	// Print all location area names
	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a location area name")
	}

	locationAreaName := args[0]
	fmt.Printf("Exploring %s...\n", locationAreaName)

	locationArea, err := getLocationArea(cfg, locationAreaName)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, enc := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a Pokemon name")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := getPokemon(cfg, pokemonName)
	if err != nil {
		return err
	}

	// Use base experience to determine catch difficulty
	// Higher base experience = harder to catch
	const maxCatchChance = 50 // Base 50% chance for Pokemon with 0 base experience
	catchChance := maxCatchChance
	if pokemon.BaseExperience > 0 {
		// Reduce catch chance based on base experience, minimum 5% chance
		catchChance = max(5, maxCatchChance-pokemon.BaseExperience/10)
	}

	// Generate random number between 1-100
	if rand.Intn(100)+1 <= catchChance {
		cfg.caughtPokemon[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
		fmt.Printf("You may now inspect it with the inspect command.\n")
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a Pokemon name")
	}

	pokemonName := args[0]

	// Check if the Pokemon has been caught
	pokemon, exists := cfg.caughtPokemon[pokemonName]
	if !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	// Display Pokemon information
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

type locationAreasResp struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type locationAreaResp struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	GameIndex         int    `json:"game_index"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func getLocationAreas(cfg *config, pageURL string) (locationAreasResp, error) {
	// Check if we have the data in cache
	if val, ok := cfg.pokeapiClient.Get(pageURL); ok {
		fmt.Printf("Using cached data for %s\n", pageURL)
		locationAreasResponse := locationAreasResp{}
		err := json.Unmarshal(val, &locationAreasResponse)
		if err != nil {
			return locationAreasResp{}, err
		}
		return locationAreasResponse, nil
	}

	fmt.Printf("Making HTTP request to %s\n", pageURL)
	res, err := http.Get(pageURL)
	if err != nil {
		return locationAreasResp{}, err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return locationAreasResp{}, err
	}

	locationAreasResponse := locationAreasResp{}
	err = json.Unmarshal(dat, &locationAreasResponse)
	if err != nil {
		return locationAreasResp{}, err
	}

	// Add to cache
	cfg.pokeapiClient.Add(pageURL, dat)

	return locationAreasResponse, nil
}

func getLocationArea(cfg *config, locationAreaName string) (locationAreaResp, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + locationAreaName

	// Check if we have the data in cache
	if val, ok := cfg.pokeapiClient.Get(url); ok {
		fmt.Printf("Using cached data for %s\n", url)
		locationAreaResponse := locationAreaResp{}
		err := json.Unmarshal(val, &locationAreaResponse)
		if err != nil {
			return locationAreaResp{}, err
		}
		return locationAreaResponse, nil
	}

	fmt.Printf("Making HTTP request to %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		return locationAreaResp{}, err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return locationAreaResp{}, err
	}

	locationAreaResponse := locationAreaResp{}
	err = json.Unmarshal(dat, &locationAreaResponse)
	if err != nil {
		return locationAreaResp{}, err
	}

	// Add to cache
	cfg.pokeapiClient.Add(url, dat)

	return locationAreaResponse, nil
}

func getPokemon(cfg *config, pokemonName string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	// Check if we have the data in cache
	if val, ok := cfg.pokeapiClient.Get(url); ok {
		fmt.Printf("Using cached data for %s\n", url)
		pokemonResponse := Pokemon{}
		err := json.Unmarshal(val, &pokemonResponse)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemonResponse, nil
	}

	fmt.Printf("Making HTTP request to %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	pokemonResponse := Pokemon{}
	err = json.Unmarshal(dat, &pokemonResponse)
	if err != nil {
		return Pokemon{}, err
	}

	// Add to cache
	cfg.pokeapiClient.Add(url, dat)

	return pokemonResponse, nil
}

func cleanInput(text string) []string {
	// Trim leading and trailing whitespace and convert to lowercase
	cleaned := strings.ToLower(strings.TrimSpace(text))

	// If the cleaned string is empty, return empty slice
	if cleaned == "" {
		return []string{}
	}

	// Split by whitespace
	words := strings.Fields(cleaned)

	return words
}
