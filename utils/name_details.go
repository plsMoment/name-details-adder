package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Details struct {
	Age         int
	Gender      string
	Nationality string
}

func NameDetails(name string) (*Details, error) {
	age, err := agify(name)
	if err != nil {
		return nil, err
	}
	gender, err := genderize(name)
	if err != nil {
		return nil, err
	}
	nationality, err := nationalize(name)
	if err != nil {
		return nil, err
	}
	return &Details{Age: age, Gender: gender, Nationality: nationality}, nil
}

type Agified struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func agify(name string) (int, error) {
	scope := "utils.name_details.agify"
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", scope, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("%s: StatusCode is not %d", scope, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", scope, err)
	}

	res := Agified{}
	if err = json.Unmarshal(body, &res); err != nil {
		return -1, fmt.Errorf("%s: %w", scope, err)
	}

	if res.Age == 0 {
		return 0, fmt.Errorf("%s: name %s haven't age", scope, name)
	}
	return res.Age, nil
}

type Genderized struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

func genderize(name string) (string, error) {
	scope := "utils.name_details.genderize"

	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "", fmt.Errorf("%s: %w", scope, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: StatusCode is not %d", scope, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", scope, err)
	}

	res := Genderized{}

	if err = json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("%s: %w", scope, err)
	}
	if res.Gender == "" {
		return "", fmt.Errorf("%s: name %s haven't gender", scope, name)
	}
	return res.Gender, nil
}

type Nationalized struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func nationalize(name string) (string, error) {
	scope := "utils.name_details.nationalize"
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", fmt.Errorf("%s: %w", scope, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: StatusCode is not %d", scope, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", scope, err)
	}

	res := Nationalized{}
	if err = json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("%s: %w", scope, err)
	}
	if len(res.Country) == 0 {
		return "", fmt.Errorf("%s: name %s haven't nationality", scope, name)
	}
	return res.Country[0].CountryID, nil
}
