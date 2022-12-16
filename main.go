package main

import (
	"git.garena.com/azmi-adhani/group-project-backend/cmd/cronservice"
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
	cronservice.Init()
	utils.RegisterTagName()
	utils.RegisterValidation()
	routers.Init(h)
}
