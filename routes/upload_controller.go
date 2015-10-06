/*********************************************************************************
*     File Name           :     api_controller.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 07:39]
*     Last Modified       :     [2015-10-06 11:11]
*     Description         :      
**********************************************************************************/
package routes

import (
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/context"
  "log"
  "encoding/json"
  "net/http"
  elastigo "github.com/mattbaird/elastigo/lib"
  "github.com/AlexsJones/crashmat/types"
)

const (
  iname string = "crashmat"
)
type uploadController struct {}

func (i *uploadController) Before(c context.Context) error {

  c.HttpResponseWriter().Header().Set("X-types.UploadController", "true")
  return nil
}

func (i *uploadController) Create(c context.Context) error {

  data, dataError := c.RequestData()
  if dataError != nil {
    log.Fatalf("Data error %s", dataError.Error())
    return goweb.API.RespondWithError(c, http.StatusInternalServerError,
    dataError.Error())
  }

  dataMap := data.(map[string]interface{})

  if dataMap["applicationid"].(string) != "" {
    if dataMap["raw"].(string) != "" {

      uploaded := types.NewUpload(dataMap["applicationid"].(string), 
      dataMap["raw"].(string))
      err := types.DatabaseConnection.Insert(&uploaded)
      if err != nil {

        return goweb.API.RespondWithError(c, http.StatusInternalServerError,
        dataError.Error())
      }
    }
  }

  return goweb.API.RespondWithData(c,nil)
}

func (i *uploadController) ReadMany(c context.Context) error {

  var results []types.Upload
  qry := elastigo.Search(iname).Pretty().Query(
    elastigo.Query().All(),
  )
  out, err := qry.Result(types.ElasticConnection)
  
  if err != nil {
      log.Println("err querying elastic connection:%v", err)
    return goweb.API.RespondWithError(c, http.StatusInternalServerError,
    err.Error())
  }

  for _, elem := range out.Hits.Hits {
    bytes, err :=  elem.Source.MarshalJSON()
    if err != nil {
      log.Println("err calling marshalJson:%v", err)
    return goweb.API.RespondWithError(c, http.StatusInternalServerError,
    err.Error())
    }
    var t types.Upload
    json.Unmarshal(bytes, &t)
    results = append(results, t) 
  }

  return goweb.API.RespondWithData(c,results)
}

func (i *uploadController) Read(applicationid string, c context.Context) error {

  var results []types.Upload

  qry := elastigo.Search(iname).Pretty().Query(
    elastigo.Query().Search(applicationid),
  )
  out, err := qry.Result(types.ElasticConnection)
  if err != nil {
      log.Println("err querying elastic connection:%v", err)
    return goweb.API.RespondWithError(c, http.StatusInternalServerError,
    err.Error())
  }

  for _, elem := range out.Hits.Hits {
    bytes, err :=  elem.Source.MarshalJSON()
    if err != nil {
      log.Println("err calling marshalJson:%v", err)
    return goweb.API.RespondWithError(c, http.StatusInternalServerError,
    err.Error())
    }
    var t types.Upload
    json.Unmarshal(bytes, &t)
    if t.ApplicationId == applicationid {
      results = append(results, t) 
    }
  }

  return goweb.API.RespondWithData(c,results)
}
