package pokeapi

import "testing"

func TestPokeMap(t *testing.T) {
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
