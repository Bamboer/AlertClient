package notification

import (
  "fmt"
  "bytes"
  "net/http"
  "encoding/json"
  "grafana/pkg/client"
  "grafana/pkg/configer"
)


func init(){
  SNS["dingding"] = DSend
}

var (
  reminders []string
  gr  interface{}
)

type dingding struct{
  dApi string
  client  *http.Client
}

func Newdingding(url string)*dingding{
   return &dingding{
      dApi: url,
      client: &http.Client{},
   }
}

func DSend(state string,msg client.SimpleInfo,b []byte)error{
  conf := configer.ConfigParse()
  dclient := Newdingding(*conf.Dingding)
  if err := dclient.Send(state,msg,b);err != nil{
     info.Println(err)
  }
}
func Text(msg string)error{
  conf := configer.ConfigParse()
  dclient := Newdingding(*conf.Dingding)
  if err := dclient.SendText(msg);err != nil{
     info.Println(err)
  }
}

func (d *dingding)SendText(msg string)error{
   data := make(map[string]interface{})
   data["msgtype"] = "text"
   data["at"] = map[string]interface{}{"atMobiles": reminders, "isAtAll": true}
   data["text"]= map[string]string{"content":  msg}

   mdata,err := json.Marshal(data)
   if err != nil{
     info.Println("Marshal error: ",err)
     return err
   }
   reader := bytes.NewReader(mdata)
   req,err := http.NewRequest("POST",d.dApi,reader)
   req.Header.Set("Content-Type","application/json; charset=utf-8")
   resp,err := d.client.Do(req)
   defer resp.Body.Close()
   if err != nil{
      info.Println("err: ",err)
      return err
   }
   err = json.NewDecoder(resp.Body).Decode(&gr)
   if err != nil{
      info.Println("err: ",err)
      return err
   }
   return nil
}

func (d *dingding)Send(state string,msg client.SimpleInfo,b []byte)error{
// state: alert status
// msg: send message body
// b: png format image 
  data := make(map[string]interface{})
  data["msgtype"] = "markdown"
  data["at"] = map[string]interface{}{"atMobiles": reminders, "isAtAll": true}
  data["markdown"]= d.RenderMsg(state,msg)
  
  mdata,err := json.Marshal(data)
  if err != nil{
     info.Println("Marshal error: ",err)
     return err
   }
  reader := bytes.NewReader(mdata)
  req,err := http.NewRequest("POST",d.dApi,reader)
  req.Header.Set("Content-Type","application/json; charset=utf-8")
  resp,err := d.client.Do(req)
  defer resp.Body.Close()
  if err != nil{
     info.Println("err: ",err)
     return err
  }
  err = json.NewDecoder(resp.Body).Decode(&gr)
  if err != nil{
     info.Println("err: ",err)
     return err
  }
   return nil
}
func (d *dingding)RenderMsg(state string,msg client.SimpleInfo)string{
  var content map[string]string
  if state == "alerting"{
      content = map[string]string{"title": "Alarm","text": fmt.Sprintf("### Alarm: %s\n> 1.Metric: %s\n> 2.Value: %s\n> 3.Dashboard: %s\n> 4.Alerting: %d\n> 5.Time: %s\n> ![screenshot](%s)\n> [详情](%s)\n",msg.Name, msg.AlertMetrics, msg.AlertValues, msg.DbSlug, *msg.AlertingNum, time.Now().UTC().String(),msg.RenderURL,msg.RenderURL)}
  }else if state == "ok"{
      content = map[string]string{"title": "Recovery","text": fmt.Sprintf("### Alarm: %s Recovery !\n> 1.Metric: %s\n> 2.Value: %s\n> 3.Dashboard: %s\n> 4.Alerting: %d\n> 5.Time: %s\n> ![screenshot](%s)\n> [详情](%s)\n",msg.Name, msg.AlertMetrics, msg.AlertValues, msg.DbSlug, *msg.AlertingNum, time.Now().UTC().String(),msg.RenderURL,msg.RenderURL)}
  }
}
