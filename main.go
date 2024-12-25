package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/asrioth/pokedexcli/pokeCache"
	"github.com/asrioth/pokedexcli/pokeapi"
)

type Config struct {
	Next     int
	Previous int
	data     *pokeCache.Cache
}

type CliCommand struct {
	name        string
	description string
	callback    func(*Config) error
	config      *Config
}

func initializeCommands() map[string]CliCommand {
	mapConfig := Config{0, 0, pokeCache.NewCache(time.Second * 5)}
	supportedCommands := map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			config:      &Config{0, 0, nil},
		},
		"help": {
			name:        "exit",
			description: "Displays a help message",
			callback:    commandHelp,
			config:      &Config{0, 0, nil},
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
	}
	return supportedCommands
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
		fmt.Println("no Cache")
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
		fmt.Println("no Cache")
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

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}
