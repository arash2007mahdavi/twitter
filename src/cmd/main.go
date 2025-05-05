package main

import (
	"twitter/src/configs"
	"twitter/src/servers"
)

func main() {
	cfg := configs.GetConfig()
	servers.Init_Server(cfg)
}