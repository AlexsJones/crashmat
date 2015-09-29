/*********************************************************************************
*     File Name           :     api_controller.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 07:39]
*     Last Modified       :     [2015-09-29 16:28]
*     Description         :      
**********************************************************************************/
package main

import (
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/context"
  "log"
  "encoding/json"
  "net/http"
  elastigo "github.com/mattbaird/elastigo/lib"

)

type uploadController struct {
  upload[] *Upload
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

  uploaded := NewUpload(dataMap["applicationid"].(string), dataMap["raw"].(string))

  elasticConnection.Index("crashmat","upload",NewGuid(),nil,uploaded) 

  return goweb.API.RespondWithData(c,nil)
}

func (i *uploadController) ReadMany(c context.Context) error {

  var results []Upload

  qry := elastigo.Search("crashmat").Pretty().Query(
    elastigo.Query().All(),
  )
  out, err := qry.Result(elasticConnection)
  if err != nil {
    log.Fatal(err)
  }

  count := 0
  for count < out.Hits.Total {

    bytes, err :=  out.Hits.Hits[count].Source.MarshalJSON()
    if err != nil {
      log.Fatalf("err calling marshalJson:%v", err)
    }

    var t Upload
    json.Unmarshal(bytes, &t)
    results = append(results, t) 
    count += 1  
  }
  return goweb.API.RespondWithData(c,results)
}

func (i *uploadController) Read(applicationid string, c context.Context) error {

  var results []Upload

  qry := elastigo.Search("crashmat").Pretty().Query(
    elastigo.Query().Search(applicationid),
  )
  out, err := qry.Result(elasticConnection)
  if err != nil {
    log.Fatal(err)
  }
  count := 0
  for count < out.Hits.Total {

    bytes, err :=  out.Hits.Hits[count].Source.MarshalJSON()
    if err != nil {
      log.Fatalf("err calling marshalJson:%v", err)
    }

    var t Upload
    json.Unmarshal(bytes, &t)
    results = append(results, t) 
    count += 1
  }
  return goweb.API.RespondWithData(c,results)
}
