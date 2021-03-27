package main
import (
  "os"
  "fmt"
  "time"
  "encoding/json"
  "path/filepath"
  "grafana/pkg/configer"
  "grafana/pkg/client"
  "grafana/pkg/log"
  "grafana/pkg/notification"
)

var (
  ConfigFile string
  Version   bool
  alerts    Alerts
  org      Org
  Dashboard dashboard
)
var (
  Version = "v1"
  DEFALTCONF = "alert_client.conf"
  DashboardURL = "/api/dashboards/uid/"
  AlertPath = "/api/alerts/"
  OrgId = "/api/org"
)

func init(){
  ConfigFile := flag.String("config","alert_client.conf","alert client configuration file.")
  Version := flag.Bool("version",false,"Print version.")

  if os.Args[1] == "version"{
    fmt.Println("grafana alert client ",Version)
    os.Exit(0)
  }
}
/*
func run(alert_client,){
  ticker := time.NewTicker(30 * time.Second)
  for _ = range ticker.C{
    data : = alert_client.Get()
  }
}*/

func main(){
  flag.Parse()
  grafana_conf := configfile()
  alert_client,_ := client.NewGrafanaClient(grafana_conf.grafana_url,grafana_conf.grafana_token)
  GetData := alert_client.Get("/api/alerts")
  err := json.NewDecoder(GetData).Decode(&alerts)
  if err != nil{
    log.Infoln(err)
  }
  log.Infoln(alerts)
}
