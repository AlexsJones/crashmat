/*********************************************************************************
*     File Name           :     types/application.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-07 16:23]
*     Last Modified       :     [2015-10-08 15:44]
*     Description         :      
**********************************************************************************/
package types

import (
  "encoding/json"
  "time"
)

type Application struct {
  Id int64
  Created int64
  ApplicationId int 
  Username string
  EncryptedPassword string
  Uploads []Upload
}

func NewApplication(applicationid int, username string, encryptedPassword string) Application {
  return Application{
    Created:time.Now().UnixNano(),
    ApplicationId:applicationid,
    Username:username,
    EncryptedPassword:encryptedPassword,
  }
}

func (u *Application) String() string {
  b, _ := json.Marshal(u)
  return string(b)
}
