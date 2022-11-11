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

	ErrWrongEmailAddress        = errors.New("wrong email address")
	ErrWrongPassword            = errors.New("wrong password")
	ErrEmailAddressAlreadyExist = errors.New("email address is already registered")

	ErrUnauthorized  = errors.New("access is denied due to invalid credentials")
	ErrForbidden     = errors.New("you don't have permission to access this resource")
	ErrRouteNotFound = errors.New("route does not exist")

	ErrInsufficientBalance      = errors.New("insufficient balance")
	ErrEventAlreadyBookmarked   = errors.New("event is already bookmarked")
	ErrEventAlreadyPaid         = errors.New("event is already paid")
	ErrMerchandiseAlreadyPaid   = errors.New("merchandise is already paid")
	ErrEventHasNotBeenPaid      = errors.New("event has not been paid")
	ErrMerchandiseAlreadyInCart = errors.New("merchandise already in cart")
	ErrExceedEventCapacity      = errors.New("event has reached maximum capacity")
	ErrExceedMerchandiseStock   = errors.New("merchandise has no stock left")

	ErrUnknown = errors.New("unknown")
)

func GetStatusCode(err error) int {
	if errors.Is(err, ErrEmptyEmailAddress) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyPassword) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyAmount) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyMembershipID) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyName) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyImage) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyDescription) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyStartTime) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyDurationMinute) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyCapacity) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyStock) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyCategoryID) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyMerchandiseID) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyQuantity) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyPrice) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrEmptyEventID) {
		return http.StatusBadRequest
	}

	if errors.Is(err, ErrInvalidEmailAddress) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidPassword) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidAmount) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidMembershipID) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidName) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidImage) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidDescription) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidStartTime) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidDurationMinute) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidCapacity) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidStock) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidCategoryID) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidMerchandiseID) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidQuantity) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidPrice) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidEventID) {
		return http.StatusBadRequest
	}

	if errors.Is(err, ErrWrongEmailAddress) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, ErrWrongPassword) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, ErrEmailAddressAlreadyExist) {
		return http.StatusConflict
	}

	if errors.Is(err, ErrUnauthorized) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, ErrForbidden) {
		return http.StatusForbidden
	}

	if errors.Is(err, ErrRouteNotFound) {
		return http.StatusNotFound
	}

	if errors.Is(err, ErrInsufficientBalance) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrEventAlreadyBookmarked) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrEventAlreadyPaid) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrMerchandiseAlreadyPaid) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrEventHasNotBeenPaid) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrMerchandiseAlreadyInCart) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrExceedEventCapacity) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrExceedMerchandiseStock) {
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
