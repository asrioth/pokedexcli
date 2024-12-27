package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// const pokeApiLocationUrl = "https://pokeapi.co/api/v2/location/%v/"
const pokeApiAreaUrl = "https://pokeapi.co/api/v2/location-area/%v/"
const pokeApiPokemonUrl = "https://pokeapi.co/api/v2/pokemon/%v/"

// const locationsPath = "pokeLocations.json"
const areasPath = "pokeapi/pokeAreas.json"
const pokemonPath = "pokeapi/pokemon/.json"

func checkPokeDataCache[PDT PokeDataType](minIndex, maxIndex int, filePath string) ([]PDT, error) {
	pokeDataFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer pokeDataFile.Close()
	var allPokeData []PDT
	decoder := json.NewDecoder(pokeDataFile)
	if err := decoder.Decode(&allPokeData); err != nil {
		return nil, err
	}
	var pokeData []PDT
	for _, pokeDatum := range allPokeData {
		if pokeDatum.GetID()-1 >= minIndex && pokeDatum.GetID()-1 <= maxIndex {
			pokeData = append(pokeData, pokeDatum)
		}
	}

	if len(pokeData) == 0 {
		return nil, errors.New("no pokeData in index range")
	}
	return pokeData, nil
}

func checkPokeDataCacheByName[PDT PokeDataType](name string, filePath string) (PDT, error) {
	pokeDataFile, err := os.Open(filePath)
	if err != nil {
		var pdt PDT
		return pdt, err
	}
	defer pokeDataFile.Close()
	var allPokeData []PDT
	decoder := json.NewDecoder(pokeDataFile)
	if err := decoder.Decode(&allPokeData); err != nil {
		var pdt PDT
		return pdt, err
	}
	var pokeDatumMatch PDT
	pokeMatch := false
	for _, pokeDatum := range allPokeData {
		if pokeDatum.GetName() == name {
			pokeDatumMatch = pokeDatum
			pokeMatch = true
			break
		}
	}

	if !pokeMatch {
		var pdt PDT
		return pdt, errors.New("no pokeData in index range")
	}
	return pokeDatumMatch, nil
}

func GetPokeDatum[PDT PokeDataType](id int, url, filePath string, cacheChecked bool) (PDT, error) {
	if !cacheChecked {
		cachePokeData, err := checkPokeDataCache[PDT](id-1, id-1, filePath)
		if err == nil {
			return cachePokeData[0], nil
		}
	}
	currentUrl := fmt.Sprintf(url, id)
	getResult, err := http.Get(currentUrl)
	if err != nil {
		var empty PDT
		return empty, err
	}
	defer getResult.Body.Close()
	var pokeDatum PDT
	decoder := json.NewDecoder(getResult.Body)
	if err := decoder.Decode(&pokeDatum); err != nil {
		var empty PDT
		return empty, err
	}
	return pokeDatum, nil
}

func GetPokeDatumByName[PDT PokeDataType](name string, url, filePath string, cacheChecked bool) (PDT, error) {
	if !cacheChecked {
		pokeDatum, err := checkPokeDataCacheByName[PDT](name, filePath)
		if err == nil {
			return pokeDatum, nil
		}
	}
	currentUrl := fmt.Sprintf(url, name)
	getResult, err := http.Get(currentUrl)
	if err != nil {
		var empty PDT
		return empty, err
	}
	defer getResult.Body.Close()
	var pokeDatum PDT
	decoder := json.NewDecoder(getResult.Body)
	if err := decoder.Decode(&pokeDatum); err != nil {
		var empty PDT
		return empty, err
	}
	return pokeDatum, nil
}

func getPokeIdByName(name, filePath string) int {
	pokeId := 0
	for i := 0; ; i++ {
		currentPath := getCurrentPath(filePath, strconv.Itoa(i))
		currentPath = getCurrentPath(currentPath, "id")
		pokeNameId, err := checkPokeDataCacheByName[PokeNameId](name, currentPath)
		if err != nil {
			if !errors.Is(err, errors.New("no pokeData in index range")) {
				break
			}
			continue
		}
		pokeId = pokeNameId.Id
		break
	}
	return pokeId
}

func GetMissingPokeData[PDT PokeDataType](missingPokeData []int, url, filePath string) ([]PDT, error) {
	pokeData := make([]PDT, len(missingPokeData))
	for index, id := range missingPokeData {
		pokeDatum, err := GetPokeDatum[PDT](id, url, filePath, true)
		if err != nil {
			return nil, err
		}
		pokeData[index] = pokeDatum
	}
	return pokeData, nil
}

func CachePokeData[PDT PokeDataType](pokeData []PDT, filePath string, hasIdmap bool) error {
	var cachedPokeData []PDT
	var cachedPokeDataNames []PokeNameId
	fileIdPath := getCurrentPath(filePath, "id")
	pokeDataFile, err := os.Open(filePath)
	if err == nil {
		decoder := json.NewDecoder(pokeDataFile)
		if err := decoder.Decode(&cachedPokeData); err != nil {
			pokeDataFile.Close()
			return err
		}
		pokeDataFile.Close()
		if hasIdmap {
			pokeNameIdFile, err := os.Open(fileIdPath)
			if err != nil {
				return err
			}
			decoder = json.NewDecoder(pokeNameIdFile)
			if err := decoder.Decode(&cachedPokeDataNames); err != nil {
				pokeDataFile.Close()
				return err
			}
			pokeNameIdFile.Close()
		}

	}

	cachedPokeData = append(cachedPokeData, pokeData...)
	if hasIdmap {
		for _, pokeDatum := range pokeData {
			pokeNameId := PokeNameId{pokeDatum.GetName(), pokeDatum.GetID()}
			cachedPokeDataNames = append(cachedPokeDataNames, pokeNameId)
		}
	}

	if errors.Is(err, fs.ErrNotExist) {
		dirSep := strings.LastIndex(filePath, "/")
		os.MkdirAll(filePath[:dirSep], 0755)
	}

	pokeDataFile, err = os.Create(filePath)
	if err != nil {
		return err
	}
	defer pokeDataFile.Close()
	encoder := json.NewEncoder(pokeDataFile)
	if err := encoder.Encode(cachedPokeData); err != nil {
		return err
	}
	if hasIdmap {
		pokeNameDataFile, err := os.Create(fileIdPath)
		if err != nil {
			return err
		}
		defer pokeNameDataFile.Close()
		encoder = json.NewEncoder(pokeNameDataFile)
		if err := encoder.Encode(cachedPokeDataNames); err != nil {
			return err
		}
	}
	return nil
}

