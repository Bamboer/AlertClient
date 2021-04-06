package notification

import (
        "gopkg.in/ini.v1"
        "grafana/pkg/client"
        "grafana/pkg/configer"
        "io"
        "log"
        "os"
        "path"
        "strings"
)

var (
        info   *log.Logger
        SNS    = make(map[string]func(state string, msg client.SimpleInfo, b []byte)error)
)

func init() {
        arg := path.Base(os.Args[0])
        logfile := strings.ToLower(arg + ".log")
        file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
                log.Println("Failed to open file: ", err)
        }
        cfg, err := ini.Load(*configer.ConfigFile)
        if err != nil {
                info.Println("Fail to read file: ", err)
                os.Exit(1)
        }
        mode := cfg.Section("").Key("mode").In("dev", []string{"dev", "debug", "prd"})
        if mode == "dev" || mode == "debug" {
                info = log.New(io.MultiWriter(os.Stdout, file), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
        } else if mode == "prd" {
                info = log.New(file, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
        }
}

type Notification interface {
        //send message to the receiver
        Send(state string) error
}

func Emit(state string, msg client.SimpleInfo, b []byte) {
        for k, send := range SNS {
                if err := send(state, msg, b); err != nil {
                        info.Println(k, "send err: ", err)
                } else {
                        info.Println(k, "send message: ")
                }
        }
}
