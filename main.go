package main

import (
	"github.com/natanaelb27/oceana-backend/cmd/cronservice"
	"github.com/natanaelb27/oceana-backend/cmd/handlers"
	"github.com/natanaelb27/oceana-backend/cmd/routers"
	"github.com/natanaelb27/oceana-backend/internal/db"
	"github.com/natanaelb27/oceana-backend/internal/utils"
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
