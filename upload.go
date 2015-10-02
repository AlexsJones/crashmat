/*********************************************************************************
*     File Name           :     upload.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 14:40]
*     Last Modified       :     [2015-10-02 14:06]
*     Description         :      
**********************************************************************************/
package main

import (
  "time"
  "encoding/json"
)
type Upload struct {
  Id int64 `db:"upload_id"`
  Created int64
  ApplicationId string `json:"Applicationid" db:"ApplicationId"`
  RawData string `json:"RawData" db:"RawData"`
}

func NewUpload(applicationid string, raw string) Upload {
  return Upload{Created:time.Now().UnixNano(), ApplicationId:applicationid, RawData:raw}
}

func (u *Upload) String() string {
  b, _ := json.Marshal(u)
  return string(b)
}
