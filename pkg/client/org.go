package client
import (
  "grafana/pkg/configer"
)

type Org struct {
	Id      int     `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Address Address `json:"address,omitempty"`
}

type Address struct {
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	City     string `json:"city,omitempty"`
	ZipCode  string `json:"zipCode,omitempty"`
	State    string `json:"state,omitempty"`
	Country  string `json:"country,omitempty"`
}

func GetOrg() (*Org, error) {
	org := &Org{}
	grafana_conf := configer.ConfigParse()
	C, _ := NewGrafanaClient(grafana_conf.Grafana_uri, grafana_conf.Grafana_token)
	if err := C.Get(OrgPath, org); err != nil {
		info.Println(err)
		return org, err
	}
	return org, nil
}
