package client

import (
	"net/http"
	"net/url"
	"encoding/json"
)

type grafana_client struct {
	uri        *url.URL
	token      string
	httpClient *http.Client
}

func NewGrafanaClient(uri, token string) (*grafana_client, error) {
	url, err := url.Parse(uri)
	token = "Bearer " + token
//        info.Println(url)
//        info.Println(token)
	if err != nil {
		info.Println(err)
		return nil, err
	}
	return &grafana_client{
		uri:        url,
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

func (c *grafana_client) Get(path string, v interface{}) error {
	uri := c.uri
	uri.Path = path
	req, err := http.NewRequest("GET", uri.String(), nil)
        if err != nil {
		info.Println(err)
		return err
	}
	req.Header.Add("Authorization", c.token)
	//        req.SetBasicAuth("admin", "admin")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
        if err != nil {
		info.Println(err)
		return err
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
		info.Println(err)
		return err
	}
//        log.Println(v)
	return nil
}
