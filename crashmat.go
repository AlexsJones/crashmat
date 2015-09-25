/*********************************************************************************
*     File Name           :     crashmat.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-24 23:14]
*     Last Modified       :     [2015-09-25 11:49]
*     Description         :
**********************************************************************************/
package main

import (
  "github.com/stretchr/goweb"
  "log"
  "net"
  "net/http"
  "os"
  "os/signal"
  "time"
)

const (
  ConfigurationPath string = "conf/app.json"
)

func main() {

  configuration := NewConfiguration(ConfigurationPath)

  mapRoutes()

  log.Print("Initialising...")
  s := &http.Server{
    Addr:           configuration.Port,
    Handler:        goweb.DefaultHttpHandler(),
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  listener, listenErr := net.Listen("tcp", ":" + configuration.Port)

  log.Printf("  visit: %s", ":" + configuration.Port)
  if listenErr != nil {
    log.Fatalf("Could not listen: %s", listenErr)
  }

  go func() {
    for _ = range c {

      listener.Close()
    }

  }()
  log.Fatalf("Error in server: %s", s.Serve(listener))
}
