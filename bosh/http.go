package bosh

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func (t Target) UA() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Host = via[0].URL.Host
			req.Header.Set("Authorization", via[0].Header.Get("Authorization"))
			return nil
		},
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
	req.Header.Set("Authorization", "Basic "+basicAuth(t.Username, t.Password))
	return t.UA().Do(req)
}

func (t Target) Post(url string, payload interface{}) (*http.Response, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", t.URL+url, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(t.Username, t.Password))
	return t.UA().Do(req)
}

func (t Target) PostYAML(url string, raw []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", t.URL+url, strings.NewReader(string(raw)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "text/yaml")
	req.Header.Set("Authorization", "Basic "+basicAuth(t.Username, t.Password))
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

func (t Target) InterpretJSONList(res *http.Response) ([][]byte, error) {
	l := make([][]byte, 0)

	sc := bufio.NewScanner(res.Body)
	for sc.Scan() {
		l = append(l, sc.Bytes())
	}
	return l, sc.Err()
}
