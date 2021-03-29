package main
import (
  "strconv"
  "time"
  "image/png"
  "grafana/pkg/client"
)

var (
  alertDict = map[int] alertinfo{}
  dbInfo    = map[string] map[string]string{}
)
type alertinfo struct {
  Name           string
  AlertMetrics    []string
  AlertValues     []interface{}
  PanelId        int
  OrgId          int
  DbUid          string
  DbSlug         string
  Frequency      int
  TempVar        map[string]string
}


func run() error{
  ticker := time.NewTicker(30 * time.Second)
  for _ = range ticker.C{

  }
}

func Alerting()error{
   data,err := client.GetAlerts()
   if err != nil{
     info.Println(err)
   }
   for _,alert := range(data){
      if alert.State == "alerting"{
         if _,ok := alertDict[alert.Id];ok{
           return nil
         }
         m := alertinfo{}
         m.Name = alert.Name
         m.PanelId = alert.PanelId
         m.DbUid  =  alert.DashboardUid
         m.DbSlug =  alert.DashboardSlug
         for _,v = range(alert.EvalData.EvalMatches){
            m.AlertMetrics = append(m.AlertMetrics,v.Metric)
            m.AlertValues = append(m.AlertValues,v.Value)
         }
         alert_info,err := client.GetAlert(strconv.Itoa(m.Id))
         if err != nil{
             info.Println("Get alert error: ",err)
         }
         m.OrgId = alert_info.OrgId
         m.Frequency = alert_info.Frequency
         alertDict[alert.Id] = m
         if k,ok := dbInfo[alert.DbUid]
      }
   }
}

func DbInfo(){
   data,err := client.GetDashboard(DashboardUid)
   if err != nil{
     info.Println(err)
   }

}

func RenderImage(dashboardUid,dashboardName,orgid,t1,t2,panelId string){
   url := host + "/render/d-solo/" + dashboardUid + "/" + dashboardName + "?orgid="+ orgid +"&from=" + t1 + "&to=" + t2 + "&panelId=" + panelId + "&width=1000&height=500&tz=Asia%2FShanghai"
   image_png := client.wget(url)
   os.Write(filename)
   return image_png
}
