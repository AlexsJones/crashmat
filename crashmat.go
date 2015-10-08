/*********************************************************************************
*     File Name           :     crashmat.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-24 23:14]
*     Last Modified       :     [2015-10-08 17:06]
*     Description         :
**********************************************************************************/
package main

import (
  "github.com/AlexsJones/crashmat/types"
  "flag"
  "log"
  "github.com/AlexsJones/crashmat/routes"
)

func main() {

  var confFlag = flag.String("conf","","Path to configuration file")

  flag.Parse()

  if *confFlag == "" {
    log.Fatal("Please provide a conf path -conf") 
    return
  }
  var configuration = types.NewConfiguration(*confFlag)

  log.Print("Map routes")
  routes.MapRoutes()

  log.Print("Initialising Elasticsearch")
  configuration.StartElasticSearch()

  log.Print("Initialising auth")
  configuration.StartAuth()

  log.Print("Initialising Database")
  configuration.StartDatabase()

  log.Print("Starting periodic fetch service...")
  configuration.StartPeriodicFetch()

  configuration.StartServer()

}
