package main

import (
	"apc_compare/pkg/apc"
	"apc_compare/pkg/compare"
	"apc_compare/pkg/mail"
	"apc_compare/server/models"
	"apc_compare/server/setting"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

func init() {
	setting.Setup()
	models.Setup()
	apc.Setup()

}

func main() {
	//go build -o main.exe ./
	hadErr, err := compare.SpaceCompare()
	if err != nil {
		log.Fatal(err)
	}

	if hadErr {
		if err := mail.SendMail(); err != nil {
			log.Fatal()
		}
	}

}
