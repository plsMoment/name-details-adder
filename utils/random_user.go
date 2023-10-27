package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FullName struct {
	Name       string     `json:"first"`
	Surname    string     `json:"last"`
	Patronymic NullString `json:"patronymic,omitempty"`
}

type UserGenerated struct {
	Results []struct {
		Fn FullName `json:"name"`
	} `json:"results"`
}

func GenerateUser() (FullName, error) {
	scope := "utils.random_user.GenerateUser"
	resp, err := http.Get("https://randomuser.me/api/?inc=name") //https://api.randomdatatools.ru/?params=LastName,FirstName
	if err != nil {
		return FullName{}, fmt.Errorf("%s: %w", scope, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FullName{}, fmt.Errorf("%s: StatusCode is not %d", scope, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FullName{}, fmt.Errorf("%s: %w", scope, err)
	}

	res := UserGenerated{}
	if err = json.Unmarshal(body, &res); err != nil {
		return FullName{}, fmt.Errorf("%s: %w", scope, err)
	}

	if len(res.Results) == 0 {
		return FullName{}, fmt.Errorf("%s: response with user name is empty", scope)
	}
	return res.Results[0].Fn, nil
}
