package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/asrioth/pokedexcli/pokeCache"
	"github.com/asrioth/pokedexcli/pokeapi"
	"github.com/asrioth/pokedexcli/pokedexData"
)

type Config struct {
	Next     int
	Previous int
	args     []string
	data     *pokeCache.Cache
	pokedex  pokedexData.PokeDex
}

type CliCommand struct {
	name        string
	description string
	callback    func(*Config) error
	config      *Config
}

func initializeCommands() map[string]CliCommand {
	pokedex := pokedexData.PokeDex{CaughtPokemon: make(map[string]pokedexData.Pokemon)}
	mapConfig := Config{0, 0, nil, pokeCache.NewCache(time.Second * 5), pokedex}
	exploreArgs := make([]string, 1)
	supportedCommands := map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			config:      &Config{0, 0, nil, nil, pokedex},
		},
		"help": {
			name:        "exit",
			description: "Displays a help message",
			callback:    commandHelp,
			config:      &Config{0, 0, nil, nil, pokedex},
		},
		"map": {
			name:        "map",
			description: "Lists the next 20 location areas",
			callback:    commandMap,
			config:      &mapConfig,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists the previous 20 location areas",
			callback:    commandMapBack,
			config:      &mapConfig,
		},
		"explore": {
			name:        "explore",
			description: "Lists all pokemon in the area, takes an area name eg. explore canalave-city-area",
			callback:    commandExplore,
			config:      &Config{0, 0, exploreArgs, nil, pokedex},
		},
		"catch": {
			name:        "catch",
			description: "Attemps to catch named pokemon. If successful adds it to pokeDex",
			callback:    commandCatch,
			config:      &Config{0, 0, exploreArgs, nil, pokedex},
		},
		"inspect": {
			name:        "inspect",
			description: "Displays pokemon data if user has attemted to catch the pokemon before",
			callback:    commandInspect,
			config:      &Config{0, 0, exploreArgs, nil, pokedex},
		},
	}
	return supportedCommands
}

func catch(baseXp int) bool {
	catchRate := float64(baseXp) / 644.0
	catchChance := rand.Float64()
	return catchChance >= catchRate
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	supportedCommands := initializeCommands()
	for {
		fmt.Print("Pokedex > ")
		input.Scan()
		words := cleanInput(input.Text())
		runCommands(words, supportedCommands)
	}
}

func runCommands(words []string, supportedCommands map[string]CliCommand) {
	for i := 0; i < len(words); i++ {
		word := words[i]
		command, ok := supportedCommands[word]
		if !ok {
			fmt.Printf("%v not a valid command.\n All of input : %v must be valid commands or part of a valid command.\n", word, words)
			break
		}
		if i+len(command.config.args) >= len(words) {
			fmt.Printf("%v expects %v arguments and command has %v arguments.\n", word, len(command.config.args), len(words)-(i+1))
			break
		}
		for argI := 0; argI < len(command.config.args); argI++ {
			i++
			command.config.args[argI] = words[i]
		}
		err := command.callback(command.config)
		if err != nil {
			fmt.Printf("command returned error: %v\n", err)
			break
		}
	}
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	supportedCommands := initializeCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for cmdName, cmd := range supportedCommands {
		fmt.Printf("%v: %v\n", cmdName, cmd.description)
	}
	return nil
}

func commandMap(config *Config) error {
	var data []string
	data = config.data.GetRange(config.Next+1, config.Next+20)
	if data == nil {
		mapStrings, err := pokeapi.GetMapStrings(config.Next, config.Next+19)
		if err != nil {
			return err
		}
		data = mapStrings
		config.data.AddAll(config.Next+1, data)
	}

	for _, datum := range data {
		fmt.Println(datum)
	}
	config.Previous = config.Next
	config.Next += 20
	return nil
}

func commandMapBack(config *Config) error {
	if config.Previous <= 0 {
		fmt.Println("you're on the first page")
		return nil
	}
	var data []string
	data = config.data.GetRange(config.Previous-19, config.Previous)
	if data == nil {
		mapStrings, err := pokeapi.GetMapStrings(config.Previous-20, config.Previous-1)
		if err != nil {
			return err
		}
		data = mapStrings
		config.data.AddAll(config.Previous-19, data)
	}
	for _, datum := range data {
		fmt.Println(datum)
	}
	config.Next = config.Previous
	config.Previous -= 20
	return nil
}

func commandExplore(config *Config) error {
	fmt.Printf("Exploring %v...\n", config.args[0])
	pokemons, err := pokeapi.GetPokemonForArea(config.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemons {
		fmt.Printf(" - %v\n", pokemon)
	}
	return nil
}

func commandCatch(config *Config) error {
	name := config.args[0]
	baseXp, err := pokeapi.GetPokemonBaseXp(name)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", name)
	if catch(baseXp) {
		fmt.Printf("%v was caught!\n", name)
		config.pokedex.Catch(name, true)
	} else {
		fmt.Printf("%v escaped!\n", name)
		config.pokedex.Catch(name, false)
	}
	return nil
}

func commandInspect(config *Config) error {
	name := config.args[0]
	pokemon, ok := config.pokedex.GetPokemon(name)
	if !ok {
		return errors.New("need to try catching a pokemon before inspecting it")
	}
	if pokemon.Description.Height == -1 {
		pokemonDescription, err := pokeapi.GetPokemonStats(name)
		if err != nil {
			return nil
		}
		config.pokedex.AddDescription(name, pokemonDescription)
	}
	pokemon, _ = config.pokedex.GetPokemon(name)
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Successful Catches: %v\n", pokemon.CatchCount)
	fmt.Printf("Failed Catches: %v\n", pokemon.FailCatchCount)
	fmt.Printf("Height: %v\n", pokemon.Description.Height)
	fmt.Printf("Weight: %v\n", pokemon.Description.Weight)
	fmt.Println("Stats:")
	printStat("hp", pokemon.Description.Hp)
	printStat("attack", pokemon.Description.Attack)
	printStat("defense", pokemon.Description.Defense)
	printStat("special-attack", pokemon.Description.SpecialAttack)
	printStat("special-defense", pokemon.Description.SpecialDefense)
	printStat("speed", pokemon.Description.Speed)
	fmt.Println("Types:")
	for _, pokeType := range pokemon.Description.Types {
		fmt.Printf(" - %v\n", pokeType)
	}
	return nil
}

func printStat(statName string, stat int) {
	fmt.Printf(" - %v: %v\n", statName, stat)
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}
