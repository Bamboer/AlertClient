package main


type obj struct {
  dingding string
  grafana_token string
  grafana_uri string
  nitifications []string
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

var Alerts []alert `json:"alerts"`

type alert stuct{
  id             int       `json:"id"`
  dashboardId    string    `json:"dashboardid"`
  dashboardSlug  string    `json:"dashboardslug"`
  panelId        int       `json:"panelId"`
  name           string    `json:"name"`
  state          string    `json:"state"`
  newStateDate   string    `json:"newstatedate"`
  evalDate       string    `json:"id"`
  evalData       evalDatas `json:"evaldata,omitempty"`
  executionError string    `json:"executionerror,omitempty"`
  url            string    `json:"id"`
}

type evalDatas struct{
  evalMatches  []evalmatche `json:"evalmatches,omitempty"`
}

type evalmatche struct{
  metric       string   `json:"metric,omitempty"`
  tags         tag      `json:"tags,omitempty"`
  value        string   `json:"value,omitempty"`
}

type dashboard struct{
   meta        metainfo        `json:"meta"`
   dashboard  dashboardinfo    `json:"dashboard"`
}

type metainfo struct{
   type       string     `json:"type"`
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
   provisioned boll      `json:"provisioned"`
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
   templating    list         `json:"templating,omitempty"`
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


func configfile()*obj{
  configuration := *obj{}
  if ConfigFile = nil{
     conf := configer.NewConfig(DEFALTCONF)
  }else{
     conf := configer.NewConfig(*ConfigFile)
  }
  configuration.dingding = conf.GetString("dingding")
  configuration.grafana_token = conf.GetString("grafana_token")
  configuration.grafana_url = conf.GetString("grafana_url")
  configuration.notifications = conf.GetString("notifications")
  configuration.grafana_url = conf.GetString("grafana_url")
  configuration.smtpServer.username  = conf.GetString("username")
  configuration.smtpServer.password = conf.GetString("password")
  configuration.smtpServer.smtpAddress = conf.GetString("smtpAddress")
  return &configuration 
}
