package main

import (
	"io/ioutil"
	"strconv"
	"time"
	//  "image/png"
	"grafana/pkg/client"
	"grafana/pkg/configer"
	"grafana/pkg/notification"
)

var (
	alertDict   = map[int] client.SimpleInfo{}
	dbInfo      = map[string]map[string]string{}
	timeSeries  = map[int][]time.Time{}
	alertingNum = 0
)


func run() error {
	ticker := time.NewTicker(30 * time.Second)
	for _ = range ticker.C {
		if err := Alerter(); err != nil {
			info.Println(err)
		}
	}
}

func Alerter() error {

	// Get alert list
	data, err := client.GetAlerts()
	if err != nil {
		info.Println(err)
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
			for _, v = range alert.EvalData.EvalMatches {
				m.AlertMetrics = append(m.AlertMetrics, v.Metric)
				m.AlertValues = append(m.AlertValues, v.Value)
			}
			//Get alert item detail info
			alert_info, err := client.GetAlert(strconv.Itoa(alert.Id))
			if err != nil {
				info.Println("Get alert error: ", err)
			}
			m.OrgId = *alert_info.OrgId
			m.Frequency = *alert_info.Frequency

			if k, ok := dbInfo[m.DbUid]; ok {
				m.TmpVar = k
			} else {
				// Get Dashboard templating variables for render image url
				db, err := client.GetDashboard(m.DbUid)
				if err != nil {
					info.Println("Get Dashboard error: ", err)
				}

				s := map[string]string{}
				Temvars = *db.Dashboard.Templating["list"]
				for _, tvar := range Temvars {
					if tvar.Current.Selected {
						s[tvar.Name] = tvar.Current.Text
					}
				}
				dbInfo[m.DbUid] = s
				m.TempVar = s
			}
			alertNum = len(alertDict)
                        m.AlertingNum = &alertNum
			b,render_url, err := RenderImage(m)
			if err != nil {
				info.Println(err)
			}
                        m.RenderURL = render_url
			alertDict[alert.Id] = m
			notification.Emit("alerting", m, b)
		}
	}

	//Recovery alert item
	for alerId, alertV := range alertDict {
		alert_info, err := client.GetAlert(strconv.Itoa(alertId))
		if err != nil {
			info.Println(err)
			return err
		}
		if *alert_info.State == "ok" {
			b,render_url,err := RenderImage(alertV)
			if err != nil {
				info.Println(err)
			}
			alertNum = len(alertDict) - 1 
                        alertV.RenderURL = render_url
			notification.Emit("ok", alertV, b)
                        delete(alertDict,alertId)
		}
	}
}

func RenderImage(m client.SimpleInfo) ([]byte, string, error) {
	//Generator time series for render image
	t1 := int(time.Now().Unix()) * 1000
	t2 := t1 - 3600000
	grafana_conf := configer.ConfigParse()

	uri, err := url.Parse(grafana_conf.Grafana_uri)
	if err != nil {
		info.Println("URL parse error: ", err)
	}
	token := "Bearer " + grafana_conf.Grafana_token
	uri.Path = "/render/d-solo/" + m.DbUid + "/" + m.DbSlug
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
	info.Println("url > ", req.URL.String())
	//request
	resp, err := c.client.Do(req)
	if err != nil {
		info.Println(err)

	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		info.Println("iotuil Read error: ", err)
		return nil, req.URL.String(), err
	}
	return b, req.URL.String(), nil
}
