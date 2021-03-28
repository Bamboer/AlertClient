package configer
import (
  "flag"
)


var (
  DEFALTCONF = "alert_client.conf"
  ConfigFile = flag.String("config","alert_client.conf","alert client configuration set.")
)
type Obj struct {
  Dingding string
  Grafana_token string
  Grafana_uri string
  Notifications string
  SmtpServer SmtpInfo
  Alert_log string
  Client_log string
}

type SmtpInfo struct {
  Username string
  Password string
  SmtpAddress string
}

func Configfile() *Obj{
  configuration := &Obj{}
  conf,_ := NewConfig(DEFALTCONF)
  if *ConfigFile != "alert_client.conf"{
     conf,_ = NewConfig(*ConfigFile)
  }
  configuration.Dingding,_ = conf.GetString("dingding")
  configuration.Grafana_token,_ = conf.GetString("grafana_token")
  configuration.Notifications,_ = conf.GetString("notifications")
  configuration.Grafana_uri,_ = conf.GetString("grafana_uri")
  configuration.SmtpServer.Username,_  = conf.GetString("username")
  configuration.SmtpServer.Password,_ = conf.GetString("password")
  configuration.SmtpServer.SmtpAddress,_ = conf.GetString("smtpAddress")
  return configuration 
}
