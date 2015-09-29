/*********************************************************************************
*     File Name           :     crashmat.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-24 23:14]
*     Last Modified       :     [2015-09-29 11:33]
*     Description         :
**********************************************************************************/
package main

import (
  "log"
)

func main() {

  var configuration = Configuration{}

  configuration.Load("conf/app.json")

  if configuration.Json.Elastic.IsEnabled {
    configuration.LoadElasticSearch()
  }

  configuration.LoadAuth()

  mapRoutes(configuration)

  log.Print("Initialising...")

  configuration.LoadServer()

}
