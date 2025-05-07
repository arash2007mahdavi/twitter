package main

import (
	"twitter/src/cache"
	"twitter/src/configs"
	"twitter/src/database"
	"twitter/src/database/migrations"
	"twitter/src/servers"
)

func main() {
	cfg := configs.GetConfig()
	err := database.InitDB(cfg)
	if err != nil {
		panic(err)
	}
	defer database.CloseDB()
	err = cache.InitRedis(cfg)
	if err != nil {
		panic(err)
	}
	defer cache.CloseRedis()
	migrations.Up1()
	servers.Init_Server(cfg)
}