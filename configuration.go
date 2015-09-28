/*********************************************************************************
*     File Name           :     configuration.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-25 11:33]
*     Last Modified       :     [2015-09-27 22:04]
*     Description         :      
**********************************************************************************/

package main

import (
  "os"
  "encoding/json"
  "log"
)

type Configuration struct {
  LocalDev bool
  Port string
  ClientSecret string
  ClientId string
  GithubAuthCallback string
}

func NewConfiguration(configurationPath string) *Configuration {
  conf,err := os.Open(configurationPath)
  if err != nil {
    log.Fatalf("opening configuration file",err.Error())
  }

  var configuration Configuration
  jsonParser := json.NewDecoder(conf)
  if err = jsonParser.Decode(&configuration); err != nil {
    log.Fatalf("parsing config file", err.Error())
  }
  return &configuration
}
