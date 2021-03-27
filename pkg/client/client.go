package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

type grafana_client struct {
	uri        *url.URL
	token      string
	httpClient *http.Client
}

func NewGrafanaClient(uri, token string) (*grafana_client, error) {
	url, err := url.Parse(uri)
	token = "Bearer " + token
	if err != nil {
		return nil, err
	}
	return &grafana_client{
		uri:        url,
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

func (c *grafana_client) Get(path string) (io.ReadCloser, error) {
	uri := c.uri
	uri.Path = path
	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.token)
	//        req.SetBasicAuth("admin", "admin")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//        err = json.NewDecoder(resp.Body).Decode(&gr)
	//        if err != nil {
	//                return nil, err
	//        }
	return resp.Body, nil
}

/*
func main(){
  C,_ := NewGrafanaClient("http://192.168.16.127:3000","eyJrIjoiZExNdVNiR3VaamdHSkxmNnNWNDdORnY2bXEyODBMT1IiLCJuIjoidGVzdCIsImlkIjoxfQ==")
//  data,err := C.Get("/api/dashboards/uid/CX9BS_wMk")
//  data,err := C.Get("/api/org")
  data,err := C.Get("/api/alerts")
  if err != nil{
    fmt.Println("client error: ",err)
  }
  switch value := (*data).(type){
  case map[string]interface{}:
    fmt.Println("map :",value)

  case []interface{}:
//    fmt.Println("slice: ",value)
    for _,i := range(value){
//      fmt.Println(k," : ",i)
      if v,ok := i.(map[string]interface{});ok{
//        fmt.Println(v["evalData"])
        if k,ok := v["evalData"].(map[string]interface{});ok{
          fmt.Println("K: ",k["evalMatches"])
          if s,ok := k["evalMatches"].([]interface{});ok{
                t := reflect.TypeOf(s[0])
                fmt.Println(t.String())
                fmt.Println("ok")
                fmt.Println("s ",s)
          }
        }
      }
    }
  default:
    fmt.Println("default: ",value)
  }
}*/
