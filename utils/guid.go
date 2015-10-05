/*********************************************************************************
*     File Name           :     guid.go
*     Created By          :     anon
*     Creation Date       :     [2015-09-29 15:06]
*     Last Modified       :     [2015-10-05 16:36]
*     Description         :      
**********************************************************************************/

package utils

import (
  "fmt"
  "os/exec"
  "log"
)

func NewGuid() string {
  out, err := exec.Command("uuidgen").Output()
  if err != nil {
    log.Fatal(err)
    return ""
  }
  return fmt.Sprintf("%s",out)
}
