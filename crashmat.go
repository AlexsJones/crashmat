/*********************************************************************************
*     File Name           :     crashmat.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-24 23:14]
*     Last Modified       :     [2015-09-25 14:22]
*     Description         :
**********************************************************************************/
package main

import (
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/gomniauth/providers/github"
  "github.com/stretchr/signature"
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
  /* Auth */
  gomniauth.SetSecurityKey(signature.RandomKey(64))
  gomniauth.WithProviders(github.New("5f587d942465196d23af",
  configuration.ClientSecret,
  configuration.GithubAuthCallback))
  /* Auth */

  /*Port and TCP connection */
  port := configuration.Port
  if os.Getenv("PORT") != "" {
    port = os.Getenv("PORT")  
    log.Print("Using environmental variable for $PORT")
  }

  mapRoutes()

  log.Print("Initialising...")
  s := &http.Server{
    Addr:           port,
    Handler:        goweb.DefaultHttpHandler(),
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  listener, listenErr := net.Listen("tcp", ":" + port)
  /*Port and TCP connection */
  
  log.Printf("  visit: %s", ":" + port)
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
