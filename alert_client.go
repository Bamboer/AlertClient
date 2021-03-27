package main
import (
  "io"
  "strings"
  "os"
  "fmt"
  "flag"
  "path"
//  "time"
  "encoding/json"
//  "path/filepath"
  "grafana/pkg/configer"
  "grafana/pkg/client"
  "log"
//  "grafana/pkg/notification"
)

var (
  alerts    []alert 
  org       Org
  Dashboard dashboard
)
var (
  info  *log.Logger
  Version = "v1"
  DEFALTCONF = "alert_client.conf"
  DashboardURL = "/api/dashboards/uid/"
  AlertPath = "/api/alerts/"
  OrgId = "/api/org"
  ConfigFile = flag.String("config","alert_client.conf","alert client configuration file.")
  version = flag.Bool("version",false,"Print version.")
)

type obj struct {
  dingding string
  grafana_token string
  grafana_uri string
  notifications string
  smtpServer *smtpServer
  alert_log string
  client_log string
}

type smtpServer struct {
  username string
  password string
  smtpAddress string
}

type Org struct{
  Id int `json:"id,omitempty"`
  Name string `json:"name,omitempty"`
  Address address  `json:"address,omitempty"`
}

type address struct{
  address1 string `json:"address1,omitempty"`
  address2 string `json:"address2,omitempty"`
  city     string `json:"city,omitempty"`
  zipCode  string `json:"zipcode,omitempty"`
  state    string `json:"state,omitempty"`
  country  string `json:"country,omitempty"`
}


type alert struct{
  id             int       `json:"id"`
  dashboardId    string    `json:"dashboardid"`
  dashboardSlug  string    `json:"dashboardslug"`
  panelId        int       `json:"panelId"`
  name           string    `json:"name"`
  state          string    `json:"state"`
  newStateDate   string    `json:"newstatedate"`
  evalDate       string    `json:"evaldate"`
  evalData       evalDatas `json:"evaldata,omitempty"`
  executionError string    `json:"executionerror,omitempty"`
  url            string    `json:"url"`
}

type evalDatas struct{
  evalMatches  []evalmatche `json:"evalmatches,omitempty"`
}

type evalmatche struct{
  metric       string   `json:"metric,omitempty"`
  tags         map[string] string`json:"tags,omitempty"`
  value        string   `json:"value,omitempty"`
}

type dashboard struct{
   meta        metainfo        `json:"meta"`
   dashboard  dashboardinfo    `json:"dashboard"`
}

type metainfo struct{
   Type       string     `json:"type"`
   canSave    bool       `json:"cansave"`
   canEdit    bool       `json:"canedit"`
   canAdmin   bool       `json:"canadmin"`
   slug       string     `json:"slug"`
   url        string     `json:"url"`
   expires    string     `json:"expires"`
   created    string     `json:"created"`
   updated    string     `json:"updated"`
   updatedBy  string     `json:"updatedby"`
   createdBy  string     `json:"createdby"`
   version    string     `json:"version"`
   hasAcl     bool       `json:"hasacl"`
   isFolder   bool       `json:"isfolder"`
   folderId   int        `json:"folderid"`
   folderUrl  string     `json:"folderurl,omitempty"`
   provisioned bool      `json:"provisioned"`
   provisionedExternalId  string     `json:"provisionedexternalid"`
}

type dashboardinfo   struct{
   annotations    interface{} `json:"annotations"`
   editable       bool        `json:"editable"`
   gnetId        string       `json:"gnetid"`
   graphTooltip  int          `json:"graphtooltip"`
   id            int          `json:"id"`
   links         []string     `json:"links,omitempty"`
   panels        interface{}  `json:"panels"`
   schemaVersion int          `json:"schemaversion"`
   style         string       `json:"style"`
   tags          []string     `json:"tags,omitempty"`
   templating    map[string] []interface{}  `json:"templating,omitempty"`
   time          times        `json:"time"`
   timepicker    interface{}  `json:"timepicker,omitempty"`
   timezone      string       `json:"timezone,omitempty"`
   uid           string       `json:"uid"`
   variables     variable     `json:"variables,omitempty"`
   version       int          `json:"version"`
}


type times struct{
   from  string  `json:"from,omitempty"`
   to    string  `json:"to,omitempty"`
}

type variable struct{
   list   []interface{}   `json:"list,omitempty"`
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
  if len(os.Args) >2 {
    if os.Args[1] == "version"{
    fmt.Println("grafana alert client ",Version)
    os.Exit(0)
    }
  }
  grafana_conf := configfile()
  alert_client,_ := client.NewGrafanaClient(grafana_conf.grafana_uri,grafana_conf.grafana_token)
  GetData,_ := alert_client.Get("/api/alerts")
  err := json.NewDecoder(GetData).Decode(&alerts)
  if err != nil{
    info.Println(err)
  }
  info.Println(alerts)
}
func init(){
  arg := path.Base(os.Args[0])
  logfile := strings.ToLower(arg + "1.log")
  file,err := os.OpenFile(logfile,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
  if err != nil{
     info.Println("Failed to open file: ",err)
  }
  info = log.New(io.MultiWriter(os.Stdout,file),"Info",log.Ldate|log.Ltime|log.Lshortfile)
}

func configfile() *obj{
  configuration := &obj{}
  fmt.Println("configfile: ",configuration)
  conf,_ := configer.NewConfig(DEFALTCONF)
  if *ConfigFile != "alert_client.conf"{
     conf,_ = configer.NewConfig(*ConfigFile)
  }
  configuration.dingding,_ = conf.GetString("dingding")
  configuration.grafana_token,_ = conf.GetString("grafana_token")
  configuration.notifications,_ = conf.GetString("notifications")
  configuration.grafana_uri,_ = conf.GetString("grafana_uri")
  configuration.smtpServer.username,_  = conf.GetString("username")
  configuration.smtpServer.password,_ = conf.GetString("password")
  configuration.smtpServer.smtpAddress,_ = conf.GetString("smtpAddress")
  return configuration 
}
