/*********************************************************************************
*     File Name           :     routes.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-25 09:51]
*     Last Modified       :     [2015-09-25 11:23]
*     Description         :
**********************************************************************************/
package main

import (
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/context"
  "log"
)

func mapRoutes() {

  goweb.MapBefore(func(c context.Context) error {
    log.Printf("%s %s %s", c.HttpRequest().RemoteAddr, c.MethodString(), c.HttpRequest().URL.Path)
    return nil
  })

  goweb.Map("GET", "/", func(c context.Context) error {
    return goweb.Respond.WithOK(c)
  })

}
