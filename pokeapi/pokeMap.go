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

// const locationsPath = "pokeLocations.json"
const areasPath = "pokeapi/pokeAreas.json"

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

func cachePokeData[PDT PokeDataType](pokeData []PDT, filePath string) error {
	var cachedPokeData []PDT
	_, err := os.Stat(filePath)
	if errors.Is(err, fs.ErrExist) {
		pokeDataFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		decoder := json.NewDecoder(pokeDataFile)
		if err := decoder.Decode(&cachedPokeData); err != nil {
			pokeDataFile.Close()
			return err
		}
		pokeDataFile.Close()
	}

	cachedPokeData = append(cachedPokeData, pokeData...)

	if errors.Is(err, fs.ErrNotExist) {
		dirSep := strings.LastIndex(filePath, "/")
		os.MkdirAll(filePath[:dirSep], 0755)
	}

	pokeDataFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer pokeDataFile.Close()
	encoder := json.NewEncoder(pokeDataFile)
	if err := encoder.Encode(cachedPokeData); err != nil {
		return err
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
		if err := cachePokeData(missingPokeData, filePath); err != nil {
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
	err = cachePokeData(pokeData, filePath)
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

func GetMapStrings(minIndex, maxIndex int) ([]string, error) {
	pathExtension := strings.Split(areasPath, ".")
	currentAreaPath := pathExtension[0] + strconv.Itoa(minIndex/20) + "." + pathExtension[1]
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
