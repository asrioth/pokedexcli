package pokeapi

import (
	"testing"
)

func TestGetMapStrings(t *testing.T) {
	expectedAreas := []string{"canalave-city-area",
		"eterna-city-area",
		"pastoria-city-area",
		"sunyshore-city-area",
		"sinnoh-pokemon-league-area",
		"oreburgh-mine-1f",
		"oreburgh-mine-b1f",
		"valley-windworks-area",
		"eterna-forest-area",
		"fuego-ironworks-area",
		"mt-coronet-1f-route-207",
		"mt-coronet-2f",
		"mt-coronet-3f",
		"mt-coronet-exterior-snowfall",
		"mt-coronet-exterior-blizzard",
		"mt-coronet-4f",
		"mt-coronet-4f-small-room",
		"mt-coronet-5f",
		"mt-coronet-6f",
		"mt-coronet-1f-from-exterior"}
	acturalAreas, err := GetMapStrings(0, 19)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if len(acturalAreas) != len(expectedAreas) {
		t.Errorf("not enough arreas returned")
		return
	}
	for index := range expectedAreas {
		if acturalAreas[index] != expectedAreas[index] {
			t.Errorf("actual area '%v' did not match expected area '%v", acturalAreas[index], expectedAreas[index])
		}
	}
}

func TestGetPokemonForArea(t *testing.T) {
	actualPokemon, err := GetPokemonForArea("eterna-city-area")
	if err != nil {
		t.Error(err)
		return
	}
	expectedPokemon := []string{"psyduck", "golduck", "magikarp", "gyarados", "barboach", "whiscash"}
	for index := range expectedPokemon {
		if actualPokemon[index] != expectedPokemon[index] {
			t.Errorf("actual pokemon '%v' did not match expected pokemon '%v", actualPokemon[index], expectedPokemon[index])
		}
	}
}

func TestGetPokemonBaseXp(t *testing.T) {
	baseXp, err := GetPokemonBaseXp("squirtle")
	if err != nil {
		t.Error(err)
		return
	}
	expectedXp := 63
	if expectedXp != baseXp {
		t.Errorf("actual Xp %v not equal to expected Xp %v", baseXp, expectedXp)
	}
}

func TestGetPokemonStats(t *testing.T) {
	pokemonDescription, err := GetPokemonStats("pikachu")
	if err != nil {
		t.Error(err)
		return
	}
	if pokemonDescription.Hp != 35 {
		t.Errorf("worng Hp %v", pokemonDescription.Hp)
	}
	if pokemonDescription.Attack != 55 {
		t.Errorf("worng Attack %v", pokemonDescription.Attack)
	}
	if pokemonDescription.Defense != 40 {
		t.Errorf("worng Defense %v", pokemonDescription.Defense)
	}
	if pokemonDescription.SpecialAttack != 50 {
		t.Errorf("worng SpecialAttack %v", pokemonDescription.SpecialAttack)
	}
	if pokemonDescription.SpecialDefense != 50 {
		t.Errorf("worng SpecialDefense %v", pokemonDescription.SpecialDefense)
	}
	if pokemonDescription.Speed != 90 {
		t.Errorf("worng Speed %v", pokemonDescription.Speed)
	}
	if pokemonDescription.Height != 4 {
		t.Errorf("worng Height %v", pokemonDescription.Height)
	}
	if pokemonDescription.Weight != 60 {
		t.Errorf("worng Weight %v", pokemonDescription.Weight)
	}
	if pokemonDescription.Types[0] != "electric" {
		t.Errorf("worng Types %v", pokemonDescription.Types)
	}
}
