/*********************************************************************************
*     File Name           :     routes.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-25 09:51]
*     Last Modified       :     [2015-09-25 17:07]
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
  goweb.Map("/auth/{provider}", func(c context.Context) error {
    log.Println("Starting authentication")
    provider, err := gomniauth.Provider(c.PathValue("provider"))
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
  /* Callback from auth */
  goweb.Map("/auth/{provider}/callback", func(c context.Context) error {
    log.Println("Authentication response")
    provider, err := gomniauth.Provider(c.PathValue("provider"))
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

    
    log.Println("Authenticated successfully!")
    log.Println(user)
    return goweb.API.Respond(c,200,"Successfully authenticated",nil)
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
