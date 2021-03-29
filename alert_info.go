package main
import (
  "strconv"
  "time"
  "image/png"
  "grafana/pkg/client"
)

var (
  alertDict = map[int] alertinfo{}
  dbInfo    = map[string] []map[string]string{}
  timeSeries     = map[int][]time.Time {}
  alertingNum    = 0
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
  TempVar        []map[string]string
}


func run() error{
  ticker := time.NewTicker(30 * time.Second)
  for _ = range ticker.C{

  }
}

func Alerting()error{
//Recovery alert item
   for alerId,alertV := range(alertDict){
       alert_info,err := client.GetAlert(strconv.Itoa(alertId))
       if err != nil {
          info.Println(err)
          return err
       }
       if *alert_info.State == "ok"{
           recovery(alertId,alertV)
       }
   }  

// Get alert list
   data,err := client.GetAlerts()
   if err != nil{
     info.Println(err)
     return err
   }
   for _,alert := range(data){
      if alert.State == "alerting"{
         if _,ok := alertDict[alert.Id];ok{
            continue
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
//Get alert item detail info
         alert_info,err := client.GetAlert(strconv.Itoa(alert.Id))
         if err != nil{
             info.Println("Get alert error: ",err)
         }
         m.OrgId = *alert_info.OrgId
         m.Frequency = *alert_info.Frequency

         if k,ok := dbInfo[m.DbUid];ok{
            m.TmpVar = k 
         }else{
            db,err := client.GetDashboard(m.DbUid)
            if err != nil{
               info.Println("Get Dashboard error: ",err)
            }
            dbInfo[m.DbUid] = *db.Dashboard.Templating["list"]
            m.TempVar = dbInfo[m.DbUid]
         }
         alertDict[alert.Id] = m
         t1 = time.
         t2 = time.
         RenderImage(m,t1,t2)
         send(m)
      }
   }
}


func RenderImage(m *alertinfo,t1,t2 time.Time){
   url := host + "/render/d-solo/" + dashboardUid + "/" + dashboardName + "?orgid="+ orgid +"&from=" + t1 + "&to=" + t2 + "&panelId=" + panelId + "&width=1000&height=500&tz=Asia%2FShanghai"
   image_png := client.wget(url)
   os.Write(filename)
   return image_png
}
