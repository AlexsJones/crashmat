/*********************************************************************************
*     File Name           :     routes.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-25 09:51]
*     Last Modified       :     [2015-09-25 12:15]
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

  goweb.Map("GET","/api",func(c context.Context) error {
    return goweb.API.Respond(c,405,"Please use POST method only",nil)
  })
}
