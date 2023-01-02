package services

import (
	"os"
	"strconv"
	"strings"
	"time"

	internalDtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"git.garena.com/azmi-adhani/group-project-backend/internal/utils"
	externalDtos "git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	WithTx(tx *gorm.DB) *authService
	Register(userCreateRequest *externalDtos.UserCreateRequest) (*externalDtos.UserCreateResponse, error)
	Login(c *gin.Context, userLoginRequest *externalDtos.UserLoginRequest) (*internalDtos.UserLoginData, error)
	RegisterWithGoogle(c *gin.Context, userRegisterWithGoogleRequest *externalDtos.UserRegisterWithGoogleRequest) (*internalDtos.UserLoginData, error)
	LoginWithGoogle(c *gin.Context, userLoginWithGoogleRequest *externalDtos.UserLoginWithGoogleRequest) (*internalDtos.UserLoginData, error)
	StepUpPassword(c *gin.Context, id uuid.UUID, request *externalDtos.AuthStepUpPasswordRequest) (token string, err error)
	SendResetPasswordEmail(c *gin.Context, request *externalDtos.UserVerifyResetPasswordRequest) (bool, error)
	ResetPassword(c *gin.Context, resetPasswordID string, request *externalDtos.ResetPasswordRequest) (*externalDtos.ResetPasswordResponse, error)
	RefreshAccessToken(c *gin.Context, refreshToken string, userID uuid.UUID) (string, error)
}

type authService struct {
	userRepository   r.UserRepository
	emailRepository  r.EmailRepository
	walletRepository r.WalletRepository
	redisRepository  r.RedisRepository
}

type AuthServiceConfig struct {
	UserRepository   r.UserRepository
	EmailRepository  r.EmailRepository
	WalletRepository r.WalletRepository
	RedisRepository  r.RedisRepository
}

func NewAuthService(c *AuthServiceConfig) AuthService {
	return &authService{
		userRepository:   c.UserRepository,
		emailRepository:  c.EmailRepository,
		walletRepository: c.WalletRepository,
		redisRepository:  c.RedisRepository,
	}
}

func (u authService) WithTx(tx *gorm.DB) *authService {
	u.userRepository = u.userRepository.WithTx(tx)
	u.emailRepository = u.emailRepository.WithTx(tx)
	u.walletRepository = u.walletRepository.WithTx(tx)

	return &u
}

