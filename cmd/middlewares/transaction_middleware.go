package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/natanaelb27/oceana-backend/internal/db"
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

		c.Set("tx", tx)
		c.Next()

		// run after handler function finishes
		if StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			// commit transactions
			if err := tx.Commit().Error; err != nil {
				log.Print("tx commit error: ", err)
			}
		} else {
			// rolling back transaction due to status code other than
			tx.Rollback()
		}
	}
}
