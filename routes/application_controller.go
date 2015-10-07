/*********************************************************************************
*     File Name           :     application_controller.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-07 11:28]
*     Last Modified       :     [2015-10-07 16:53]
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

func (a *applicationController) Create(c context.Context) error {

  data, dataError := c.RequestData()
  if dataError != nil {
    log.Println("Data error %s", dataError.Error())
    return goweb.API.RespondWithError(c, http.StatusInternalServerError,
    dataError.Error())
  }

  dataMap := data.(map[string]interface{})

  basicd,ok := dataMap["Authorization"]

  if !ok {
    log.Println("Post missing applicationid")
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Post missing applicationid")
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

  return goweb.API.RespondWithData(c,nil)
}

