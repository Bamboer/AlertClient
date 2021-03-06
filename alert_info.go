package main

import (
        "net/url"
        "net/http"
        "time"
        "strconv"
        "io/ioutil"
        "path/filepath"
        "encoding/json"
        "grafana/pkg/client"
        "grafana/pkg/configer"
        "grafana/pkg/notification"
        "github.com/gomodule/redigo/redis"
)

var (
        alertDict   = map[int]client.SimpleInfo{}
        dbInfo      = map[string]map[string]string{}
        timeSeries  = map[int][]time.Time{}
        alertNum = 0
        RedisServer string
)

func run() {
        if err := CacheRe();err != nil{
          info.Println(err)
        }

        go client.FileServer()

        ticker := time.NewTicker(30 * time.Second)
        for _ = range ticker.C {
                if err := Alerter(); err != nil {
                        info.Println(err)
                }
        }
        info.Println("Agent Down.")
}

func Alerter() error {

        // Get alert list
        data, err := client.GetAlerts()
        if err != nil {
//                info.Println(err)
                return err
        }
        for _, alert := range data {
                if alert.State == "alerting" {
                        if _, ok := alertDict[alert.Id]; ok {
                                continue
                        }
                        m := client.SimpleInfo{}
                        m.Name = alert.Name
                        m.PanelId = alert.PanelId
                        m.DbUid = alert.DashboardUid
                        m.DbSlug = alert.DashboardSlug
                        //Scan the alerting data
                        for _, v := range alert.EvalData.EvalMatches {
                                m.AlertMetrics = m.AlertMetrics +" " + v.Metric
                                value,err := strconv.ParseFloat(fmt.Sprintf("%.2f",v.Value),32)
                                if err != nil{
                                    m.AlertValues = append(m.AlertValues,v.Value)
                                 }else{
                                    m.AlertValues = append(m.AlertValues,float32(value))
                                 }

                        }
                        //Get alert item detail info
                        alert_info, err := client.GetAlert(strconv.Itoa(alert.Id))
                        if err != nil {
                                info.Println("Get alert error: ", err)
                        }
                        m.OrgId = alert_info.OrgId
                        m.Frequency = alert_info.Frequency

                        if k, ok := dbInfo[m.DbUid]; ok {
                                m.TempVar = k
                        } else {
                                // Get Dashboard templating variables for render image url
                                db, err := client.GetDashboard(m.DbUid)
                                if err != nil {
                                        info.Println("Get Dashboard error: ", err)
                                }

                                s := map[string]string{}
                                Temvars := db.Dashboard.Templating["list"]
                                for _, tvar := range Temvars {
                                        if tvar.Current.Selected {
                                                s[tvar.Name] = tvar.Current.Text
                                        }
                                }
                                dbInfo[m.DbUid] = s
                                m.TempVar = s
                                if err := CacheSet("dbinfo",&dbInfo);err != nil{
                                   info.Println(err)
                                }
                        }
                        b, render_url,imgurl,err := RenderImage(m)
                        if err != nil {
                                info.Println(err)
                        }
                        m.RenderURL = render_url
                        m.ImgURL = imgurl
                        alertDict[alert.Id] = m
                        alertNum = len(alertDict)
                        m.AlertNum = &alertNum
                        notification.Emit("alerting", m, b)
                        if err := CacheSet("alertdict",&alertDict);err != nil{
                            info.Println(err)
                        }
                }
        }

        //Recovery alert item
        for alertId, alertV := range alertDict {
                alert_info, err := client.GetAlert(strconv.Itoa(alertId))
                if err != nil {
                        info.Println(err)
                        return err
                }
                if alert_info.State == "ok" {
                        info.Println("Recovery alert: ",alertV.Name)
                        b, render_url,imgurl, err := RenderImage(alertV)
                        if err != nil {
                                info.Println(err)
                        }
                        alertNum = len(alertDict) - 1
                        alertV.RenderURL = render_url
                        alertV.ImgURL = imgurl
                        notification.Emit("ok", alertV, b)
                        delete(alertDict, alertId)
                        if err := CacheSet("alertdict",&alertDict);err != nil{
                            info.Println(err)
                        }
                }
        }
        return nil
}

