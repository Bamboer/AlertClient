package main
import (
  "os"
  "fmt"
  "time"
  "path/filepath"
  "grafana/pkg/configer"
  "grafana/pkg/client"
  "grafana/pkg/log"
  "grafana/pkg/notification"
)

var (
  ConfigFile string
  Version   bool
)
var (
  Version = "v1"
  DEFALTCONF = "alert_client.conf"
)

func init(){
  ConfigFile := flag.String("config","alert_client.conf","alert client configuration file.")
  Version := flag.Bool("version",false,"Print version.")

  if os.Args[1] == "versions"{
    fmt.Println("grafana alert client ",Version)
    os.Exit(0)
  }
}

func run(){
  ticker := time.NewTicker(30 * time.Second)
  for _ = range ticker.C{

  }
}

func main(){
  flag.Parse()
  client,_ := client.NewGrafanaClient("http://10.x.43.17x:3000","eyJrIjoia1AyeWhsbUpyMTE2VU41TlBCRnZHcVNOFoyM1FEcWIiLCJuIjoidGVzdDEiLCJpZCI6MX0=")
  fmt.Println(data)
}
