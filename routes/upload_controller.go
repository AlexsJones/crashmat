/*********************************************************************************
*     File Name           :     api_controller.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 07:39]
*     Last Modified       :     [2015-10-08 15:17]
*     Description         :      
**********************************************************************************/
package routes

import (
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/context"
  "log"
  "fmt"
  "encoding/json"
  "strconv"
  "net/http"
  elastigo "github.com/mattbaird/elastigo/lib"
  "github.com/AlexsJones/crashmat/types"
  "github.com/AlexsJones/crashmat/utils"
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

  isValid,appd,userd,passd,rawd := utils.CheckHeaderIsValidWithBasicAuthAndRawData(c)

  if isValid == false {
    return goweb.API.RespondWithError(c, http.StatusBadRequest,
    "Bad request in POST header")
  }

  var result types.Application

  types.DatabaseConnection.Where(&types.Application{ 
    ApplicationId:appd}).First(&result)

    if result.ApplicationId == appd {

      uploaded := types.NewUpload(appd, 
      rawd)

      if result.Username != userd {
        log.Println("Post bad username")
        return goweb.API.RespondWithError(c, http.StatusBadRequest,
        "Bad credentials")
      }

      if utils.DoesPasswordMatchHash(result.EncryptedPassword,passd)  {
        log.Println("Password matches for post")

        types.DatabaseConnection.Create(&uploaded)

        /* Update this into the application */

        result.Uploads = append(result.Uploads,uploaded)

        types.DatabaseConnection.Model(&result).Updates(types.Application {
          Uploads:result.Uploads})

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
  func (i *uploadController) ReadMany(c context.Context) error {

    var results []types.Upload
    qry := elastigo.Search(iname).Pretty().Query(
      elastigo.Query().All(),
    )
    out, err := qry.Result(types.ElasticConnection)
    if err != nil {
      fmt.Println("err querying elastic connection:%v", err)
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
      appi,_ := strconv.Atoi(applicationid)
      if t.ApplicationID == appi{
        results = append(results, t) 
      }
    }

    return goweb.API.RespondWithData(c,results)
  }