func GetPokeData[PDT PokeDataType](minIndex, maxIndex int, url, filePath string) ([]PDT, error) {
	cachedPokeData, err := checkPokeDataCache[PDT](minIndex, maxIndex, filePath)
	if err == nil {
		if len(cachedPokeData) == maxIndex-minIndex+1 {
			return cachedPokeData, nil
		}
		allIds := make([]int, maxIndex-minIndex+1)
		for _, pokeDatum := range cachedPokeData {
			curId := pokeDatum.GetID()
			allIds[curId-1] = curId
		}
		var missingPokeDataIds []int
		for index, id := range allIds {
			if id == 0 {
				missingPokeDataIds = append(missingPokeDataIds, index+1)
			}
		}
		missingPokeData, err := GetMissingPokeData[PDT](missingPokeDataIds, url, filePath)
		if err != nil {
			return nil, err
		}
		if err := CachePokeData(missingPokeData, filePath, true); err != nil {
			return nil, err
		}
	}
	var pokeData []PDT
	for i := minIndex + 1; i < maxIndex+2; i++ {
		pokeDatum, err := GetPokeDatum[PDT](i, url, filePath, true)
		if err != nil {
			return nil, err
		}
		pokeData = append(pokeData, pokeDatum)
	}
	err = CachePokeData(pokeData, filePath, true)
	if err != nil {
		return nil, err
	}
	return pokeData, nil
}

/*func mapifyPokeData[PDT PokeDataType](pokeData []PDT) map[int]PDT {
	pokeMap := make(map[int]PDT)
	for _, pokeDatum := range pokeData {
		pokeMap[pokeDatum.GetID()] = pokeDatum
	}
	return pokeMap
}*/

func getCurrentPath(filepath string, fileId string) string {
	pathExtension := strings.Split(filepath, ".")
	currentPath := pathExtension[0] + fileId + "." + pathExtension[1]
	return currentPath
}

func GetMapStrings(minIndex, maxIndex int) ([]string, error) {
	currentAreaPath := getCurrentPath(areasPath, strconv.Itoa(minIndex/20))
	path, err := filepath.Localize(currentAreaPath)
	if err != nil {
		return nil, err
	}
	areas, err := GetPokeData[Area](minIndex, maxIndex, pokeApiAreaUrl, path)
	if err != nil {
		return nil, err
	}
	areaNames := make([]string, 20)
	for index := range areas {
		areaNames[index] = areas[index].Name
	}
	return areaNames, nil
}

func GetPokemonForArea(areaName string) ([]string, error) {
	id := getPokeIdByName(areaName, areasPath)
	var pokemonNames []string
	var area Area
	var err error
	if id != 0 {
		area, err = GetPokeDatum[Area](id, pokeApiAreaUrl, getCurrentPath(areasPath, strconv.Itoa(id/20)), false)
	} else {
		area, err = GetPokeDatumByName[Area](areaName, pokeApiAreaUrl, "", true)
	}
	if err != nil {
		return nil, err
	}
	for _, ecnounter := range area.PokemonEncounters {
		pokemonNames = append(pokemonNames, ecnounter.AreaPokemon.Name)
	}
	return pokemonNames, nil
}

func GetPokemonBaseXp(name string) (int, error) {
	currentpath := getCurrentPath(pokemonPath, name)
	pokemon, err := GetPokeDatumByName[Pokemon](name, pokeApiPokemonUrl, currentpath, false)
	if err != nil {
		return 0, err
	}
	pokemonData := []Pokemon{pokemon}
	CachePokeData(pokemonData, currentpath, false)
	return pokemon.BaseExperience, nil
}

func GetPokemonStats(name string) (PokemonDescription, error) {
	currentpath := getCurrentPath(pokemonPath, name)
	pokemon, err := GetPokeDatumByName[Pokemon](name, pokeApiPokemonUrl, currentpath, false)
	if err != nil {
		return PokemonDescription{}, err
	}
	pokemonData := []Pokemon{pokemon}
	CachePokeData(pokemonData, currentpath, false)

	var pokeTypes []string
	for _, pokeType := range pokemon.Types {
		pokeTypes = append(pokeTypes, pokeType.Type.Name)
	}
	pokemonStats := PokemonStats{Hp: pokemon.Stats[0].BaseStat, Attack: pokemon.Stats[1].BaseStat, Defense: pokemon.Stats[2].BaseStat, SpecialAttack: pokemon.Stats[3].BaseStat, SpecialDefense: pokemon.Stats[4].BaseStat, Speed: pokemon.Stats[5].BaseStat}
	pokemonDescription := PokemonDescription{Height: pokemon.Height, Weight: pokemon.Weight, PokemonStats: pokemonStats, Types: pokeTypes}
	return pokemonDescription, nil
}
