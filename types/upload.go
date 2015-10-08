/*********************************************************************************
*     File Name           :     upload.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 14:40]
*     Last Modified       :     [2015-10-08 15:43]
*     Description         :      
**********************************************************************************/
package types

import (
  "time"
  "encoding/json"
)
type Upload struct {
  Id int64 
  Created int64
  RawData string
  ApplicationID int `sql:"index"`
}

func NewUpload(applicationid int, raw string) Upload {
  return Upload{Created:time.Now().UnixNano(), 
  ApplicationID:applicationid, RawData:raw}
}

func (u *Upload) String() string {
  b, _ := json.Marshal(u)
  return string(b)
}
