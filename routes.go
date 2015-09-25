/*********************************************************************************
*     File Name           :     routes.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-25 09:51]
*     Last Modified       :     [2015-09-25 15:51]
*     Description         :
**********************************************************************************/
package main

import (
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/context"
  "log"
  "net/http"
)

func generateApiRoutes() {
  goweb.Map("GET","/api",func(c context.Context) error {
    return goweb.API.Respond(c,http.StatusMethodNotAllowed,
    "Please use POST method only",nil)
  })
  /* Perform the auth */
  goweb.Map("/auth", func(c context.Context) error {
    log.Println("Starting authentication")
    provider, err := gomniauth.Provider("github")
    log.Println("Created new provider")
    if err != nil {
      return err
    }
    state := gomniauth.NewState("after","success")
    log.Println("Set to new state")
    authUrl, err := provider.GetBeginAuthURL(state,nil)
    log.Println("Getting auth url")
    if err != nil {
      return err
    }
    log.Println("Responding with redirect")
    return goweb.Respond.WithRedirect(c,authUrl)
  })
  /* Callback from github */
  goweb.Map("/auth/callback", func(c context.Context) error {
    log.Println("Authentication response")
    provider, err := gomniauth.Provider("github")
    if err != nil {
      log.Fatalf("Error with provider")
      return err
    }
    creds, err := provider.CompleteAuth(c.QueryParams())
    log.Println("Completing authentication")
    if err != nil {
      log.Fatalf("Error completing authentication")
      return err
    }
    log.Println("Getting user credentials")
    user, userErr := provider.GetUser(creds)
    if userErr != nil {
      log.Fatalf("Get user error")
      return userErr
    }
    return goweb.API.RespondWithData(c,user)
  })
}

func mapRoutes() {
  goweb.MapBefore(func(c context.Context) error {
    log.Printf("%s %s %s", c.HttpRequest().RemoteAddr, 
    c.MethodString(), c.HttpRequest().URL.Path)
    return nil
  })
  generateApiRoutes()
}
