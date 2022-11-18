package main

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/handlers"
	"git.garena.com/azmi-adhani/group-project-backend/cmd/routers"
	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
)

func main() {
	utils.InitEnv()
	db.Connect()
	db.RedisConnect()
	h := handlers.Init()
	utils.RegisterTagName()
	routers.Init(h)
}
