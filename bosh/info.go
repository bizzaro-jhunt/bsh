package bosh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Feature struct {
	Status bool                   `json:"status"`
	Extras map[string]interface{} `json:"extras"`
}

type UserAuth struct {
	Type    string                 `json:"type"`
	Options map[string]interface{} `json:"options"`
}

type Info struct {
	Name     string             `json:"name"`
	UUID     string             `json:"uuid"`
	Version  string             `json:"version"`
	User     string             `json:"user"`
	CPI      string             `json:"cpi"`
	Features map[string]Feature `json:"features"`
	UserAuth UserAuth           `json:"user_authentication"`
}

func (t Target) GetInfo() (Info, error) {
	var info Info

	req, err := http.NewRequest("GET", t.URL+"/info", nil)
	if err != nil {
		return info, err
	}

	r, err := t.UA().Do(req)
	if err != nil {
		return info, err
	}

	if r.StatusCode != 200 {
		return info, fmt.Errorf("BOSH API returned %s", r.Status)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal(b, &info)
	if err != nil {
		return info, err
	}

	return info, nil
}
