package main
import (
  "flag"
  "fmt"
)



var config = flag.String("config","alert_client.conf","configuration file.")

func main(){
  flag.Parse()
  fmt.Println("config: ",*config)
}
