package pokedexData

import "testing"

func TestPokedexCatch(t *testing.T) {
	pokedex := PokeDex{CaughtPokemon: make(map[string]Pokemon)}
	name := "testachu"
	pokedex.Catch(name, true)
	if pokedex.CaughtPokemon[name].CatchCount != 1 {
		t.Errorf("Catch count (%v) not incremented on successful catch", pokedex.CaughtPokemon[name].CatchCount)
	}
	if pokedex.CaughtPokemon[name].FailCatchCount != 0 {
		t.Errorf("Fail count (%v) incremented on successful catch when it shouldn't be", pokedex.CaughtPokemon[name].CatchCount)
	}
	pokedex.Catch(name, false)
	if pokedex.CaughtPokemon[name].CatchCount != 1 {
		t.Errorf("Catch count (%v) incremented on failed catch when it shouldn't be", pokedex.CaughtPokemon[name].CatchCount)
	}
	if pokedex.CaughtPokemon[name].FailCatchCount != 1 {
		t.Errorf("Fail count (%v) not incremented on failed catch", pokedex.CaughtPokemon[name].CatchCount)
	}
}

func TestPokedexSaveLoad(t *testing.T) {
	pokedex := PokeDex{CaughtPokemon: make(map[string]Pokemon)}
	name := "testachu"
	pokedex.Catch(name, true)
	pokedex.Save()
	pokedex = PokeDex{}
	pokedex.Load()
	pokemon, ok := pokedex.CaughtPokemon[name]
	if !ok {
		t.Error("pokemon failed to load")
	}
	if pokemon.Name != name {
		t.Error("load incorrectly stored pokemon")
	}
}
