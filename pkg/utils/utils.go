package utils

import (
        "context"
        "gopkg.in/ini.v1"
        "io"
        "log"
        "os"
        "path"
        "strings"
        "grafana/pkg/configer"
)

var (
        info  *log.Logger
        UTILS = []func(ctx context.Context){}
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

func UtilRuner(ctx context.Context) {
        for _, v := range UTILS {
                go v(ctx)
        }
}
