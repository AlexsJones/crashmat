/*********************************************************************************
*     File Name           :     api_controller.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 07:39]
*     Last Modified       :     [2015-09-29 08:36]
*     Description         :      
**********************************************************************************/
package main

import (
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/context"
  "log"
  "net/http"
)

type upload struct {
  applicationid string
  raw string
}

type uploadController struct {
  upload[] *upload
}

func (i *uploadController) Before(c context.Context) error {

  c.HttpResponseWriter().Header().Set("X-UploadController", "true")
  return nil
}

func (i *uploadController) Create(c context.Context) error {

  data, dataError := c.RequestData()
  if dataError != nil {
    log.Fatalf("Data error %s", dataError.Error())
    return goweb.API.RespondWithError(c, http.StatusInternalServerError, dataError.Error())
  }

  dataMap := data.(map[string]interface{})

  log.Printf(dataMap["applicationid"].(string))
  log.Printf(dataMap["raw"].(string))

  incomingData := new(upload)
  incomingData.applicationid = dataMap["applicationid"].(string)
  incomingData.raw = dataMap["raw"].(string)
  i.upload = append(i.upload, incomingData)
  return goweb.Respond.WithStatus(c, http.StatusCreated)
}

func (i *uploadController) ReadMany(c context.Context) error {

  if i.upload == nil {
    i.upload = make([] *upload, 0)
  }
  return goweb.API.RespondWithData(c,i.upload)
}

func (i *uploadController) Read(applicationid string, c context.Context) error {

  var results[] *upload

  for _, incomingData := range i.upload {
    if incomingData.applicationid == applicationid {
      results = append(results, incomingData)
    }
  }
  return goweb.API.RespondWithData(c, results)
}
