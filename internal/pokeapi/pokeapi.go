package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const base_url = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
}

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) GetLocationAreas(url *string) (LocationAreasResponse, error) {
	endpoint := "/location-area"
	full_url := base_url + endpoint
	if url != nil {
		full_url = *url
	}

	req, err := http.NewRequest("GET", full_url, nil)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return LocationAreasResponse{}, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	resp_2 := LocationAreasResponse{}
	err = json.Unmarshal(data, &resp_2)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	return resp_2, nil
}

