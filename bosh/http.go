package bosh

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (t Target) UA() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: t.Insecure,
			},
		},
	}
}

func (t Target) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", t.URL+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(t.Username, t.Password))
	return t.UA().Do(req)
}

func (t Target) InterpretJSON(res *http.Response, v interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

func (t Target) InterpretJSONList(res *http.Response) ([]string, error) {
	/* FIXME get a line reader, and parse stuff into a list of JSON-ish strings */
	return nil, nil
}
