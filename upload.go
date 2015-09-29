/*********************************************************************************
*     File Name           :     upload.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 14:40]
*     Last Modified       :     [2015-09-29 17:48]
*     Description         :      
**********************************************************************************/
package main

import (
  "time"
  "encoding/json"
)
type Upload struct {
  Applicationid string `json:"Applicationid"`
  Date time.Time `json:"Date"`
  Raw string  `json:"Raw"`
}

func NewUpload(applicationid string, raw string) Upload {
  return Upload{Date:time.Now(), Applicationid:applicationid, Raw:raw}
}

func (u *Upload) String() string {
  b, _ := json.Marshal(u)
  return string(b)
}
