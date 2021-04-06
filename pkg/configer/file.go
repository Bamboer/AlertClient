package configer

import (
        "gopkg.in/ini.v1"
        "io"
        "log"
        "os"
        "flag"
        "path"
//        "strconv"
        "strings"
)

var (
        info       *log.Logger
        ConfigFile = flag.String("config","alert_client.conf","set default configuration for this app.")
)

type Obj struct {
        Mode                string
        Dingding            string
        Grafana_token       string
        Grafana_uri         string
        NotificationsLeader string

        Notifications     string
        Notifications_cc  string
        Notifications_bcc string
        Notifications_dau string

        SmtpServer SmtpInfo

        Alert_log  string
        Client_log string
}

type SmtpInfo struct {
        Username    string
        Password    string
        SmtpAddress string
        Port        string
        StartTime   string
        EndTime     string
}

func init() {
        arg := path.Base(os.Args[0])
        logfile := strings.ToLower(arg + ".log")
        file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
                log.Println("Failed to open file: ", err)
        }
        cfg, err := ini.Load(*ConfigFile)
        if err != nil {
                info.Println("Fail to read file: ", err)
                os.Exit(1)
        }
        mode := cfg.Section("").Key("mode").In("dev", []string{"dev", "debug", "prd"})
        if mode == "dev" || mode == "debug" {
                info = log.New(io.MultiWriter(os.Stdout, file), "", log.Ldate|log.Ltime|log.Lshortfile)
        } else if mode == "prd" {
                info = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
        }
}

func ConfigParse() *Obj {
        configuration := &Obj{}
        cfg, err := ini.Load(*ConfigFile)
        if err != nil {
                info.Println("Fail to read file: ", err)
                os.Exit(1)
        }
        configuration.Mode = cfg.Section("").Key("mode").In("dev", []string{"dev", "debug", "prd"})
        configuration.Dingding = cfg.Section("").Key("dingding").String()
        configuration.Grafana_token = cfg.Section("grafana").Key("grafana_token").String()
        configuration.Grafana_uri = cfg.Section("grafana").Key("grafana_uri").String()

        configuration.Notifications = cfg.Section("email").Key("notifications").String()
        configuration.Notifications_cc = cfg.Section("email").Key("notifications_cc").String()
        configuration.Notifications_bcc = cfg.Section("email").Key("notifications_bcc").String()
        configuration.Notifications_dau = cfg.Section("email").Key("notifications_dau").String()

        configuration.SmtpServer.Username = cfg.Section("smtp_server").Key("username").String()
        configuration.SmtpServer.Password = cfg.Section("smtp_server").Key("password").String()
        configuration.SmtpServer.SmtpAddress = cfg.Section("smtp_server").Key("smtpAddress").String()
        configuration.SmtpServer.Port = cfg.Section("smtp_server").Key("smtpPort").String()
        configuration.SmtpServer.StartTime = cfg.Section("smtp_server").Key("start_time").String()
        configuration.SmtpServer.EndTime = cfg.Section("smtp_server").Key("end_time").String()

        configuration.Alert_log = cfg.Section("log").Key("alert_log").String()
        configuration.Client_log = cfg.Section("log").Key("client_log").String()
//        info.Println("configuration: ",configuration)
        return configuration
}
