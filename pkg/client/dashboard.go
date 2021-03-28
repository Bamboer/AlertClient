package client
import (
   "grafana/pkg/configer"
)

type Dashboard struct {
	Meta      Metainfo      `json:"meta"`
	Dashboard Dashboardinfo `json:"dashboard"`
}

type Metainfo struct {
	Type                  string `json:"type"`
	CanSave               bool   `json:"canSave"`
	CanEdit               bool   `json:"canEdit"`
	CanAdmin              bool   `json:"canAdmin"`
	CanStar               bool   `json:"canStar"`
	Slug                  string `json:"slug"`
	Url                   string `json:"url"`
	Expires               string `json:"expires"`
	Created               string `json:"created"`
	Updated               string `json:"updated"`
	UpdatedBy             string `json:"updatedBy"`
	CreatedBy             string `json:"createdBy"`
	Version               int    `json:"version"`
	HasAcl                bool   `json:"hasAcl"`
	IsFolder              bool   `json:"isFolder"`
	FolderId              int    `json:"folderId"`
	FolderUrl             string `json:"folderUrl,omitempty"`
	Provisioned           bool   `json:"provisioned"`
	ProvisionedExternalId string `json:"provisionedExternalId"`
}

type Dashboardinfo struct {
	Annotations   interface{}              `json:"annotations"`
	Editable      bool                     `json:"editable"`
	GnetId        string                   `json:"gnetId"`
	GraphTooltip  int                      `json:"graphTooltip"`
	Id            int                      `json:"id"`
	Links         []string                 `json:"links,omitempty"`
	Panels        interface{}              `json:"panels"`
	SchemaVersion int                      `json:"schemaVersion"`
	Style         string                   `json:"style"`
	Tags          []string                 `json:"tags,omitempty"`
	Templating    map[string][]interface{} `json:"templating,omitempty"`
	Time          Times                    `json:"time"`
	Timepicker    interface{}              `json:"timepicker,omitempty"`
	Timezone      string                   `json:"timezone,omitempty"`
	title         string                   `json:"title,omitempty"`
	Uid           string                   `json:"uid"`
	Variables     Variable                 `json:"variables,omitempty"`
	Version       int                      `json:"version"`
}

type Times struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

type Variable struct {
	List []interface{} `json:"list,omitempty"`
}

func GetDashboard(DashboardUid string) (*Dashboard, error) {
	db := &Dashboard{}
	grafana_conf := configer.ConfigParse()
	C, _ := NewGrafanaClient(grafana_conf.Grafana_uri, grafana_conf.Grafana_token)
	if err := C.Get(DashboardPath + DashboardUid, db); err != nil {
		info.Println(err)
		return db, err
	}
	return db, nil
}
