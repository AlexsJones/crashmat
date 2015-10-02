/*********************************************************************************
*     File Name           :     crashmat.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-24 23:14]
*     Last Modified       :     [2015-10-02 12:14]
*     Description         :
**********************************************************************************/
package main

import (
  "log"
)

func main() {

  var configuration = Configuration{}

  configuration.Load("conf/app.json")

  configuration.LoadElasticSearch()

  log.Print("Initialising Elasticsearch")

  configuration.LoadAuth()

  log.Print("Initialising auth")

  mapRoutes(configuration)

  log.Print("Map routes")

  configuration.LoadDatabase()

  log.Print("Initialising Database")

  configuration.LoadServer()

  defer configuration.DbMap.Db.Close()
}
