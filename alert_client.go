package main
import (
  "io"
  "strings"
  "os"
  "fmt"
  "flag"
  "path"
//  "time"
//  "path/filepath"
//  "grafana/pkg/configer"
  "grafana/pkg/client"
  "log"
//  "grafana/pkg/notification"
)

var (
  info  *log.Logger
  Version = "v1"
  DashboardURL = "/api/dashboards/uid/"
  AlertPath = "/api/alerts/"
  OrgId = "/api/org"
  version = flag.Bool("version",false,"Print version.")
)


func init(){
  arg := path.Base(os.Args[0])
  logfile := strings.ToLower(arg + "1.log")
  file,err := os.OpenFile(logfile,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
  if err != nil{
     info.Println("Failed to open file: ",err)
  }
  info = log.New(io.MultiWriter(os.Stdout,file),"Info: ",log.Ldate|log.Ltime|log.Lshortfile)
}


/*
func run(alert_client,){
  ticker := time.NewTicker(30 * time.Second)
  for _ = range ticker.C{

  }
}*/

func main(){
  flag.Parse()
  if len(os.Args) >2 {
    if os.Args[1] == "version"{
    fmt.Println("grafana alert client ",Version)
    os.Exit(0)
    }
   data,err := client.GetAlerts()
   if err != nil{
      info.Println(err)
   }
   log.Println(data)
  }
}

