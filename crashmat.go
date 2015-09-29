/*********************************************************************************
*     File Name           :     crashmat.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-24 23:14]
*     Last Modified       :     [2015-09-29 14:29]
*     Description         :
**********************************************************************************/
package main

import (
  "log"
)

func main() {

  var configuration = Configuration{}

  configuration.Load("conf/app.json")

  configuration.LoadAuth()

  mapRoutes(configuration)

  log.Print("Initialising...")

  configuration.LoadServer()

}
