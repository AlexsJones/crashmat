/*********************************************************************************
*     File Name           :     database_connection_test.go
*     Created By          :     anon
Creation Date       :     [2015-10-02 08:39]
*     Last Modified       :     [2015-10-02 13:58]
*     Description         :      
**********************************************************************************/
package main

import (
  "database/sql"
  "gopkg.in/gorp.v1"
  "log"
  _ "github.com/mattn/go-sqlite3"
  "testing"
)

func initDb(t *testing.T) *gorp.DbMap {

  configuration := Configuration{}

  configuration.Load("conf/app.json")

  db, err := sql.Open("sqlite3",configuration.Json.Database.LocalPath)
  checkErr(err, "sql.Open failed")
  log.Println("Sucessfully connected to database")

  dbmap := &gorp.DbMap{ Db: db, Dialect: gorp.SqliteDialect{}}

  dbmap.AddTableWithName(Upload{},"upload_entries").SetKeys(true, "Id")

  err = dbmap.CreateTablesIfNotExists()
  checkErr(err, "Create tables failed")
  return dbmap
}

func TestDatabaseConnectionTest(t *testing.T) {

  dbmap := initDb(t)

  err := dbmap.TruncateTables()
  checkErr(err,"truncate tables failed")

  u := NewUpload("ApplicationIdExample","Loads of raw data")
  ut := NewUpload("ApplicationIdExample2","Loads of raw data 2")

  err = dbmap.Insert(&u)
  checkErr(err,"Insert failed")
  err = dbmap.Insert(&ut)
  checkErr(err,"Insert failed")

  /* Count */
  var uploads[] Upload

  _, err = dbmap.Select(&uploads,"select * from upload_entries")
  checkErr(err,"Select failed")

  for x, p := range uploads {
    log.Printf("%d: %v\n",x,p)
  } 

  count,err := dbmap.Delete(&u)
  log.Println(count)
  checkErr(err,"Delete fail")
  count,err = dbmap.Delete(&ut)
  log.Println(count)
  checkErr(err,"Delete fail")

  defer dbmap.Db.Close()

}

