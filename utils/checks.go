/*********************************************************************************
*     File Name           :     utils/checks.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-08 14:04]
*     Last Modified       :     [2015-10-08 15:41]
*     Description         :      
**********************************************************************************/

package utils

import "log"
import "github.com/stretchr/goweb/context"
import "strconv"

func CheckHeaderIsValidWithBasicAuth(c context.Context) (didpass bool, applicationid int, username string, password string) {

  data, dataError := c.RequestData()
  if dataError != nil {
    log.Println("Data error %s", dataError.Error())
    return false,0,"",""
  }

  dataMap := data.(map[string]interface{})

  appd, ok := dataMap["applicationid"]

  if !ok {
    log.Println("Header missing applicationid")
    return false,0,"",""
  }

  basicd,ok := dataMap["authorization"]

  if !ok {
    log.Println("Header missing authorization")
    return false,0,"",""
  }
  if basicd != "Basic" {
    log.Println("Header missing basic Authorization")
    return false,0,"",""
  }

  userd, ok := dataMap["username"]

  if !ok {
    log.Println("Header missing username")
    return false,0,"",""
  }

  passd, ok := dataMap["password"]

  if !ok {
    log.Println("Header missing password")
    return false,0,"",""
  }

  appi, _ := strconv.Atoi(appd.(string))
  return true,appi,userd.(string),passd.(string)
}

func CheckHeaderIsValidWithBasicAuthAndRawData(c context.Context) (didpass bool, applicationid int, username string, password string,raw string) {

  data, dataError := c.RequestData()
  if dataError != nil {
    log.Println("Data error %s", dataError.Error())
    return false,0,"","",""
  }

  dataMap := data.(map[string]interface{})

  rawd, ok := dataMap["raw"]
  if !ok {
    log.Println("Header missing raw data")
    return false,0,"","",""
  }
  valid,appid,userd,passd := CheckHeaderIsValidWithBasicAuth(c)
  return valid,appid,userd,passd,rawd.(string)
}
