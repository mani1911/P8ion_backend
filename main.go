package main

import (
	"p8ion/config"
	"p8ion/database"
	"p8ion/server"
	"p8ion/server/model"
)

func main() {
	config.InitConfig()
	database.Connect()
	model.MigrateDB()
	server.Run()
}
