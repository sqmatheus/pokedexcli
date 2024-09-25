package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sqmatheus/pokedexcli/internal/cache"
)

const (
	BaseUrl         string = "https://pokeapi.co/api/v2"
	LocationUrl     string = BaseUrl + "/location"
	LocationAreaUrl string = BaseUrl + "/location-area"
	PokemonUrl      string = BaseUrl + "/pokemon"
)

var apiCache = cache.NewCache(time.Second * 30)

type Pagination[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

type LocationEntry struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Location struct {
	Areas []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"areas"`
}

type LocationArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func newClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

func Get[T any](url string) (T, error) {
	var result T

	if data, ok := apiCache.Get(url); ok {
		if err := json.Unmarshal(data, &result); err != nil {
			return result, err
		}
		return result, nil
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}
	client := newClient()

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return result, errors.New("not found")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	// var prettyJSON bytes.Buffer
	// if err = json.Indent(&prettyJSON, data, "", "\t"); err != nil {
	// 	return result, err
	// }
	// fmt.Println(prettyJSON.String())

	if err := json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	apiCache.Add(url, data)

	return result, nil
}

func GetLocations(url string) (Pagination[LocationEntry], error) {
	return Get[Pagination[LocationEntry]](url)
}

func GetLocation(location string) (Location, error) {
	return Get[Location](fmt.Sprintf("%s/%s", LocationUrl, location))
}

func GetLocationArea(area string) (LocationArea, error) {
	return Get[LocationArea](fmt.Sprintf("%s/%s", LocationAreaUrl, area))
}

func GetPokemon(pokemon string) (Pokemon, error) {
	return Get[Pokemon](fmt.Sprintf("%s/%s", PokemonUrl, pokemon))
}
