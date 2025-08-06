package main

import (
	"testing"
	"time"

	"github.com/see-why/Pokedex/internal/pokecache"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "singleword",
			expected: []string{"singleword"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// Check the length of the actual slice against the expected slice
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned slice of length %d, expected %d",
				c.input, len(actual), len(c.expected))
			continue
		}

		// Check each word in the slice
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q) returned %q at index %d, expected %q",
					c.input, word, i, expectedWord)
			}
		}
	}
}

func TestGetCommands(t *testing.T) {
	commands := getCommands()

	// Test that we have exactly the expected number of commands
	expectedCount := 5
	if len(commands) != expectedCount {
		t.Errorf("Expected %d commands, got %d", expectedCount, len(commands))
	}

	// Test that all expected commands are present
	expectedCommands := []string{"help", "exit", "map", "mapb", "explore"}

	for _, expectedCmd := range expectedCommands {
		_, exists := commands[expectedCmd]
		if !exists {
			t.Errorf("Expected command %q to be registered", expectedCmd)
		}
	}

	// Test that no unexpected commands are present
	for cmdName := range commands {
		found := false
		for _, expectedCmd := range expectedCommands {
			if cmdName == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected command %q found in registry", cmdName)
		}
	}

	// Test specific command properties
	if cmd, exists := commands["map"]; exists {
		if cmd.name != "map" {
			t.Errorf("map command name = %q, expected %q", cmd.name, "map")
		}
		if cmd.description != "Displays the names of 20 location areas" {
			t.Errorf("map command description = %q, expected %q", cmd.description, "Displays the names of 20 location areas")
		}
		if cmd.callback == nil {
			t.Error("map command callback is nil")
		}
	}

	if cmd, exists := commands["mapb"]; exists {
		if cmd.name != "mapb" {
			t.Errorf("mapb command name = %q, expected %q", cmd.name, "mapb")
		}
		if cmd.description != "Displays the previous 20 location areas" {
			t.Errorf("mapb command description = %q, expected %q", cmd.description, "Displays the previous 20 location areas")
		}
		if cmd.callback == nil {
			t.Error("mapb command callback is nil")
		}
	}
}

func TestCommandMapb_FirstPage(t *testing.T) {
	// Test mapb command when on first page (previousLocationURL is nil)
	cfg := &config{
		pokeapiClient:       pokecache.NewCache(5 * time.Minute),
		nextLocationURL:     "https://pokeapi.co/api/v2/location-area",
		previousLocationURL: nil,
	}

	// This should not return an error and should print "you're on the first page"
	err := commandMapb(cfg)
	if err != nil {
		t.Errorf("commandMapb returned unexpected error: %v", err)
	}

	// Config should remain unchanged when on first page
	if cfg.nextLocationURL != "https://pokeapi.co/api/v2/location-area" {
		t.Errorf("nextLocationURL changed unexpectedly: %q", cfg.nextLocationURL)
	}
	if cfg.previousLocationURL != nil {
		t.Errorf("previousLocationURL should remain nil on first page, got: %v", cfg.previousLocationURL)
	}
}

func TestConfig(t *testing.T) {
	// Test config struct initialization
	cfg := &config{
		pokeapiClient:       pokecache.NewCache(5 * time.Minute),
		nextLocationURL:     "https://pokeapi.co/api/v2/location-area",
		previousLocationURL: nil,
	}

	if cfg.nextLocationURL != "https://pokeapi.co/api/v2/location-area" {
		t.Errorf("nextLocationURL = %q, expected %q", cfg.nextLocationURL, "https://pokeapi.co/api/v2/location-area")
	}

	if cfg.previousLocationURL != nil {
		t.Errorf("previousLocationURL should be nil initially, got: %v", cfg.previousLocationURL)
	}

	// Test updating previousLocationURL
	testURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	cfg.previousLocationURL = &testURL

	if cfg.previousLocationURL == nil {
		t.Error("previousLocationURL should not be nil after assignment")
	}

	if *cfg.previousLocationURL != testURL {
		t.Errorf("previousLocationURL = %q, expected %q", *cfg.previousLocationURL, testURL)
	}
}

func TestCliCommandStruct(t *testing.T) {
	// Test that the cliCommand struct works as expected
	testCallback := func(cfg *config, args ...string) error {
		return nil
	}

	cmd := cliCommand{
		name:        "test",
		description: "Test command",
		callback:    testCallback,
	}

	if cmd.name != "test" {
		t.Errorf("command name = %q, expected %q", cmd.name, "test")
	}

	if cmd.description != "Test command" {
		t.Errorf("command description = %q, expected %q", cmd.description, "Test command")
	}

	if cmd.callback == nil {
		t.Error("command callback should not be nil")
	}

	// Test calling the callback
	cfg := &config{
		pokeapiClient: pokecache.NewCache(5 * time.Minute),
	}
	err := cmd.callback(cfg)
	if err != nil {
		t.Errorf("callback returned unexpected error: %v", err)
	}
}

func TestCommandExplore_NoArgs(t *testing.T) {
	cfg := &config{
		pokeapiClient: pokecache.NewCache(5 * time.Minute),
	}
	
	err := commandExplore(cfg)
	if err == nil {
		t.Error("expected error when no arguments provided")
	}
	
	expectedError := "you must provide a location area name"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}
