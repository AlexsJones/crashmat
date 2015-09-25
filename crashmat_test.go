/*********************************************************************************
*     File Name           :     crashmat_test.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-25 12:15]
*     Last Modified       :     [2015-09-25 12:18]
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
func TestRoutes(t *testing.T) {
  // make a test HttpHandler and use it
  codecService := new(services.WebCodecService)
  handler := handlers.NewHttpHandler(codecService)
  goweb.SetDefaultHttpHandler(handler)
  // call the target code
  mapRoutes()
  goweb.Test(t, "GET /api", func(t *testing.T, response *testifyhttp.TestResponseWriter) {
    // should not be allowed (405)
    assert.Equal(t, 405, response.StatusCode, "Status code should be correct")
  })
}
