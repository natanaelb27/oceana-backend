package errors

import (
	"errors"
	"net/http"
)

var (
	ErrEmptyEmailAddress   = errors.New("email address should not be empty")
	ErrEmptyPassword       = errors.New("password should not be empty")
	ErrEmptyAmount         = errors.New("amount should not be empty")
	ErrEmptyMembershipID   = errors.New("membership id should not be empty")
	ErrEmptyName           = errors.New("name should not be empty")
	ErrEmptyImage          = errors.New("image should not be empty")
	ErrEmptyDescription    = errors.New("description should not be empty")
	ErrEmptyStartTime      = errors.New("start time should not be empty")
	ErrEmptyDurationMinute = errors.New("duration minute should not be empty")
	ErrEmptyCapacity       = errors.New("capacity should not be empty")
	ErrEmptyStock          = errors.New("stock should not be empty")
	ErrEmptyCategoryID     = errors.New("category id should not be empty")
	ErrEmptyMerchandiseID  = errors.New("merchandise id should not be empty")
	ErrEmptyQuantity       = errors.New("quantity should not be empty")
	ErrEmptyPrice          = errors.New("price should not be empty")
	ErrEmptyEventID        = errors.New("event id should not be empty")

	ErrInvalidEmailAddress   = errors.New("invalid email address")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidAmount         = errors.New("invalid amount")
	ErrInvalidMembershipID   = errors.New("invalid membership id")
	ErrInvalidName           = errors.New("invalid name")
	ErrInvalidImage          = errors.New("invalid image")
	ErrInvalidDescription    = errors.New("invalid description")
	ErrInvalidStartTime      = errors.New("invalid start time")
	ErrInvalidDurationMinute = errors.New("invalid duration minute")
	ErrInvalidCapacity       = errors.New("invalid capacity")
	ErrInvalidStock          = errors.New("invalid stock")
	ErrInvalidCategoryID     = errors.New("invalid category id")
	ErrInvalidMerchandiseID  = errors.New("invalid merchandise id")
	ErrInvalidQuantity       = errors.New("invalid quantity")
	ErrInvalidPrice          = errors.New("invalid price")
	ErrInvalidEventID        = errors.New("invalid event id")
	ErrDataNotFound          = errors.New("data not found")
	ErrInvalidRequest        = errors.New("invalid request")
	ErrInvalidIdentifier     = errors.New("invalid identifier")
	ErrInvalidToken          = errors.New("invalid token")

	ErrWrongEmailAddress        = errors.New("wrong email address")
	ErrWrongPassword            = errors.New("wrong password")
	ErrEmailAddressAlreadyExist = errors.New("email address is already registered")

	ErrUnauthorized         = errors.New("access is denied due to invalid credentials")
	ErrForbidden            = errors.New("you don't have permission to access this resource")
	ErrResourceDoesNotExist = errors.New("resource does not exist")

	ErrInsufficientBalance      = errors.New("insufficient balance")
	ErrEventAlreadyBookmarked   = errors.New("event is already bookmarked")
	ErrEventAlreadyPaid         = errors.New("event is already paid")
	ErrMerchandiseAlreadyPaid   = errors.New("merchandise is already paid")
	ErrEventHasNotBeenPaid      = errors.New("event has not been paid")
	ErrMerchandiseAlreadyInCart = errors.New("merchandise already in cart")
	ErrExceedEventCapacity      = errors.New("event has reached maximum capacity")
	ErrExceedMerchandiseStock   = errors.New("merchandise has no stock left")

	ErrInsertFailed                  = errors.New("insert failed")
	ErrDatabaseOperationFailed       = errors.New("database operation failed")
	ErrAuthenticationOperationFailed = errors.New("authentication operation failed")

	ErrUnknown = errors.New("unknown")
)

func GetStatusCode(err error) int {
	if errors.Is(err, ErrEmptyEmailAddress) ||
		errors.Is(err, ErrEmptyPassword) ||
		errors.Is(err, ErrEmptyAmount) ||
		errors.Is(err, ErrEmptyMembershipID) ||
		errors.Is(err, ErrEmptyName) ||
		errors.Is(err, ErrEmptyImage) ||
		errors.Is(err, ErrEmptyDescription) ||
		errors.Is(err, ErrEmptyStartTime) ||
		errors.Is(err, ErrEmptyDurationMinute) ||
		errors.Is(err, ErrEmptyCapacity) ||
		errors.Is(err, ErrEmptyStock) ||
		errors.Is(err, ErrEmptyCategoryID) ||
		errors.Is(err, ErrEmptyMerchandiseID) ||
		errors.Is(err, ErrEmptyQuantity) ||
		errors.Is(err, ErrEmptyPrice) ||
		errors.Is(err, ErrEmptyEventID) ||
		errors.Is(err, ErrInvalidEmailAddress) ||
		errors.Is(err, ErrInvalidPassword) ||
		errors.Is(err, ErrInvalidAmount) ||
		errors.Is(err, ErrInvalidMembershipID) ||
		errors.Is(err, ErrInvalidName) ||
		errors.Is(err, ErrInvalidImage) ||
		errors.Is(err, ErrInvalidDescription) ||
		errors.Is(err, ErrInvalidStartTime) ||
		errors.Is(err, ErrInvalidDurationMinute) ||
		errors.Is(err, ErrInvalidCapacity) ||
		errors.Is(err, ErrInvalidStock) ||
		errors.Is(err, ErrInvalidCategoryID) ||
		errors.Is(err, ErrInvalidMerchandiseID) ||
		errors.Is(err, ErrInvalidQuantity) ||
		errors.Is(err, ErrInvalidPrice) ||
		errors.Is(err, ErrInvalidEventID) ||
		errors.Is(err, ErrInvalidIdentifier) ||
		errors.Is(err, ErrDataNotFound) ||
		errors.Is(err, ErrInvalidToken) ||
		errors.Is(err, ErrInvalidRequest) {
		return http.StatusBadRequest
	} else if errors.Is(err, ErrWrongEmailAddress) ||
		errors.Is(err, ErrWrongPassword) ||
		errors.Is(err, ErrUnauthorized) {
		return http.StatusUnauthorized
	} else if errors.Is(err, ErrEmailAddressAlreadyExist) {
		return http.StatusConflict
	} else if errors.Is(err, ErrForbidden) {
		return http.StatusForbidden
	} else if errors.Is(err, ErrResourceDoesNotExist) {
		return http.StatusNotFound
	} else if errors.Is(err, ErrInsufficientBalance) ||
		errors.Is(err, ErrEventAlreadyBookmarked) ||
		errors.Is(err, ErrEventAlreadyPaid) ||
		errors.Is(err, ErrMerchandiseAlreadyPaid) ||
		errors.Is(err, ErrEventHasNotBeenPaid) ||
		errors.Is(err, ErrMerchandiseAlreadyInCart) ||
		errors.Is(err, ErrExceedEventCapacity) ||
		errors.Is(err, ErrExceedMerchandiseStock) {
		return http.StatusUnprocessableEntity
	}

	return http.StatusInternalServerError
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func IsNot(err, target error) bool {
	return !errors.Is(err, target)
}
