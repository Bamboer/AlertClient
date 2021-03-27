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
