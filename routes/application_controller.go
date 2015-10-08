/*********************************************************************************
*     File Name           :     application_controller.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-07 11:28]
*     Last Modified       :     [2015-10-08 16:54]
*     Description         :      
**********************************************************************************/

package routes

import (
  "log"
  "net/http"
  "fmt"
  "github.com/AlexsJones/crashmat/utils"
  "github.com/AlexsJones/crashmat/types"
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/context"
)
type applicationController struct {}

func (a *applicationController) ReadMany(c context.Context) error {

  var applications []types.Application

  types.DatabaseConnection.Find(&applications)

  appId := make(map[string]string)

  for _, elem := range applications {
    appId[fmt.Sprintf("%v",elem.Id)] = fmt.Sprintf("%v",elem.ApplicationId)
  }

  return goweb.API.RespondWithData(c,appId)
}

func (a *applicationController) DeleteMany(c context.Context) error {

  isValid,appd,userd,passd := utils.CheckHeaderIsValidWithBasicAuth(c)

  if isValid == false {
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Bad request in POST header")
  }

  var result types.Application

  types.DatabaseConnection.Where(&types.Application{ 
    ApplicationId:appd}).First(&result)

    if result.ApplicationId == appd{

      if result.Username != userd {
        log.Println("Post bad username")
        return goweb.API.RespondWithError(c, http.StatusBadRequest,
        "Bad credentials")
      }

      if utils.DoesPasswordMatchHash(result.EncryptedPassword,passd)  {
        log.Println("Password matches for post")

        types.DatabaseConnection.Delete(&result)
      }else {
        log.Println("Post bad password")
        return goweb.API.RespondWithError(c, http.StatusBadRequest,
        "Bad credentials")
      }
    } else {

      log.Println("Application not found")
      return goweb.API.RespondWithError(c, http.StatusBadRequest,
      "Application not found")
    }

    return goweb.API.RespondWithData(c,nil)
  }

  func (a *applicationController) Create(c context.Context) error {

    isValid,appd,userd,passd := utils.CheckHeaderIsValidWithBasicAuth(c)

    if isValid == false {
      return goweb.API.RespondWithError(c, http.StatusBadRequest,
      "Bad request in POST header")
    }

    hashedPassword := utils.PasswordToHash(passd)

    app := types.NewApplication(appd,userd,hashedPassword) 

    var results types.Application

    types.DatabaseConnection.Where(&types.Application{ 
      ApplicationId:app.ApplicationId}).First(&results)

      if results.ApplicationId == app.ApplicationId {
        log.Println("Found existing application registered with this id")  
        return goweb.API.RespondWithError(c, http.StatusBadRequest,
        "Application already exists")
      }else {
        types.DatabaseConnection.Create(&app)
      }

      return goweb.API.RespondWithData(c,nil)
    }

