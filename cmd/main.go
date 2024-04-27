package main

import (
	"log"
	"restful/app"
	"restful/cmd/server"

	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := app.MysqlDb(mysql.Config{
		User:   "root",
		Passwd: "r23password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "restful_new",

		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	server := server.New(":3000", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
