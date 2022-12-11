package cronservice

import (
	"log"
)

func (c *Cron) Per15Seconds() {
	c.scheduler.Every(15).Seconds().Do(func() {
		// cron example, service method that will get roles from db and print its count to console
		err := c.roleService.CronRolesCount()
		if err != nil {
			log.Print("CRON: ERROR: roleService.CronExample:", err)
		}

		err = nil
		err = c.orderTrackingService.CronOrderTrackingUpdater()
		if err != nil {
			log.Print("CRON: ERROR: orderTrackingService.CronOrderTrackingUpdater:", err)
		}

		err = nil
		err = c.refundService.CronAutoApproveRefundByExpirationTime()
		if err != nil {
			log.Print("CRON: ERROR: refundService.CronAutoApproveRefundByExpirationTime:", err)
		}
		// add any services that will run per 15 seconds here
	})
}
