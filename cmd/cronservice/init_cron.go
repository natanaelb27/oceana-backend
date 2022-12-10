package cronservice

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"github.com/go-co-op/gocron"
)

func Init() {
	roleRepository := repositories.NewRoleRepository(&repositories.RoleRepositoryConfig{
		DB: db.Get(),
	})
	roleService := services.NewRoleService(&services.RoleServiceConfig{
		RoleRepository: roleRepository,
	})

	c := New(&CronConfig{
		Scheduler:   *gocron.NewScheduler(time.UTC),
		RoleService: roleService,
	})

	// cron that runs per 15 second
	c.Per15Seconds()

	c.scheduler.StartAsync()
}
