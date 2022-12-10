package cronservice

import "log"

func (c *Cron) Per15Seconds() {
	c.scheduler.Every(15).Seconds().Do(func() {
		// cron example, service method that will get roles from db and print its count to console
		err := c.roleService.CronRolesCount()
		if err != nil {
			log.Print("CRON: roleService.CronExample:", err)
		}

		// add any services that will run per 15 seconds here
	})

	c.scheduler.StartAsync()
}
