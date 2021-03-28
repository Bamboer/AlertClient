package client

import (
        "io"
	"log"
	"os"
	"path"
	"strings"
	"grafana/pkg/configer"
)

var (
	info          *log.Logger
	DashboardPath = "/api/dashboards/uid/"
	AlertPath     = "/api/alerts/"
	OrgPath       = "/api/org"
)

func init() {
	arg := path.Base(os.Args[0])
	logfile := strings.ToLower(arg + "1.log")
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open file: ", err)
	}
	info = log.New(io.MultiWriter(os.Stdout, file), "Info ", log.Ldate|log.Ltime|log.Lshortfile)
}

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

func GetAlerts() ([]Alert, error) {
	alerts := []Alert{}
	grafana_conf := configer.Configfile()
        info.Println(grafana_conf)
	C, _ := NewGrafanaClient(grafana_conf.Grafana_uri, grafana_conf.Grafana_token)
	if err := C.Get(AlertPath, &alerts); err != nil {
		info.Println(err)
		return alerts, err
	}
        log.Println("data: ",alerts)
	return alerts, nil
}
