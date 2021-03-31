package notification

import (
  "fmt"
  "bytes"
  "net/http"
  "encoding/json"
  "grafana/pkg/configer"
)


func init(){
  conf := configer.ConfigParse()
  SNS["dingding"] = Newdingding(*conf.Dingding)
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

func (d *dingding)Send(state string,alertNum int,msg interface{})error{
   data := make(map[string]interface{})
   data["msgtype"] = "text"
   data["at"] = map[string]interface{}{"atMobiles": reminders, "isAtAll": true}
   data["text"]= map[string]string{"content": msg}

   mdata,err := json.Marshal(data)
//   fmt.Println("mdata: ",mdata)
   if err != nil{
     info.Println("Marshal error: ",err)
   }
   reader := bytes.NewReader(mdata)
   req,err := http.NewRequest("POST",d.dApi,reader)
   req.Header.Set("Content-Type","application/json; charset=utf-8")
   resp,err := d.client.Do(req)
   defer resp.Body.Close()
   if err != nil{
      info.Println("err: ",err)
   }
   err = json.NewDecoder(resp.Body).Decode(&gr)
   if err != nil{
      info.Println("err: ",err)
   }
   return nil
}
