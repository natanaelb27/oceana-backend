package utils

import (
	"os"
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

func NewRefreshToken(user *models.User) (refreshToken *models.RefreshToken) {
	refreshTokenDaysDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS_DURATION"))
	refreshTokenMonthsDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MONTHS_DURATION"))
	refreshTokenYearsDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_YEARS_DURATION"))

	refreshToken = &models.RefreshToken{
		UserID:    user.ID,
		Code:      uuid.New().String(),
		ExpiresAt: time.Now().AddDate(refreshTokenYearsDuration, refreshTokenMonthsDuration, refreshTokenDaysDuration),
	}
	return
}
