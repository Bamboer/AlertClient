package client
import (
//  "fmt"
  "encoding/json"
  "net/http"
  "net/url"
)


var  gr interface{}

type grafana_client struct {
  uri  *url.URL
  token string
  httpClient *http.Client
}

func NewGrafanaClient(uri,token string)(*grafana_client,error){
  url,err := url.Parse(uri)
  token = "Bearer " + token
  if err != nil{
    return nil,err
  }
  return &grafana_client{
    uri: url,
    token: token,
    httpClient: &http.Client{},
  },nil
}

func (c *grafana_client) Get(path string) (interface{}, error) {
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
//        fmt.Println("Head: ",req)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()
        err = json.NewDecoder(resp.Body).Decode(&gr)
 //       fmt.Println(resp)
        if err != nil {
                return nil, err
        }
        return gr, nil
}
