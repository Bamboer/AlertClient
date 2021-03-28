package main
import (
  "time"
  "grafana/pkg/client"
)

type alertinfo struct {
  Id   int  
  OrgId  int 
}

func run() error{
  ticker := time.NewTicker(30 * time.Second)
  for _ = range ticker.C{
      
  }
}

func Alerting()(alertlist,error){
   data,err := client.GetAlerts()
   if err != nil{
     info.Println(err)
   }
   for index,alert := range(data){
      if alert.State == "alerting"{

      }
   }
}

func DbInfo(){
   data,err := client.GetDashboard(DashboardUid)
   if err != nil{
     info.Println(err)
   }

}

func RenderImage(){
   url := host + "/render/d-solo/" + dashboardUid + "/" + dashboardName + "?orgid="+ orgid +"&from=" + t1 + "&to=" + t2 + "&panelId=" + panelId + "&width=1000&height=500&tz=Asia%2FShanghai"
   image_png := client.wget(url)
   os.Write(filename)
   return image_png
}
