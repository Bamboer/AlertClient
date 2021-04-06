package client

import "grafana/pkg/configer"

var (
        DashboardPath = "/api/dashboards/uid/"
        AlertPath     = "/api/alerts/"
        OrgPath       = "/api/org"
)

type Alert struct {
        Id             int       `json:"id"`
        DashboardId    int       `json:"dashboardId"`
        DashboardUid   string    `json:"dashboardUid"`
        DashboardSlug  string    `json:"dashboardSlug"`
        PanelId        int       `json:"panelId"`
        Name           string    `json:"name"`
        State          string    `json:"state"`
        NewStateDate   string    `json:"newStateDate"`
        EvalDate       string    `json:"evalDate"`
        EvalData       EvalDatas `json:"evalData,omitempty"`
        ExecutionError string    `json:"executionError,omitempty"`
        Url            string    `json:"url"`
}

type EvalDatas struct {
        EvalMatches []Evalmatche `json:"evalmatches,omitempty"`
}

type Evalmatche struct {
        Metric string            `json:"metric,omitempty"`
        Tags   map[string]string `json:"tags,omitempty"`
        Value  float32           `json:"value,omitempty"`
}

type AlertInfo struct {
        Id             int         `json:"Id"`
        Version        int         `json:"Version"`
        OrgId          int         `json:"OrgId"`
        DashboardId    int         `json:"DashboardId"`
        PanelId        int         `json:"PanelId"`
        Name           string      `json:"Name"`
        Message        string      `json:"Message,omitempty"`
        Serverity      string      `json:"Serverity,omitempty"`
        State          string      `json:"State"`
        Handler        int         `json:"Handler"`
        Silenced       bool        `json:"Silenced"`
        ExecutionError string      `json:"ExecutionError,emitempty"`
        Frequency      int         `json:"Frequency"`
        For            int         `json:"For"`
        EvalData       EvalDatas   `json:"EvalData"`
        NewStateDate   string      `json:"NewStateDate"`
        StateChanges   int         `json:"StateChanges"`
        Created        string      `json:"Created"`
        Updated        string      `json:"Updated"`
        Settings       interface{} `json:"Settings"`
}

type SimpleInfo struct {
        Name         string
        AlertMetrics string
        AlertValues  []float32
        PanelId      int
        OrgId        int
        DbUid        string
        DbSlug       string
        Frequency    int
        AlertNum  *int
        TempVar      map[string]string
        RenderURL    string
}

func GetAlerts() ([]Alert, error) {
        alerts := []Alert{}
        grafana_conf := configer.ConfigParse()
        //        info.Println(grafana_conf)
        C, _ := NewGrafanaClient(grafana_conf.Grafana_uri, grafana_conf.Grafana_token)
        if err := C.Get(AlertPath, &alerts); err != nil {
                info.Println(err)
                return alerts, err
        }
        return alerts, nil
}

func GetAlert(AlertId string) (*AlertInfo, error) {
        alerts := &AlertInfo{}
        grafana_conf := configer.ConfigParse()
        //        info.Println(grafana_conf)
        C, _ := NewGrafanaClient(grafana_conf.Grafana_uri, grafana_conf.Grafana_token)
        if err := C.Get(AlertPath+AlertId, &alerts); err != nil {
                info.Println(err)
                return alerts, err
        }
        return alerts, nil
}
