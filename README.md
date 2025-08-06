# Pokedex CLI

A command-line Pokemon exploration and collection tool built in Go. Explore the Pokemon world, catch Pokemon, and build your own Pokedex collection using data from the [PokeAPI](https://pokeapi.co/).

## Features

### 🌍 World Exploration

- **Location Navigation**: Browse through Pokemon location areas with pagination
- **Area Exploration**: Discover which Pokemon can be found in specific locations
- **Intelligent Caching**: Fast response times with built-in HTTP response caching

### 🎮 Pokemon Interaction  

- **Pokemon Catching**: Attempt to catch Pokemon with randomized success rates based on difficulty
- **Collection Management**: Keep track of all Pokemon you've successfully caught
- **Pokemon Inspection**: View detailed stats, types, height, and weight of caught Pokemon

### 🖥️ Interactive Experience

- **REPL Interface**: Command-line interface with persistent state
- **Help System**: Built-in command documentation and usage examples
- **Real-time Feedback**: Immediate responses and error handling

## Installation

```bash
# Clone the repository
git clone https://github.com/see-why/Pokedex.git
cd Pokedex

# Build the application
go build

# Run the Pokedex
./Pokedex
```

## Commands

| Command | Arguments | Description |
|---------|-----------|-------------|
| `help` | none | Display all available commands and their descriptions |
| `map` | none | Show the next 20 location areas |
| `mapb` | none | Show the previous 20 location areas |
| `explore` | `<location-area>` | List all Pokemon that can be found in the specified location |
| `catch` | `<pokemon-name>` | Attempt to catch a Pokemon (success varies by Pokemon difficulty) |
| `inspect` | `<pokemon-name>` | View detailed information about a caught Pokemon |
| `pokedex` | none | Display a list of all Pokemon you have caught |
| `exit` | none | Exit the application |

## Usage Examples

```bash
# Start the Pokedex
$ ./Pokedex
Pokedex > 

# Get help
Pokedex > help
Welcome to the Pokedex!
Usage:
help: Displays a help message
map: Displays the names of 20 location areas
...

# Explore the world
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...

# Discover Pokemon in an area
Pokedex > explore pallet-town-area
Exploring pallet-town-area...
Found Pokemon:
 - bulbasaur
 - charmander
 - squirtle
 - pikachu
...

# Catch Pokemon
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!
You may now inspect it with the inspect command.

# View your collection
Pokedex > pokedex
Your Pokedex:
 - pikachu

# Inspect caught Pokemon
Pokedex > inspect pikachu
Name: pikachu
Height: 4
Weight: 60
Stats:
  -hp: 35
  -attack: 55
  -defense: 40
  -special-attack: 50
  -special-defense: 50
  -speed: 90
Types:
  - electric
```

## Project Structure

```
Pokedex/
├── main.go              # Main application with REPL and command implementations
├── repl_test.go         # Comprehensive test suite
├── go.mod              # Go module definition
├── internal/
│   └── pokecache/
│       ├── cache.go     # HTTP response caching with TTL
│       └── cache_test.go# Cache testing
└── README.md           # This file
```

## Architecture

### Core Components

- **REPL Loop**: Interactive command-line interface with command registry
- **Command System**: Modular command architecture with consistent error handling
- **HTTP Client**: Integration with PokeAPI for real-time Pokemon data
- **Caching Layer**: Thread-safe HTTP response caching with automatic cleanup
- **State Management**: Persistent Pokemon collection during session

### Technical Details

- **Language**: Go 1.24+
- **API Integration**: RESTful calls to [PokeAPI](https://pokeapi.co/)
- **Concurrency**: Thread-safe operations with mutex protection
- **Testing**: Comprehensive test coverage with table-driven tests
- **Error Handling**: Graceful error handling with user-friendly messages

## API Integration

The Pokedex integrates with the following PokeAPI endpoints:

- `/location-area` - Location area listings with pagination
- `/location-area/{area}` - Detailed area information and Pokemon encounters  
- `/pokemon/{name}` - Individual Pokemon data including stats and types

## Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test file
go test ./repl_test.go
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the MIT License.

## Acknowledgments

- [PokeAPI](https://pokeapi.co/) - Free RESTful Pokemon API
- The Pokemon Company - For creating the Pokemon universe
- Go Community - For excellent tooling and libraries
