package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"log"
)

//
var db *gorm.DB

var server = "10.10.135.235"
var port = 1433
var user = "exl_user_V2"
var password = "hwtjdaup"
var database = "RPT_DB"

//var server = "192.168.1.41"
//var port = 1433
//var user = "sa"
//var password = "Zz01470147"
//var database = "master"

// Setup initializes the database instance
func Setup() {
	var err error

	db, err = gorm.Open("mssql", fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable",
		server, user, password, port, database))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	db.SingularTable(true)

}
