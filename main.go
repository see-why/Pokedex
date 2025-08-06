package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

		// Look up command in registry
		commands := getCommands()
		if command, exists := commands[commandName]; exists {
			err := command.callback(config)
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
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
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

func commandMap(cfg *config) error {
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

func commandMapb(cfg *config) error {
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

type locationAreasResp struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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
