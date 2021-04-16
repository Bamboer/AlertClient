package utils

import(
  "time"
  "bytes"
  "strings"
  "strconv"
  "context"
  "net/http"
  "net/url"
  "net/smtp"
  "io/ioutil"
  "encoding/json"
  "path/filepath"
  "html/template"
  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
  "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
  "grafana/pkg/configer"
  "grafana/pkg/notification"
)

func init(){
    UTILS = append(UTILS, DAU)
}

var(
  GraphiteURL = "http://10.50.24.197:7001/render"
  User = "8bf3584d"
  Pwd  = "6f30e011"
)

type Tres struct{
 Datapoints [][]float64
 Target  string
}

type DailyReport struct{
  Timer  time.Time
  WK     int
  WeekDay map[int]DayData
}

type DayData struct{
  Access  int
  Health  int
}

func DAU(ctx context.Context){
        info.Println("Start DAU Reporter.")
        if err := RenderDau();err != nil{
          info.Println(err)
        }
        for {
            if time.Now().Hour() == 0{
              if err := RenderDau();err != nil{
                 info.Println(err)
              }
            }
                select {
                case <-ctx.Done():
                        info.Println("done")
                        return
                default:
                }
                time.Sleep(1 * time.Hour)
        }
}


func MailSend(b []byte)error{
        buffer := bytes.NewBuffer(nil)
        conf := configer.ConfigParse()
        notifications  := strings.Split(conf.Notifications,",")
        notifications_cc := strings.Split(conf.Notifications_cc,",")
        notifications_bcc := strings.Split(conf.Notifications_bcc,",")
        message := &notification.Message{From: "SVoice " + conf.SmtpServer.Username,
                To:   notifications,
                Cc:   notifications_cc,
                Bcc:  notifications_bcc,
                Attachment: notification.Attachment{
                        WithFile:    true,
                        ContentType: "image/png",
                        Name:        "graph.png",
                },
        }
        m, _ := notification.NewMail(conf.SmtpServer.Username, conf.SmtpServer.Password, conf.SmtpServer.SmtpAddress, conf.SmtpServer.Port)
        boundary := "BamboerBoundary"
        Header := make(map[string]string)
        Header["From"] = message.From
        Header["To"] = strings.Join(message.To, ";")
        Header["Cc"] = strings.Join(message.Cc, ";")
        Header["Bcc"] = strings.Join(message.Bcc, ";")
        Header["Subject"] = "SVoice  Daily Report"
        Header["Content-Type"] = "multipart/related;boundary=" + boundary
//        Header["Mime-Version"] = "1.0"
//        Header["Date"] = time.Now().UTC().String()

        m.WriteHeader(buffer, Header)

        if message.Attachment.WithFile {
                attachment := "\r\n--" + boundary + "\r\n"
                attachment += "Content-Transfer-Encoding:base64\r\n"
                //                attachment += "Content-Disposition:attachment\r\n"
                attachment += "Content-Type:" + message.Attachment.ContentType + ";name=\"" + message.Attachment.Name + "\"\r\n"
                attachment += "Content-ID: <" + message.Attachment.Name + "> \r\n\r\n"
                buffer.WriteString(attachment)
                defer func() {
                        if err := recover(); err != nil {
                                info.Println("Error: ", err)
                        }
                }()
                m.WriteFile(buffer, b)
        }

        body := "\r\n--" + boundary + "\r\n"
        body += "Content-Type: text/html; charset=UTF-8 \r\n"
        buffer.WriteString(body)
//        body += renderMessage(state, imgsrc, msg)
        err := GHtml(buffer)
        if err != nil{
           return err
        }

        buffer.WriteString("\r\n--" + boundary + "--")
        if err := smtp.SendMail(m.Host+m.Port, m.Auth, m.User, message.To, buffer.Bytes());err !=nil{
            return err
        }
        return nil
}

func GHtml(b *bytes.Buffer)error{
  conf := configer.ConfigParse()
  elb := conf.AWSELBName
  region := conf.AWSRegion
  tpPath := conf.DauTpPath
  t := time.Now()
  td := int(t.Weekday())
  et := int(time.Date(t.Year(),t.Month(),t.Day(),0,0,0,0,time.Local).Unix())
  
  _,wk := t.ISOWeek()
  y,m,d := t.Date()
  h,M,s := t.Clock()

  DReport =  DailyReport{
         Timer: fmt.Sprintf("%d/%d/%d %d:%d:%d UTC",int(m),d,y,h,M,s),
         WK: wk,
         WeekDay: map[int]DayData{},
  }
  
  access,err := Access(region,elb)
  if err != nil{
    info.Println(err)
  }
  health,err := Health()
  if err != nil{
    info.Println(err)
  }



  DReport.WeekDay[td] = DayData{Health: health[et]*100 }
  for i := 1; i <= int(td);i++{
     t := time.Unix(int64(et - i*86400),0)
     wk := int(t.Weekday())
     DReport.WeekDay[wk] = DayData{Access: access[et - i*86400], Health: health[et - i*86400]*100}
  }
  for i := int(td);i < 6; i++{
   t := time.Unix(int64(et + (6-i)*86400),0)
   wk := int(t.Weekday())
   DReport.WeekDay[wk] = DayData{}
  }

  absPath,err := filepath.Abs(tpPath)
  if err != nil{
     info.Println(err)
  }
  tp,err := template.ParseFiles(absPath)
  if err != nil{
     info.Println(err)
  }
  if err := tp.Execute(b, DReport);err != nil{
     return err
  }
  return nil
}