func RenderImage(m client.SimpleInfo) ([]byte, string, string, error) {
        //Generator time series for render image
        t1 := int(time.Now().Unix()) * 1000
        t2 := t1 - 3600000

        imgname := m.DbSlug + strconv.Itoa(m.PanelId) + ".png"
        grafana_conf := configer.ConfigParse()
        imgdir,err:= filepath.Abs(grafana_conf.ImgDir)
        if err != nil{
           info.Println(err)
        }
        imgpath := filepath.Join(imgdir,imgname)
        imgurl := "http://" + grafana_conf.ImgServerName + grafana_conf.ImgServerPort + "/" + imgname
        info.Println("imgurl: ",imgurl)

        uri, err := url.Parse(grafana_conf.Grafana_uri)
        if err != nil {
                info.Println("URL parse error: ", err)
        }
        token := "Bearer " + grafana_conf.Grafana_token
        uri.Path = "/render/d-solo/" + m.DbUid + "/" + m.DbSlug
        c,_ := NewRender(grafana_conf.Grafana_uri,token)
        req, err := http.NewRequest("GET", uri.String(), nil)
        if err != nil {
                info.Println("request generator error: ", err)
        }
        // request header add
        req.Header.Add("Authorization", token)
        req.Header.Add("Content-Type", "application/json")
        req.Header.Add("Accept", "application/json")
        //add query data
        q := req.URL.Query()
        q.Add("orgid", strconv.Itoa(m.OrgId))
        q.Add("from", strconv.Itoa(t2))
        q.Add("to", strconv.Itoa(t1))
        q.Add("panelId", strconv.Itoa(m.PanelId))
        q.Add("width", "1000")
        q.Add("height", "500")
        for k, v := range m.TempVar {
                q.Add(k, v)
        }
        req.URL.RawQuery = q.Encode()
        info.Println("Render url: ", req.URL.String())
        //request
        resp, err := c.client.Do(req)
        if err != nil {
                info.Println(err)

        }
        defer resp.Body.Close()
        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                info.Println("iotuil Read error: ", err)
                return nil, req.URL.String(),imgurl, err
        }
       if err := ioutil.WriteFile(imgpath,b,0644);err !=nil{
            info.Println(err)
       }
        return b, req.URL.String(),imgurl, nil
}

type Render struct{
  uri    *url.URL
  token  string
  client *http.Client
}

func NewRender(uri ,token string)(*Render,error){
   url,err := url.Parse(uri)
   if err != nil{
     info.Println("Error: ",err)
     return nil,err
   }
   token = "Bearer " + token
   return &Render{
       uri : url,
       token:  token,
       client: &http.Client{},
   },nil
}

func CacheRe()error{
   RedisServer = configer.ConfigParse().RedisServer
   c,err := redis.Dial("tcp",RedisServer)
   if err != nil{
     info.Println(err)
     return err
   }
   defer c.Close()
   if alertdict,err := c.Do("get","alertdict");err != nil{
      info.Println(err)
   }else{
      if v,ok := alertdict.([]byte);ok{
            json.Unmarshal(v,&alertDict)
      }
   }
   if dbinfo,err := c.Do("get","dbinfo");err != nil{
      info.Println(err)
   }else{
      if v,ok := dbinfo.([]byte);ok{
            json.Unmarshal(v,&dbInfo)
      }
   }
   return nil
}

func CacheSet(name string,v interface{})error{
   p,err := json.Marshal(v)
   if err != nil{
     info.Println(err)
   }
   c,err := redis.Dial("tcp",RedisServer)
   if err != nil{
     return err
   }
   defer c.Close()
   if _,err := c.Do("set",name,p);err != nil{
      return err
   }
   info.Println("set ",name,"to cache server successed.")
   return nil
}
