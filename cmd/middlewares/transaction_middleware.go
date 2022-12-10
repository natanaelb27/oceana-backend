package middlewares

import (
	"log"
	"net/http"

	"git.garena.com/azmi-adhani/group-project-backend/internal/db"
	"github.com/gin-gonic/gin"
)

// checks if the given status is in the list
func StatusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

func DBTransactionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := db.Get()
		// begin database transaction
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		c.Set("db_trx", tx)
		c.Next()

		// run afterhandler function finishes
		if StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			// commit transactions
			if err := tx.Commit().Error; err != nil {
				log.Print("trx commit error: ", err)
			}
		} else {
			// rolling back transaction due to status code other than
			tx.Rollback()
		}
	}
}
