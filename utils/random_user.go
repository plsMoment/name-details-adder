package utils

import (
	"encoding/json"
	"errors"
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

func GenerateUser() (fn FullName, e error) {
	scope := "utils.random_user.GenerateUser"
	resp, err := http.Get("https://randomuser.me/api/?inc=name") //https://api.randomdatatools.ru/?params=LastName,FirstName
	if err != nil {
		e = fmt.Errorf("%s: %w", scope, err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			e = errors.Join(e, fmt.Errorf("%s: %w", scope, err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		e = fmt.Errorf("%s: StatusCode is not %d", scope, http.StatusOK)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e = fmt.Errorf("%s: %w", scope, err)
		return
	}

	res := UserGenerated{}
	if err = json.Unmarshal(body, &res); err != nil {
		e = fmt.Errorf("%s: %w", scope, err)
		return
	}

	if len(res.Results) == 0 {
		e = fmt.Errorf("%s: response with user name is empty", scope)
		return
	}

	fn = res.Results[0].Fn
	return
}
