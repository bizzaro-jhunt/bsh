package bosh

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (t Target) Get(uri string) (*http.Response, error) {
	req, err := http.NewRequest("GET", t.URL+uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+basicAuth(t.Username, t.Password))
	return t.UA().Do(req)
}

func (t Target) GetJSON(uri string, v interface{}) error {
	r, err := t.Get(uri)
	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		return fmt.Errorf("BOSH API returned %s", r.Status)
	}

	err = t.InterpretJSON(r, v)
	if err != nil {
		return err
	}

	return nil
}

func (t Target) Delete(uri string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", t.URL+uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+basicAuth(t.Username, t.Password))

	res, err := t.UA().Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 302 {
		u, err := url.Parse(res.Header.Get("Location"))
		if err != nil {
			return nil, err
		}
		return t.Get(u.Path) /* bosh never redirs to querystrings... */
	}
	return res, err
}

func (t Target) Post(uri string, payload interface{}) (*http.Response, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", t.URL+uri, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(t.Username, t.Password))
	return t.UA().Do(req)
}

func (t Target) PostYAML(uri string, raw []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", t.URL+uri, strings.NewReader(string(raw)))
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
