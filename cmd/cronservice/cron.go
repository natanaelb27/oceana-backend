package cronservice

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/services"
	"github.com/go-co-op/gocron"
)

type Cron struct {
	scheduler   gocron.Scheduler
	roleService services.RoleService
}

type CronConfig struct {
	Scheduler   gocron.Scheduler
	RoleService services.RoleService
}

func New(c *CronConfig) *Cron {
	return &Cron{
		scheduler:   c.Scheduler,
		roleService: c.RoleService,
	}
}
