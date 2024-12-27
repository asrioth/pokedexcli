package pokedexData

import "github.com/asrioth/pokedexcli/pokeapi"

const PokedexPath string = "pokedexData/pokedex.json"

type Pokemon struct {
	Name           string                     `json:"name"`
	Description    pokeapi.PokemonDescription `json:"description"`
	CatchCount     int                        `json:"catch_count"`
	FailCatchCount int                        `json:"fail_catch_count"`
}

type PokeDex struct {
	CaughtPokemon map[string]Pokemon `json:"caught_pokemon"`
}

func (P PokeDex) GetID() int {
	return 0
}

func (P PokeDex) GetName() string {
	return ""
}

func (P PokeDex) Save() {
	pArray := []PokeDex{P}
	pokeapi.CachePokeData(pArray, PokedexPath, false)
}

func (P PokeDex) Catch(name string, caught bool) Pokemon {
	pokemon, ok := P.CaughtPokemon[name]
	if !ok {
		pokemon = Pokemon{Name: name, CatchCount: 0, FailCatchCount: 0, Description: pokeapi.PokemonDescription{Height: -1}}
	}
	if caught {
		pokemon.CatchCount += 1
	} else {
		pokemon.FailCatchCount += 1
	}
	P.CaughtPokemon[name] = pokemon
	return pokemon
}

func (P PokeDex) GetPokemon(name string) (Pokemon, bool) {
	pokemon, ok := P.CaughtPokemon[name]
	return pokemon, ok
}

func (P PokeDex) AddDescription(name string, description pokeapi.PokemonDescription) {
	pokemon, ok := P.CaughtPokemon[name]
	if ok {
		pokemon.Description = description
		P.CaughtPokemon[name] = pokemon
	}
}
