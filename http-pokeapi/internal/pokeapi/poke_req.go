// Package pokeapi provides a client for interacting with the Pokémon API (https://pokeapi.co/).
package pokeapi

import (
	"encoding/json"
	"fmt"
	models "http-pokeapi/internal/models"
	"io"
	"log"
	"net/http"
	"os"
)

const baseURL = "https://pokeapi.co/api/v2"

type Client struct {
	http.Client
}

func (client *Client) Pokemon(name string) (models.Pokemon, error) {
	endURL := "/pokemon/"
	fullURL := baseURL + endURL + name

	req, err := http.NewRequest("GET", fullURL, nil)

	if err != nil {
		return models.Pokemon{}, err
	}

	res, err := client.Do(req)

	if err != nil {
		return models.Pokemon{}, err
	}

	if res.StatusCode > 399 {
		return models.Pokemon{}, fmt.Errorf("bad status code : %v", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Pokemon{}, err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Printf("%s", err)
			os.Exit(1)
		}
	}()

	pokeinfo := models.Pokemon{}
	err = json.Unmarshal(data, &pokeinfo)
	if err != nil {
		return models.Pokemon{}, err
	}
	return pokeinfo, nil

}

func (client *Client) LocationArearesponse() (models.Location, error) {
	endURL := "/location-area"
	fullURL := baseURL + endURL

	req, err := http.NewRequest("GET", fullURL, nil)

	if err != nil {
		return models.Location{}, err
	}

	res, err := client.Do(req)

	if err != nil {
		return models.Location{}, err
	}

	if res.StatusCode > 399 {
		return models.Location{}, fmt.Errorf("bad status code : %v", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Location{}, err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Printf("%s", err)
			os.Exit(1)
		}
	}()

	LocationAreaValues := models.Location{}
	err = json.Unmarshal(data, &LocationAreaValues)
	if err != nil {
		return models.Location{}, err
	}
	return LocationAreaValues, nil

}

func (client *Client) Pokelocationres(arg string) (models.Pokelocation, error) {
	endURL := "/location-area/"
	fullURL := baseURL + endURL + arg

	req, err := http.NewRequest("GET", fullURL, nil)

	if err != nil {
		return models.Pokelocation{}, err
	}

	res, err := client.Do(req)

	if err != nil {
		return models.Pokelocation{}, err
	}

	if res.StatusCode > 399 {
		return models.Pokelocation{}, fmt.Errorf("bad status code : %v", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Pokelocation{}, err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Printf("%s", err)
			os.Exit(1)
		}
	}()

	LocationAreaValues := models.Pokelocation{}
	err = json.Unmarshal(data, &LocationAreaValues)
	if err != nil {
		return models.Pokelocation{}, err
	}
	return LocationAreaValues, nil

}