func (u *authService) Register(userCreateRequest *externalDtos.UserCreateRequest) (*externalDtos.UserCreateResponse, error) {
	_, err := u.emailRepository.GetByEmail(userCreateRequest.Email)
	if err != e.ErrDataNotFound {
		return nil, e.ErrEmailAddressAlreadyExist
	}
	_, err = u.userRepository.GetByUsername(userCreateRequest.Username)
	if err != e.ErrDataNotFound {
		return nil, e.ErrUsernameAlreadyExist
	}
	_, err = u.userRepository.GetByPhoneNumber(userCreateRequest.PhoneNumber)
	if err != e.ErrDataNotFound {
		return nil, e.ErrPhoneNumberAlreadyExist
	}
	email := &models.Email{
		Address: userCreateRequest.Email,
	}
	createdEmail, err := u.emailRepository.Create(email)
	if err != nil {
		return nil, err
	}
	user := userCreateRequest.ToModel()
	user.EmailID = email.ID
	if userCreateRequest.Photo != nil {
		photoFile, err := userCreateRequest.Photo.Open()
		if err != nil {
			return nil, err
		}
		defer photoFile.Close()
		imageUrl, err := utils.UploadMedia(photoFile)
		if err != nil {
			return nil, err
		}
		user.Photo = &imageUrl
	}
	hashedPassword, err := utils.Hash(userCreateRequest.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	createdUser, err := u.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	wallet := &models.Wallet{
		UserID:          createdUser.ID,
		Balance:         0,
		Pin:             "",
		IsActive:        false,
		InvalidPinCount: 0,
	}
	_, err = u.walletRepository.Create(wallet)
	if err != nil {
		return nil, err
	}
	responseData := externalDtos.NewUserCreateResponse(createdUser)
	responseData.Email = createdEmail.Address
	return responseData, nil
}

func (u *authService) Login(c *gin.Context, userLoginRequest *externalDtos.UserLoginRequest) (*internalDtos.UserLoginData, error) {
	if userLoginRequest.Password == "" {
		// to prevent users who refreshed/closed the page during username/phone-number/password form
		// in register-google from logging in with no password
		return nil, e.ErrInvalidPassword
	}
	user, err := u.userRepository.GetByUsername(userLoginRequest.LoginCredentials)
	if err == nil {
		userLoginData, err := utils.AuthenticateUser(user, userLoginRequest.Password)
		if err != nil {
			return nil, err
		}
		refreshTokenCode, err := utils.GenerateRefreshTokenString(user)
		if err != nil {
			return nil, err
		}
		refreshTokenDaysDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS_DURATION"))
		refreshTokenDuration := time.Duration(refreshTokenDaysDuration) * 24 * time.Hour
		err = u.redisRepository.Set(c, "refresh-token:"+user.ID.String(), refreshTokenCode, refreshTokenDuration)
		if err != nil {
			return nil, err
		}
		userLoginData.RefreshToken = refreshTokenCode
		userLoginData.ExpireDuration = int(refreshTokenDuration.Seconds())
		return userLoginData, nil
	}

	user, err = u.userRepository.GetByPhoneNumber(userLoginRequest.LoginCredentials)
	if err == nil {
		userLoginData, err := utils.AuthenticateUser(user, userLoginRequest.Password)
		if err != nil {
			return nil, err
		}
		refreshTokenCode, err := utils.GenerateRefreshTokenString(user)
		if err != nil {
			return nil, err
		}
		refreshTokenDaysDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS_DURATION"))
		refreshTokenDuration := time.Duration(refreshTokenDaysDuration) * 24 * time.Hour
		err = u.redisRepository.Set(c, "refresh-token:"+user.ID.String(), refreshTokenCode, refreshTokenDuration)
		if err != nil {
			return nil, err
		}
		userLoginData.RefreshToken = refreshTokenCode
		userLoginData.ExpireDuration = int(refreshTokenDuration.Seconds())
		return userLoginData, nil
	}

	email, err := u.emailRepository.GetByEmail(userLoginRequest.LoginCredentials)
	if err == nil {
		user, err = u.userRepository.GetByEmailID(email.ID)
		if err != nil {
			return nil, err
		}
		userLoginData, err := utils.AuthenticateUser(user, userLoginRequest.Password)
		if err != nil {
			return nil, err
		}
		refreshTokenCode, err := utils.GenerateRefreshTokenString(user)
		if err != nil {
			return nil, err
		}
		refreshTokenDaysDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS_DURATION"))
		refreshTokenDuration := time.Duration(refreshTokenDaysDuration) * 24 * time.Hour
		err = u.redisRepository.Set(c, "refresh-token:"+user.ID.String(), refreshTokenCode, refreshTokenDuration)
		if err != nil {
			return nil, err
		}
		userLoginData.RefreshToken = refreshTokenCode
		userLoginData.ExpireDuration = int(refreshTokenDuration.Seconds())
		return userLoginData, nil
	}
	return nil, e.ErrUserNotFound
}
func (u *authService) RegisterWithGoogle(c *gin.Context, userRegisterWithGoogleRequest *externalDtos.UserRegisterWithGoogleRequest) (*internalDtos.UserLoginData, error) {
	_, err := u.userRepository.GetByUsername(userRegisterWithGoogleRequest.Username)
	if err != e.ErrDataNotFound {
		return nil, e.ErrUsernameAlreadyExist
	}
	_, err = u.userRepository.GetByPhoneNumber(userRegisterWithGoogleRequest.PhoneNumber)
	if err != e.ErrDataNotFound {
		return nil, e.ErrPhoneNumberAlreadyExist
	}

	email := &models.Email{
		Address: userRegisterWithGoogleRequest.Email,
	}
	_, err = u.emailRepository.Create(email)
	if err != nil {
		return nil, err
	}
	user := userRegisterWithGoogleRequest.ToModel()
	user.EmailID = email.ID

	hashedGoogleID, err := utils.Hash(*userRegisterWithGoogleRequest.GoogleID)
	if err != nil {
		return nil, err
	}
	user.GoogleID = &hashedGoogleID

	hashedPassword, err := utils.Hash(userRegisterWithGoogleRequest.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	createdUser, err := u.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	wallet := &models.Wallet{
		UserID:          createdUser.ID,
		Balance:         0,
		Pin:             "",
		IsActive:        false,
		InvalidPinCount: 0,
	}
	_, err = u.walletRepository.Create(wallet)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateUserTokenString(createdUser)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshTokenString(user)
	if err != nil {
		return nil, err
	}
	userLoginData := &internalDtos.UserLoginData{
		AccessToken:  token,
		RefreshToken: refreshToken,
		Role:         createdUser.Role,
	}
	refreshTokenDaysDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS_DURATION"))
	refreshTokenDuration := time.Duration(refreshTokenDaysDuration) * 24 * time.Hour
	err = u.redisRepository.Set(c, "refresh-token:"+user.ID.String(), token, refreshTokenDuration)
	if err != nil {
		return nil, err
	}
	userLoginData.ExpireDuration = int(refreshTokenDuration.Seconds())
	return userLoginData, nil
}

func (u *authService) LoginWithGoogle(c *gin.Context, userLoginWithGoogleRequest *externalDtos.UserLoginWithGoogleRequest) (*internalDtos.UserLoginData, error) {
	email, err := u.emailRepository.GetByEmail(userLoginWithGoogleRequest.Email)
	if err != nil {
		return nil, e.ErrRegisterGoogle
	}
	user, err := u.userRepository.GetByEmailID(email.ID)
	if err != nil {
		return nil, err
	}

	if user.GoogleID != nil {
		err = utils.CheckHash(*user.GoogleID, userLoginWithGoogleRequest.GoogleID)
		if err != nil {
			return nil, err
		}
	}

	token, err := utils.GenerateUserTokenString(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRefreshTokenString(user)
	if err != nil {
		return nil, err
	}
	userLoginData := &internalDtos.UserLoginData{
		AccessToken:  token,
		RefreshToken: refreshToken,
		Role:         user.Role,
	}
	refreshTokenDaysDuration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS_DURATION"))
	refreshTokenDuration := time.Duration(refreshTokenDaysDuration) * 24 * time.Hour
	err = u.redisRepository.Set(c, "refresh-token:"+user.ID.String(), token, refreshTokenDuration)
	if err != nil {
		return nil, err
	}
	userLoginData.ExpireDuration = int(refreshTokenDuration.Seconds())
	return userLoginData, nil
}

func (u *authService) StepUpPassword(c *gin.Context, id uuid.UUID, request *externalDtos.AuthStepUpPasswordRequest) (token string, err error) {
	user, err := u.userRepository.GetByID(id)
	if err != nil {
		return
	}
	err = utils.CheckHash(user.Password, request.Password)
	if err != nil {
		err = e.ErrInvalidPassword
		return
	}

	jwtWrapper, err := utils.NewJwtWrapper()
	if err != nil {
		return
	}
	jwtData := internalDtos.NewJwtStepUp("level1")
	token, err = jwtWrapper.JwtStepUpGenerateToken(jwtData)
	if err != nil {
		return
	}

	return
}

func (u *authService) SendResetPasswordEmail(c *gin.Context, request *externalDtos.UserVerifyResetPasswordRequest) (bool, error) {
	email, err := u.emailRepository.GetByEmail(request.Email)
	if err != nil {
		return false, e.ErrEmailNotFound
	}
	user, err := u.userRepository.GetByEmailID(email.ID)
	if err != nil {
		return false, e.ErrUserNotFound
	}
	resetPasswordID := uuid.New().String()
	err = u.redisRepository.Set(c, "reset-password:"+resetPasswordID, user.ID.String(), 15*time.Minute)
	if err != nil {
		return false, err
	}
	err = u.redisRepository.Set(c, "reset-password:"+user.ID.String(), resetPasswordID, 15*time.Minute)
	if err != nil {
		return false, err
	}

	to := []string{email.Address}
	subject := "Reset Password"
	message := utils.WriteResetPasswordEmailMessage(os.Getenv("FRONTEND_URL") + "/reset-password/" + resetPasswordID)

	err = utils.SendMail(to, subject, message)
	if err != nil {
		return false, e.ErrSendingEmailFailed
	}

	return true, nil
}

func (u *authService) ResetPassword(c *gin.Context, resetPasswordID string, request *externalDtos.ResetPasswordRequest) (*externalDtos.ResetPasswordResponse, error) {
	userID, err := u.redisRepository.Get(c, "reset-password:"+resetPasswordID)
	if err != nil {
		return nil, e.ErrInvalidResetPasswordID
	}
	getResetPasswordID, err := u.redisRepository.Get(c, "reset-password:"+userID)
	if err != nil {
		return nil, e.ErrInvalidResetPasswordID
	}
	if resetPasswordID != getResetPasswordID {
		return nil, e.ErrInvalidResetPasswordID
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepository.GetByID(userUUID)
	if err != nil {
		return nil, err
	}
	err = utils.CheckHash(user.Password, request.Password)
	if err == nil {
		return nil, e.ErrInvalidPassword
	}
	if strings.Contains(strings.ToLower(request.Password), strings.ToLower(user.Username)) {
		return nil, e.ErrInvalidPassword
	}
	hashedPassword, err := utils.Hash(request.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	_, err = u.userRepository.Update(user.ID, user)
	if err != nil {
		return nil, err
	}
	err = u.redisRepository.Delete(c, "reset-password:"+resetPasswordID)
	if err != nil {
		return nil, err
	}
	err = u.redisRepository.Delete(c, "reset-password:"+userID)
	if err != nil {
		return nil, err
	}
	responseData := externalDtos.NewResetPasswordResponse(true)
	return responseData, nil
}

func (u *authService) RefreshAccessToken(c *gin.Context, refreshToken string, userID uuid.UUID) (string, error) {
	user, err := u.userRepository.GetByID(userID)
	if err != nil {
		return "", e.ErrUserNotFound
	}
	redisRefreshToken, err := u.redisRepository.Get(c, "refresh-token:"+userID.String())
	if err != nil {
		return "", e.ErrExpiredRefreshToken
	}
	if refreshToken != redisRefreshToken {
		return "", e.ErrInvalidRefreshToken
	}
	accessToken, err := utils.GenerateUserTokenString(user)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
