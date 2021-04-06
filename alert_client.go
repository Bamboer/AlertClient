package main

import (
        "flag"
        "fmt"
        "io"
        "os"
        "log"
        "path"
        "context"
        "strings"
        "grafana/pkg/utils"
        "grafana/pkg/configer"
        "gopkg.in/ini.v1"
)

var (
        info    *log.Logger
        Version = "v1"
        version = flag.Bool("version", false, "Print version.")
)

func init() {
  arg := path.Base(os.Args[0])
  logfile := strings.ToLower(arg + ".log")
  file,err := os.OpenFile(logfile,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
  if err != nil{
     log.Println("Failed to open file: ",err)
  }
  cfg,err := ini.Load(*configer.ConfigFile)
  if err != nil{
    info.Println("Fail to read file: ",err)
    os.Exit(1)
  }
  mode := cfg.Section("").Key("mode").In("dev",[]string{"dev","debug","prd"})
  if mode == "dev" || mode == "debug"{
    info = log.New(io.MultiWriter(os.Stdout,file),"Info: ",log.Ldate|log.Ltime|log.Lshortfile)
  }else if mode == "prd"{
    info = log.New(file,"Info: ",log.Ldate|log.Ltime|log.Lshortfile)
  }
}

func main() {
        flag.Parse()
        if len(os.Args) > 2 {
                fmt.Println(len(os.Args))
                if os.Args[1] == "version" {
                        info.Println("grafana alert client ", Version)
                        os.Exit(0)
                }
        }
        //   data,err := client.GetDashboard(DashboardUid)
        //   data,err := client.GetAlerts()
        /*data, err := client.GetOrg()
        if err != nil {
                info.Println(err)
        }
        fmt.Println(data)*/
      ctx,cancel := context.WithCancel(context.Background())
      defer cancel()

      utils.UtilRuner(ctx)
      run()
}
