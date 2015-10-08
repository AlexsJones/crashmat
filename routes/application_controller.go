/*********************************************************************************
*     File Name           :     application_controller.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-07 11:28]
*     Last Modified       :     [2015-10-08 13:38]
*     Description         :      
**********************************************************************************/

package routes

import (
  "log"
  "net/http"
  "github.com/AlexsJones/crashmat/utils"
  "github.com/AlexsJones/crashmat/types"
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/context"
)
type applicationController struct {}

func (a *applicationController) ReadMany(c context.Context) error {

  var applications []types.Application

  types.DatabaseConnection.Find(&applications)

  var appId []string

  for _, elem := range applications {
    appId = append(appId,elem.ApplicationId)
  }

  return goweb.API.RespondWithData(c,appId)
}

func (a *applicationController) Delete(applicationid string,
c context.Context) error {
  
  log.Println("Delete")
  data, dataError := c.RequestData()
  if dataError != nil {
    log.Println("Data error %s", dataError.Error())
    return goweb.API.RespondWithError(c, http.StatusInternalServerError,
    dataError.Error())
  }

  dataMap := data.(map[string]interface{})

  basicd,ok := dataMap["authorization"]

  if basicd != "Basic" {
    log.Println("Post missing basic Authorization")
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Post missing correct Authorization")
  }

  userd, ok := dataMap["username"]

  if !ok {
    log.Println("Post missing username")
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Post missing username")
  }

  passd, ok := dataMap["password"]

  if !ok {
    log.Println("Post missing password")
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Post missing password")
  }

  var result types.Application

  types.DatabaseConnection.Where(&types.Application{ 
    ApplicationId:applicationid}).First(&result)

    if result.ApplicationId == applicationid {

      if result.Username != userd {
        log.Println("Post bad username")
        return goweb.API.RespondWithError(c, http.StatusBadRequest,
        "Bad credentials")
      }

      if utils.DoesPasswordMatchHash(result.EncryptedPassword,passd.(string))  {
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

  data, dataError := c.RequestData()
  if dataError != nil {
    log.Println("Data error %s", dataError.Error())
    return goweb.API.RespondWithError(c, http.StatusInternalServerError,
    dataError.Error())
  }

  dataMap := data.(map[string]interface{})

  log.Println(dataMap)

  basicd,ok := dataMap["authorization"]

  if !ok {
    log.Println("Post missing authorization")
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Post missing authorization")
  }
  /* Setting up to use basic auth */
  if basicd != "Basic" {
    log.Println("Post missing basic Authorization")
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Post missing correct Authorization")
  }

  appd, ok := dataMap["applicationid"]

  if !ok {
    log.Println("Post missing applicationid")
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Post missing applicationid")
  }

  userd, ok := dataMap["username"]

  if !ok {
    log.Println("Post missing username")
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Post missing username")
  }

  passd, ok := dataMap["password"]

  if !ok {
    log.Println("Post missing password")
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Post missing password")
  }

  hashedPassword := utils.PasswordToHash(passd.(string))

  app := types.NewApplication(appd.(string),userd.(string),hashedPassword) 
  
  /* Check and see whether we have already registred this app */

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

