/*********************************************************************************
*     File Name           :     utils/checks.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-08 14:04]
*     Last Modified       :     [2015-10-08 14:35]
*     Description         :      
**********************************************************************************/

package utils

import "log"
import "github.com/stretchr/goweb/context"

func CheckHeaderIsValidWithBasicAuth(c context.Context) (didpass bool, applicationid string, username string, password string) {

  data, dataError := c.RequestData()
  if dataError != nil {
    log.Println("Data error %s", dataError.Error())
    return false,"","",""
  }

  dataMap := data.(map[string]interface{})

  appd, ok := dataMap["applicationid"]

  if !ok {
    log.Println("Header missing applicationid")
    return false,"","",""
  }

  basicd,ok := dataMap["authorization"]

  if !ok {
    log.Println("Header missing authorization")
    return false,"","",""
  }
  if basicd != "Basic" {
    log.Println("Header missing basic Authorization")
    return false,"","",""
  }

  userd, ok := dataMap["username"]

  if !ok {
    log.Println("Header missing username")
    return false,"","",""
  }

  passd, ok := dataMap["password"]

  if !ok {
    log.Println("Header missing password")
    return false,"","",""
  }

  return true,appd.(string),userd.(string),passd.(string)
}

func CheckHeaderIsValidWithBasicAuthAndRawData(c context.Context) (didpass bool, applicationid string, username string, password string,raw string) {

  data, dataError := c.RequestData()
  if dataError != nil {
    log.Println("Data error %s", dataError.Error())
    return false,"","","",""
  }

  dataMap := data.(map[string]interface{})

  rawd, ok := dataMap["raw"]
  if !ok {
    log.Println("Header missing raw data")
    return false,"","","",""
  }
  valid,appid,userd,passd := CheckHeaderIsValidWithBasicAuth(c)
  return valid,appid,userd,passd,rawd.(string)
}
