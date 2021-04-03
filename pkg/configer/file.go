package configer
import (
  "flag"
)


var (
  err   error
  DEFALTCONF = "alert_client.conf"
  ConfigFile = flag.String("config","alert_client.conf","alert client configuration set.")
)
type Obj struct {
  Dingding string
  Grafana_token string
  Grafana_uri string
  NotificationsLeader string
  
  Notifications string
  Notifications_cc string
  Notifications_bcc string
  SmtpServer SmtpInfo
  Alert_log string
  Client_log string
}

type SmtpInfo struct {
  Username string
  Password string
  SmtpAddress string
  Port      int
}

func ConfigParse() *Obj{
  configuration := &Obj{}
  conf,_ := NewConfig(DEFALTCONF)
  //info.Println(conf)
  if *ConfigFile != "alert_client.conf"{
     conf,_ = NewConfig(*ConfigFile)
  }

  configuration.Dingding,err = conf.GetString("dingding")
  if err != nil{
     info.Println(err)
  }
  configuration.Grafana_token,_ = conf.GetString("grafana_token")
  configuration.Notifications,_ = conf.GetString("notifications")
  configuration.Notifications_cc,_ = conf.GetString("notifications_cc")
  configuration.Notifications_bcc,_ = conf.GetString("notifications_bcc")
  configuration.Grafana_uri,_ = conf.GetString("grafana_uri")
  configuration.SmtpServer.Username,_  = conf.GetString("username")
  configuration.SmtpServer.Password,_ = conf.GetString("password")
  configuration.SmtpServer.SmtpAddress,_ = conf.GetString("smtpAddress")
  configuration.SmtpServer.Port,_ = conf.GetInt("smtpPort")
  configuration.Alert_log,_ = conf.GetString("alert_log")
  configuration.Client_log,_ = conf.GetString("client_log")
  return configuration
}
