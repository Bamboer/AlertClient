package client

import (
        "encoding/json"
        "gopkg.in/ini.v1"
        "io"
        "log"
        "net/http"
        "net/url"
        "os"
        "path"
        "strings"
        "grafana/pkg/configer"
)

var (
  info *log.Logger
)

func init() {
        arg := path.Base(os.Args[0])
        logfile := strings.ToLower(arg + ".log")
        file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
                log.Println("Failed to open file: ", err)
        }
        cfg, err := ini.Load(*configer.ConfigFile)
        if err != nil {
                info.Println("Fail to read file: ", err)
                os.Exit(1)
        }
        mode := cfg.Section("").Key("mode").In("dev", []string{"dev", "debug", "prd"})
        if mode == "dev" || mode == "debug" {
                info = log.New(io.MultiWriter(os.Stdout, file), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
        } else if mode == "prd" {
                info = log.New(file, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
        }
}

type grafana_client struct {
        uri        *url.URL
        token      string
        httpClient *http.Client
}

func NewGrafanaClient(uri, token string) (*grafana_client, error) {
        url, err := url.Parse(uri)
        token = "Bearer " + token
        //        info.Println(url)
        //        info.Println(token)
        if err != nil {
                info.Println(err)
                return nil, err
        }
        return &grafana_client{
                uri:        url,
                token:      token,
                httpClient: &http.Client{},
        }, nil
}

func (c *grafana_client) Get(path string, v interface{}) error {
        uri := c.uri
        uri.Path = path
        req, err := http.NewRequest("GET", uri.String(), nil)
        if err != nil {
                info.Println(err)
                return err
        }
        req.Header.Add("Authorization", c.token)
        //        req.SetBasicAuth("admin", "admin")
        req.Header.Add("Content-Type", "application/json")
        req.Header.Add("Accept", "application/json")
        resp, err := c.httpClient.Do(req)
        if err != nil {
                info.Println(err)
                return err
        }
        defer resp.Body.Close()
        if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
                info.Println(err)
                return err
        }
        //        log.Println(v)
        return nil
}
