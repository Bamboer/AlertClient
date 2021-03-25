package main

type obj struct {
  dingding string
  grafana_token string
  grafana_uri string
  nitification_list []string
  smtpServer *smtpServer
  alert_log string
  client_log string
}

type smtpServer struct {
  username string
  password string
  smtpAddress string
  alert_time  [2]int8
}

func conf()*configer.Config{
  if ConfigFile = nil{
     conf := configer.NewConfig(DEFALTCONF)
  }else{
     conf := configer.NewConfig(*ConfigFile)
  }
  return conf
}
