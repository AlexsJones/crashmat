/*********************************************************************************
*     File Name           :     configuration.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-25 11:33]
*     Last Modified       :     [2015-09-28 17:36]
*     Description         :      
**********************************************************************************/

package main

import (
  "os/signal"
  "time"
  "net"
  "net/http"
  "github.com/stretchr/goweb"
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/signature"
  "github.com/stretchr/gomniauth/providers/github"
  "os"
  "encoding/json"
  "log"
  "fmt"
  "github.com/olivere/elastic"
)

type Elastic struct {
  HostAddress string
}
type Json struct {
  LocalDev bool
  Port string
  ClientSecret string
  ClientId string
  GithubAuthCallback string
  Elastic Elastic
}
type Configuration struct {
  Json Json
  ElasticClient *elastic.Client
  HttpServer *http.Server
}

func (c *Configuration)Load(configurationPath string) {

  conf,err := os.Open(configurationPath)
  if err != nil {
    log.Fatalf("opening configuration file",err.Error())
  }

  jsonParser := json.NewDecoder(conf)
  if err = jsonParser.Decode(&c.Json); err != nil {
    log.Fatalf("parsing config file", err.Error())
  }
}

func (c *Configuration) LoadElasticSearch() {

  log.Printf("Connecting to elastic search %s\n", c.Json.Elastic.HostAddress)

  client, err := elastic.NewClient(elastic.SetURL(c.Json.Elastic.HostAddress))
  if err != nil {
    log.Fatal(`Could not connect to the Elasticsearch service - 
    Please make sure configuration is correct`)
    panic(err)
  }

  info, code, err := client.Ping().Do()
  if err != nil {
    log.Fatal(`Could not connect to the Elasticsearch service -
    Please make sure configuration is correct`)
    panic(err)
  }

  fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

  c.ElasticClient = client
}

func (c *Configuration) LoadAuth() {

  gomniauth.SetSecurityKey(signature.RandomKey(64))
  gomniauth.WithProviders(github.New(c.Json.ClientId,
  c.Json.ClientSecret,
  c.Json.GithubAuthCallback))
}

func (c *Configuration) LoadServer() {

  port := c.Json.Port
  if os.Getenv("PORT") != "" {
    port = os.Getenv("PORT")  
    log.Print("Using environmental variable for $PORT")
  }

  c.HttpServer = &http.Server{
    Addr:           port,
    Handler:        goweb.DefaultHttpHandler(),
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  b := make(chan os.Signal, 1)

  signal.Notify(b, os.Interrupt)
  listener, listenErr := net.Listen("tcp", ":" + port)
  log.Printf("  visit: %s", ":" + port)
  if listenErr != nil {
    log.Fatalf("Could not listen: %s", listenErr)
  }

  go func() {
    for _ = range b {
      listener.Close()
    }

  }()
  log.Fatalf("Error in server: %s", c.HttpServer.Serve(listener))
}
