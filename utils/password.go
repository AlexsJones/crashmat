/*********************************************************************************
*     File Name           :     utils/password.go
*     Created By          :     anon
*     Creation Date       :     [2015-10-07 16:42]
*     Last Modified       :     [2015-10-07 16:51]
*     Description         :      
**********************************************************************************/
package utils

import(
  "golang.org/x/crypto/bcrypt"
)

func PasswordToHash(password string) string {

  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),
  bcrypt.DefaultCost)

  if err != nil {
    panic(err)
  }
  return string(hashedPassword)
}

func DoesPasswordMatchHash(hashedPassword string, password string) bool {

  err := bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
  if err == nil {
    return true
  }
  return false
}
