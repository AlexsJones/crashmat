/*********************************************************************************
*     File Name           :     crashmat.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-24 23:14]
*     Last Modified       :     [2015-10-05 19:21]
*     Description         :
**********************************************************************************/
package main

import (
  "github.com/AlexsJones/crashmat/types"
  "log"
  "github.com/AlexsJones/crashmat/routes"
)

func main() {

  var configuration = types.NewConfiguration("conf/app.json")

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

  defer configuration.DbMap.Db.Close()
}
