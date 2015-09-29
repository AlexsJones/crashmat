/*********************************************************************************
*     File Name           :     upload_controller_test.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 09:12]
*     Last Modified       :     [2015-09-29 09:17]
*     Description         :      
**********************************************************************************/
package main
import (
  "github.com/stretchr/codecs/services"
  "github.com/stretchr/goweb"
  "github.com/stretchr/goweb/handlers"
  "github.com/stretchr/testify/assert"
  testifyhttp "github.com/stretchr/testify/http"
  "testing"
)
func TestUploadControllerGet(t *testing.T) {
  // make a test HttpHandler and use it
  codecService := new(services.WebCodecService)
  handler := handlers.NewHttpHandler(codecService)
  goweb.SetDefaultHttpHandler(handler)
  // call the target code
  var configuration = Configuration{}

  configuration.Load("conf/app.json")
  
  mapRoutes(configuration)
  
  goweb.Test(t, "GET /upload", func(t *testing.T, response *testifyhttp.TestResponseWriter) {
    
    assert.Equal(t, 200, response.StatusCode, "Status code should be correct")
  })
}

