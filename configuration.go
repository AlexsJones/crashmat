/*********************************************************************************
*     File Name           :     configuration.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-25 11:33]
*     Last Modified       :     [2015-10-02 12:11]
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
  "database/sql"
  "gopkg.in/gorp.v1"
  _ "github.com/mattn/go-sqlite3"
  elastigo "github.com/mattbaird/elastigo/lib"
)

type Elastic struct {
  IsEnabled bool
  HostAddress string
}

type Json struct {
  Port string
  ClientSecret string
  ClientId string
  GithubAuthCallback string
  Elastic Elastic
  Database Database
}

type Database struct {
  LocalPath string
}

type Configuration struct {
  Json Json
  HttpServer *http.Server
  DbMap *gorp.DbMap
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

/* TODO:Until I know how to write Go better I'll store a ref here to ES */
/* Global */
var elasticConnection *elastigo.Conn
/* Global */

func (c *Configuration) LoadElasticSearch() {

  elasticConnection = elastigo.NewConn()
  elasticConnection.SetFromUrl(c.Json.Elastic.HostAddress)

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

func (c *Configuration) LoadDatabase() {

  db, err := sql.Open("sqlite3", c.Json.Database.LocalPath)
  checkErr(err, "sql.Open failed")
  dbmap := &gorp.DbMap{ Db: db, Dialect: gorp.SqliteDialect{}}
  dbmap.AddTableWithName(Upload{},"upload_entries").SetKeys(true, "Id")
  err = dbmap.CreateTablesIfNotExists()
  checkErr(err, "Create tables failed")
  c.DbMap = dbmap
  //c.DbMap.TruncateTables()
}
func checkErr(err error, msg string) {
  if err != nil {
    log.Fatalln(msg, err)
  }
}