func Access(region,elb string)(map[int]int,error){
  data := map[int]int{}
  t := time.Now()
  st := time.Date(t.Year(),t.Month(),t.Day()-6,0,0,0,0,time.Local)
  et := time.Date(t.Year(),t.Month(),t.Day(),0,0,0,0,time.Local)
  cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
  if err != nil {
          info.Println("unable to load SDK config, %v", err)
          return data,err
  }
  client := cloudwatch.NewFromConfig(cfg)
  input := &cloudwatch.GetMetricStatisticsInput{
      StartTime : aws.Time(st),
      EndTime : aws.Time(et),
      MetricName: aws.String("RequestCount"),
      Namespace: aws.String("AWS/ELB"),
      Period: aws.Int32(86400),
      Dimensions: []types.Dimension{{Name: aws.String("LoadBalancerName"),Value: aws.String(elb)}},
      Statistics: []types.Statistic{types.StatisticSum},
  }
  output,err := client.GetMetricStatistics(context.Background(),input)
  if err != nil {
    return data,err
  }
  for _,v := range output.Datapoints{
    data[int((*v.Timestamp).Unix())] = int(*v.Sum)
  }
  return data,nil
}

func Health()(map[int]int,error){
  gr := []Tres{}
  data := map[int]int{}
  c,_ := NewRender(GraphiteURL,User,Pwd)
  req,err := http.NewRequest("GET",c.uri.String(),nil)
  if err != nil{
    return data,err
  }
  req.SetBasicAuth(c.user,c.password)

  q := req.URL.Query()
  q.Add("target",`alias(summarize(averageSeries(ec2-cn-north-1-svoice-idg-rel.timers.application.dummy-client.*.*vdt.health.mean, *), "1d", "avg", false), "Overall ")`)
  q.Add("from","-144hours")
  q.Add("format","json")
  req.URL.RawQuery = q.Encode()
  resp, err := c.client.Do(req)
  if err != nil {
      return data,err
  }
  defer resp.Body.Close()

  b,err := ioutil.ReadAll(resp.Body)
  if err != nil{
      return data,err
  }
  if err := json.Unmarshal(b,&gr);err != nil{
      return data,err
  }
  for _, v := range gr[0].Datapoints{
      data[int(v[1])] = int(v[0])
  }
  return data,nil
}


type Render struct{
  uri    *url.URL
  user   string
  password string
  client *http.Client
}

func NewRender(uri, user,password string)(*Render,error){
   url,err := url.Parse(uri)
   if err != nil{
     info.Println("Error: ",err)
     return nil,err
   }
   return &Render{
       uri : url,
       user: user,
       password: password,
       client: &http.Client{},
   },nil
}

func RenderDau()error{
        //Generator time series for render image
        t1 := int(time.Now().Unix()) * 1000
        t2 := t1 - 86400000*7

        grafana_conf := configer.ConfigParse()
        uri, err := url.Parse(grafana_conf.Grafana_uri)
        if err != nil {
                info.Println("URL parse error: ", err)
        }
        token := "Bearer " + grafana_conf.Grafana_token
        uri.Path = "/render/d-solo/" + "000000221" + "/" + "user-count"
        c,_ := newrender(grafana_conf.Grafana_uri,token)
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
        q.Add("orgid", strconv.Itoa(2))
        q.Add("from", strconv.Itoa(t2))
        q.Add("to", strconv.Itoa(t1))
        q.Add("panelId", strconv.Itoa(15))
        q.Add("var-FRAM", "ec2-cn-north-1-svoice-idg-rel")
        q.Add("width", "1000")
        q.Add("height", "500")
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
                return  err
        }
        if err := MailSend(b);err != nil{
                info.Println(err)
                return err
        }
        return  nil
}
type render struct{
  uri    *url.URL
  token  string
  client *http.Client
}

func newrender(uri ,token string)(*render,error){
   url,err := url.Parse(uri)
   if err != nil{
     info.Println("Error: ",err)
     return nil,err
   }
   token = "Bearer " + token
   return &render{
       uri : url,
       token:  token,
       client: &http.Client{},
   },nil
}
