/*********************************************************************************
*     File Name           :     routes.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-25 09:51]
*     Last Modified       :     [2015-10-14 16:01]
*     Description         :
**********************************************************************************/
package routes

import (
  "bytes"
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/context"
  "html/template"
  "io/ioutil"
  "log"
)

func generateAuthRoutes() {
  /* Perform the auth */
  goweb.Map("/auth/{provider}", func(c context.Context) error {
    log.Println("Starting authentication")
    provider, err := gomniauth.Provider(c.PathValue("provider"))
    log.Println("Created new provider")
    if err != nil {
      return err
    }
    state := gomniauth.NewState("after", "success")
    log.Println("Set to new state")
    authUrl, err := provider.GetBeginAuthURL(state, nil)
    log.Println("Getting auth url")
    if err != nil {
      return err
    }
    log.Println("Responding with redirect")
    return goweb.Respond.WithRedirect(c, authUrl)
  })
  /* Callback from auth */
  goweb.Map("/auth/{provider}/callback", func(c context.Context) error {
    log.Println("Authentication response")
    provider, err := gomniauth.Provider(c.PathValue("provider"))
    if err != nil {
      log.Fatalf("Error with provider")
      return goweb.Respond.WithRedirect(c, "/auth/status/failed")
    }
    creds, err := provider.CompleteAuth(c.QueryParams())
    log.Println("Completing authentication")
    if err != nil {
      log.Fatalf("Error completing authentication")
      return goweb.Respond.WithRedirect(c, "/auth/status/failed")
    }
    log.Println("Getting user credentials")
    user, userErr := provider.GetUser(creds)
    if userErr != nil {
      log.Fatalf("Get user error")
      return goweb.Respond.WithRedirect(c, "/auth/status/failed")
    }

    log.Println("Authenticated successfully!")
    log.Println("Username: %s User email: %s", user.Name(), user.Email())
    return goweb.Respond.WithRedirect(c, "/auth/status/successful")
  })
  /* Complete auth notification */
  goweb.Map("/auth/status/successful", func(c context.Context) error {

    return goweb.Respond.With(c, 200, []byte("Authentication completed successfully"))
  })
  /* Failed auth notification */
  goweb.Map("/auth/status/failed", func(c context.Context) error {
    return goweb.Respond.With(c, 400, []byte("Authentication failed"))
  })
}

func generateControllers() {

  uploadController := new(uploadController)
  goweb.MapController(uploadController)

  applicationController := new(applicationController)
  goweb.MapController(applicationController)
}

type Content struct{
  Body string
}

func generateHomePage() {

  goweb.Map("/", func(c context.Context) error {
    s1, _ := template.ParseFiles("tmpl/header.tmpl", 
    "tmpl/content.tmpl", "tmpl/footer.tmpl")

    fbody, err := ioutil.ReadFile("views/home.html")
    if err != nil {

    }
    var buffer bytes.Buffer
    s1.ExecuteTemplate(&buffer, "header", nil)
    s1.ExecuteTemplate(&buffer, "content", template.HTML(string(fbody)))
    s1.ExecuteTemplate(&buffer, "footer", nil)
    return goweb.Respond.With(c, 200, buffer.Bytes())
  })

}
func MapRoutes() {
  goweb.MapBefore(func(c context.Context) error {
    log.Printf("%s %s %s", c.HttpRequest().RemoteAddr,
    c.MethodString(), c.HttpRequest().URL.Path)
    return nil
  })

  generateAuthRoutes()

  generateHomePage()

  generateControllers()
}
