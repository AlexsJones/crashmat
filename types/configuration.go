/*********************************************************************************
*     File Name           :     types/configuration.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-05 15:36]
*     Last Modified       :     [2015-10-05 20:13]
*     Description         :      
**********************************************************************************/
package types

import (
  "gopkg.in/gorp.v1"
  "time"
  "database/sql"
  "github.com/stretchr/goweb"
  "fmt"
  _ "github.com/mattn/go-sqlite3"
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/signature"
  "github.com/stretchr/gomniauth/providers/github"
  "os"
  "log"
  "github.com/AlexsJones/crashmat/utils"
  "encoding/json"
  "net"
  "net/http"
  "os/signal"
  elastigo "github.com/mattbaird/elastigo/lib"
)
const (
  iname string = "crashmat"
)

type Elastic struct {
  IsEnabled bool
  HostAddress string
}
type FetchUpdate struct {
  MillisecondFrequency int
}
type Json struct {
  Port string
  ClientSecret string
  ClientId string
  GithubAuthCallback string
  Elastic Elastic
  Database Database
  FetchUpdate FetchUpdate
}

type Database struct {
  LocalPath string
  TableName string
  SelectOnRange string
}

type Configuration struct {
  Json Json
  HttpServer *http.Server
  HttpPort string
  Listener net.Listener
  DbMap *gorp.DbMap
  ElasticConnection *elastigo.Conn
}

/* TODO:Until I know how to write Go better I'll store a ref here to ES */
/* Global */
var ElasticConnection *elastigo.Conn
var DatabaseConnection *gorp.DbMap
/* Global */

func parseJson(configurationPath string) Json {

  var j Json
  conf,err := os.Open(configurationPath)
  if err != nil {
    log.Fatalf("opening configuration file",err.Error())
  }

  jsonParser := json.NewDecoder(conf)
  if err = jsonParser.Decode(&j); err != nil {
    log.Fatalf("parsing config file", err.Error())
  }
  return j
}

func (c *Configuration)StartElasticSearch() {

  elasticConnection := elastigo.NewConn()

  address := c.Json.Elastic.HostAddress

  if os.Getenv("ElasticHostAddress") != "" {
    address = os.Getenv("ElasticHostAddress")
  }

  elasticConnection.SetFromUrl(address)

  ElasticConnection = elasticConnection
  c.ElasticConnection = elasticConnection
}

func (c *Configuration)StartAuth() {

  clientSecret := c.Json.ClientSecret
  if os.Getenv("ClientSecret") != "" {
    clientSecret = os.Getenv("ClientSecret")
  }

  clientId := c.Json.ClientId
  if os.Getenv("ClientId") != "" {
    clientId = os.Getenv("ClientId")
  }

  gomniauth.SetSecurityKey(signature.RandomKey(64))
  gomniauth.WithProviders(github.New(clientId,
  clientSecret,
  c.Json.GithubAuthCallback))
}

func (c *Configuration) StartServer() {

  port := c.Json.Port
  if os.Getenv("PORT") != "" {
    port = os.Getenv("PORT")  
    log.Print("Using environmental variable for $PORT")
  }

  c.HttpPort = port

  c.HttpServer = &http.Server{
    Addr:           port,
    Handler:        goweb.DefaultHttpHandler(),
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  b := make(chan os.Signal, 1)

  signal.Notify(b, os.Interrupt)
  listener, listenErr := net.Listen("tcp", ":" + c.HttpPort)
  log.Printf("  visit: %s", ":" + c.HttpPort)
  if listenErr != nil {
    log.Fatalf("Could not listen: %s", listenErr)
  }
  c.Listener = listener

  go func() {
    for _ = range b {
      c.Listener.Close()
    }

  }()
  log.Fatalf("Error in server: %s", c.HttpServer.Serve(listener))
}

func (c *Configuration) StartDatabase() {

  db, err := sql.Open("sqlite3", c.Json.Database.LocalPath)
  utils.CheckErr(err, "sql.Open failed")
  dbmap := &gorp.DbMap{ Db: db, Dialect: gorp.SqliteDialect{}}
  dbmap.AddTableWithName(Upload{},c.Json.Database.TableName).SetKeys(true, "Id")
  err = dbmap.CreateTablesIfNotExists()
  utils.CheckErr(err, "Create tables failed")
  c.DbMap = dbmap
  DatabaseConnection = c.DbMap
}

func NewConfiguration(configurationPath string) Configuration {

  c := Configuration{
    Json:parseJson(configurationPath),
  }
  return c
} 

func (c *Configuration)fetchLastIndexFromES() int64 {
  var results []Upload

  qry := elastigo.Search(iname).Pretty().Query(
    elastigo.Query().All(),
  )
  out, err := qry.Result(ElasticConnection)

  utils.CheckErr(err,"Elastic Connection")    

  if out.Hits.Total == 0 {
    log.Println("No indice data for updating fetch information")
    return 0
  }

  for _, elem := range out.Hits.Hits {
    bytes, err :=  elem.Source.MarshalJSON()
    if err != nil {
      log.Fatalf("err calling marshalJson:%v", err)
    }
    var t Upload
    json.Unmarshal(bytes, &t)
    results = append(results, t) 
  }

  var highestId int64 = 0

  for _,p := range results {
    if p.Id >= highestId{
      highestId = p.Id
    }
  }
  return highestId
}

func (c *Configuration) StartPeriodicFetch() {

  log.Println("Starting periodic fetch")

  go func() {

    updateFrequency := c.Json.FetchUpdate.MillisecondFrequency

    var chunkSize int64 = 10
    for {
      var StartIndex = c.fetchLastIndexFromES()

      log.Println("Starting index from ",StartIndex)

      var uploads[] Upload

      _, err := DatabaseConnection.Select(&uploads, 
      fmt.Sprintf(c.Json.Database.SelectOnRange, 
      StartIndex,StartIndex + chunkSize))
  
      utils.CheckErr(err,"Select failed")      

      if len(uploads) == 0 {
        log.Printf("Nothing to parse from database")       
      }else{
        for x, p := range uploads {
          log.Printf("%d: %v\n",x,p)
          ElasticConnection.Index("crashmat","upload",utils.NewGuid(),nil,p)

        }
      }
      time.Sleep(time.Duration(updateFrequency) * time.Millisecond)
    }
  }()
}
