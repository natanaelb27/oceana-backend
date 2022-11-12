package main

import (
	"fmt"

	"git.garena.com/azmi-adhani/group-project-backend/cmd/routers"
	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	"github.com/google/uuid"
)

func main() {
	utils.InitEnv()
	db.Connect()
	db := db.Get()
	role := models.Role{
		ID:   uuid.New(),
		Name: "Test",
	}
	result := db.Create(&role)
	fmt.Println(result.Error)
	routers.InitRouters()
}
