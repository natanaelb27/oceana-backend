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
	ErrEmptyRequstBody     = errors.New("request body should not be empty")
	ErrRequestBodyTooLarge = errors.New("request body must not be larger than 1MB")
	ErrEmptyVariants       = errors.New("variants cannot be empty")

	ErrInvalidEmailAddress            = errors.New("invalid email address")
	ErrInvalidPassword                = errors.New("invalid password")
	ErrInvalidAmount                  = errors.New("invalid amount")
	ErrInvalidOrderDetailID           = errors.New("invalid order detail id")
	ErrInvalidMembershipID            = errors.New("invalid membership id")
	ErrInvalidName                    = errors.New("invalid name")
	ErrInvalidImage                   = errors.New("invalid image")
	ErrInvalidDescription             = errors.New("invalid description")
	ErrInvalidStartTime               = errors.New("invalid start time")
	ErrInvalidDurationMinute          = errors.New("invalid duration minute")
	ErrInvalidCapacity                = errors.New("invalid capacity")
	ErrInvalidStock                   = errors.New("invalid stock")
	ErrInvalidCategoryID              = errors.New("invalid category id")
	ErrInvalidMerchandiseID           = errors.New("invalid merchandise id")
	ErrInvalidQuantity                = errors.New("invalid quantity")
	ErrInvalidPrice                   = errors.New("invalid price")
	ErrInvalidEventID                 = errors.New("invalid event id")
	ErrInvalidSignature               = errors.New("invalid signature")
	ErrDataNotFound                   = errors.New("data not found")
	ErrMerchantNotFound               = errors.New("merchant not found")
	ErrEmailNotFound                  = errors.New("email not found")
	ErrInvalidRequest                 = errors.New("invalid request")
	ErrInvalidIdentifier              = errors.New("invalid identifier")
	ErrInvalidToken                   = errors.New("invalid token")
	ErrPasswordCannotContainsUsername = errors.New("password cannot contains username")
	ErrInvalidVerificationCode        = errors.New("invalid verification code")
	ErrInvalidSealabsPayID            = errors.New("invalid sealabs pay id")
	ErrInvalidSealabsPayCardNumber    = errors.New("invalid sealabs pay card number")
	ErrInvalidMerchantVoucherCode     = errors.New("invalid voucher code")
	ErrInvalidParams                  = errors.New("invalid params")
	ErrInvalidResetPasswordID         = errors.New("invalid reset password id")
	ErrInvalidModel                   = errors.New("invalid model")
	ErrInvalidVariantsMedia           = errors.New("invalid variants media")
	ErrInvalidDiscountType            = errors.New("invalid discount type")
	ErrInvalidDiscount                = errors.New("invalid discount")

	ErrWrongEmailAddress                         = errors.New("wrong email address")
	ErrWrongPassword                             = errors.New("wrong password")
	ErrEmailAddressAlreadyExist                  = errors.New("email address is already registered")
	ErrUsernameAlreadyExist                      = errors.New("username is already registered")
	ErrPhoneNumberAlreadyExist                   = errors.New("phone number is already registered")
	ErrGoogleIDAlreadyExist                      = errors.New("google id is already registered")
	ErrWalletAlreadyActive                       = errors.New("wallet is already active")
	ErrPasswordCannotBeTheSameAsPreviousPassword = errors.New("password cannot be the same as previous password")
	ErrSealabsPayCardNumberAlreadyExist          = errors.New("account number is already registered")
	ErrMerchantVoucherCodeAlreadyExist           = errors.New("voucher code already exist")
	ErrMerchantNameAlreadyExist                  = errors.New("merchant name is already registered")
	ErrMerchantPromotionProductConflict          = errors.New("one or more products have conflicting promotion periods")
	ErrDeleteShopAddressConflict                 = errors.New("cannot delete shop address")
	ErrDeleteDefaultAddressConflict              = errors.New("cannot delete default address")
	ErrWrongWalletPin                            = errors.New("wrong wallet pin")
	ErrRefundAlreadyExist                        = errors.New("refund already exist")

	ErrUnauthorized                     = errors.New("access is denied due to invalid credentials")
	ErrUserNotFound                     = errors.New("user not found")
	ErrForbidden                        = errors.New("you don't have permission to access this resource")
	ErrResourceDoesNotExist             = errors.New("resource does not exist")
	ErrStepUpNotFound                   = errors.New("step up authenticaiton not found")
	ErrStepUpScopeIncorrect             = errors.New("incorrect step up scope")
	ErrVerificationCodeNotFound         = errors.New("verification code not found")
	ErrProductNotFound                  = errors.New("product not found")
	ErrVoucherNotFound                  = errors.New("voucher not found")
	ErrPromotionNotFound                = errors.New("promotion not found")
	ErrProductMerchantPromotionNotFound = errors.New("product merchant promotion not found")
	ErrOrderNotFound                    = errors.New("order not found")

	ErrForbiddenMerchant = errors.New("user is not the merchant")
	ErrForbiddenOrder    = errors.New("user is not the buyer")

	ErrInsufficientBalance                                        = errors.New("insufficient balance")
	ErrEventAlreadyBookmarked                                     = errors.New("event is already bookmarked")
	ErrEventAlreadyPaid                                           = errors.New("event is already paid")
	ErrMerchandiseAlreadyPaid                                     = errors.New("merchandise is already paid")
	ErrEventHasNotBeenPaid                                        = errors.New("event has not been paid")
	ErrMerchandiseAlreadyInCart                                   = errors.New("merchandise already in cart")
	ErrExceedEventCapacity                                        = errors.New("event has reached maximum capacity")
	ErrExceedMerchandiseStock                                     = errors.New("merchandise has no stock left")
	ErrWalletIsNotActive                                          = errors.New("wallet is not active")
	ErrDiscountExceedMinimumSpendAmount                           = errors.New("discount is exceeding the minimum spend amount")
	ErrExceedDiscountPercentage                                   = errors.New("discount is more than 100%")
	ErrMaximumDiscountAmountExceedMinimumSpendAmount              = errors.New("maximum discount amount is exceeding the minimum spend amount")
	ErrVoucherStartDatetimeEarlierThanNow                         = errors.New("voucher start period is earlier than now")
	ErrVoucherEndDatetimeEarlierThanVoucherStartDatetime          = errors.New("voucher end period is earlier than the start period")
	ErrVoucherPeriodMinimumOneHour                                = errors.New("voucher period must be at least 1 hour")
	ErrVoucherStartDatetimeMinimumOneDayFromNow                   = errors.New("voucher start period must be at least 1 day from now")
	ErrVoucherHasEnded                                            = errors.New("voucher has ended")
	ErrOngoingStartDatetimeUpdated                                = errors.New("ongoing voucher/promotion start period cannot be updated")
	ErrMaximumQuotaLowerThanCurrentQuota                          = errors.New("voucher/promotion maximum quota is lower than current quota")
	ErrMaximumQuotaLowerThanMaximumUserLimit                      = errors.New("promotion maximum quota is than maximum user limit")
	ErrVoucherEndDatetimeEarlierThanNow                           = errors.New("voucher end period is earlier than now")
	ErrVoucherNotUpcoming                                         = errors.New("voucher is not upcoming")
	ErrStockIsNotAvailable                                        = errors.New("stock is not available")
	ErrStatusIsNotPendingSeller                                   = errors.New("status is not pending seller")
	ErrStatusIsNotOnProcess                                       = errors.New("status is not on process")
	ErrStatusIsNotDelivered                                       = errors.New("status is not delivered")
	ErrUpdatedVoucherStartDatetimeEarlierThanCurrentStartDatetime = errors.New("new voucher start period is earlier than the current start period")
	ErrUpdatedVoucherEndDatetimeExceedCurrentEndDatetime          = errors.New("new voucher end period is exceeding the current end period")
	ErrVoucherPeriodExceedMaximumDays                             = errors.New("voucher period must be less than 180 days")
	ErrBlockedWallet                                              = errors.New("wallet is currently blocked")

	ErrPromotionStartDatetimeEarlierThanNow                = errors.New("promotion start period is earlier than now")
	ErrPromotionEndDatetimeEarlierThanVoucherStartDatetime = errors.New("promotion end period is earlier than the start period")
	ErrPromotionPeriodMinimumOneHour                       = errors.New("promotion period must be at least 1 hour")
	ErrPromotionStartDatetimeMinimumOneDayFromNow          = errors.New("promotion start period must be at least 1 day from now")
	ErrPromotionPeriodExceedMaximumDays                    = errors.New("promotion period must be less than 180 days")
	ErrPromotionIsOngoing                                  = errors.New("promotion is ongoing")
	ErrPromotionHasEnded                                   = errors.New("promotion has ended")
	ErrPromotionDatetimeUpdatedBeyondRange                 = errors.New("updated promotion start/end period must be between current start/end period")
	ErrPromotionUpdateProductIDListMismatch                = errors.New("updated promotion product id list and product merchant promotion are mismatched")

	ErrInsertFailed                   = errors.New("insert failed")
	ErrDatabaseOperationFailed        = errors.New("database operation failed")
	ErrCacheOperationFailed           = errors.New("cache operation failed")
	ErrAuthenticationOperationFailed  = errors.New("authentication operation failed")
	ErrHashingFailed                  = errors.New("hashing failed")
	ErrDatabaseConnectionFailed       = errors.New("database connection failed")
	ErrSendingEmailFailed             = errors.New("sending email failed")
	ErrInternalOperationFailed        = errors.New("internal operation failed")
	ErrSealabsPayOperationFailed      = errors.New("sealabs pay operation failed")
	ErrInvalidSealabsPayAction        = errors.New("invalid sealabs pay action")
	ErrDatabaseTransactionSetUpFailed = errors.New("database transactions set up failed")

	ErrUnknown        = errors.New("unknown")
	ErrRegisterGoogle = errors.New("register with google")

	// bad request
	ErrRequestedVoucherNotMerchant                             = errors.New("the requested voucher is not the merchant's voucher")
	ErrRequestedPromotionNotProduct                            = errors.New("product does not have that promotion")
	ErrMerchantDoesNotHaveThatProductItem                      = errors.New("merchant does not have that product item")
	ErrProductDoesNotHaveThatProductItem                       = errors.New("product does not have that product item")
	ErrOrderDoesNotHaveThatOrderDetail                         = errors.New("order does not have that order detail")
	ErrPurchasedAmountNotExceedMerchantVoucherMinimumAmount    = errors.New("purchased amount does not exceed the merchant voucher's minimum purchased amount")
	ErrPurchasedAmountNotExceedMarketplaceVoucherMinimumAmount = errors.New("purchased amount does not exceed the marketplace voucher's minimum purchased amount")

	// 422
	ErrMarketplaceVoucherNotOngoing                           = errors.New("marketplace voucher is not ongoing")
	ErrMarketplaceVoucherQuotaFull                            = errors.New("marketplace voucher's quota is full")
	ErrMerchantVoucherNotOngoing                              = errors.New("merchant voucher is not ongoing")
	ErrMerchantVoucherQuotaFull                               = errors.New("merchant voucher's quota is full")
	ErrMerchantVoucherUsedByUser                              = errors.New("merchant voucher has been used by the user")
	ErrMerchantPromotionNotOngoing                            = errors.New("merchant promotion is not ongoing")
	ErrProductStockNotEnough                                  = errors.New("product stock is not enough")
	ErrPurchasedQuantityExceedRemainingMerchantPromotionQuota = errors.New("the amount purchased exceeds the remaining merchant promotion quota")
	ErrPurchasedQuantityExceedRemainingUserLimit              = errors.New("the amount purchased exceeds the remaining user's limit")
	ErrDifferentTotalPayment                                  = errors.New("there are changes in some promotion or discount")
	ErrWalletBalanceNotEnough                                 = errors.New("wallet balance is not enough")
	ErrProductStatusInactive                                  = errors.New("product is delisted")
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
		errors.Is(err, ErrEmptyRequstBody) ||
		errors.Is(err, ErrEmptyVariants) ||
		errors.Is(err, ErrRequestBodyTooLarge) ||
		errors.Is(err, ErrInvalidEmailAddress) ||
		errors.Is(err, ErrInvalidPassword) ||
		errors.Is(err, ErrInvalidAmount) ||
		errors.Is(err, ErrInvalidOrderDetailID) ||
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
		errors.Is(err, ErrInvalidSignature) ||
		errors.Is(err, ErrInvalidIdentifier) ||
		errors.Is(err, ErrInvalidToken) ||
		errors.Is(err, ErrPasswordCannotContainsUsername) ||
		errors.Is(err, ErrInvalidRequest) ||
		errors.Is(err, ErrInvalidVerificationCode) ||
		errors.Is(err, ErrInvalidRequest) ||
		errors.Is(err, ErrInvalidSealabsPayID) ||
		errors.Is(err, ErrInvalidSealabsPayCardNumber) ||
		errors.Is(err, ErrInvalidMerchantVoucherCode) ||
		errors.Is(err, ErrInvalidParams) ||
		errors.Is(err, ErrInvalidResetPasswordID) ||
		errors.Is(err, ErrInvalidModel) ||
		errors.Is(err, ErrInvalidVariantsMedia) ||
		errors.Is(err, ErrInvalidDiscountType) ||
		errors.Is(err, ErrInvalidDiscount) ||
		errors.Is(err, ErrRequestedVoucherNotMerchant) ||
		errors.Is(err, ErrRequestedPromotionNotProduct) ||
		errors.Is(err, ErrMerchantDoesNotHaveThatProductItem) ||
		errors.Is(err, ErrProductDoesNotHaveThatProductItem) ||
		errors.Is(err, ErrPurchasedAmountNotExceedMerchantVoucherMinimumAmount) ||
		errors.Is(err, ErrPurchasedAmountNotExceedMarketplaceVoucherMinimumAmount) ||
		errors.Is(err, ErrOrderDoesNotHaveThatOrderDetail) {
		return http.StatusBadRequest
	} else if errors.Is(err, ErrWrongEmailAddress) ||
		errors.Is(err, ErrWrongPassword) ||
		errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrStepUpNotFound) ||
		errors.Is(err, ErrStepUpScopeIncorrect) ||
		errors.Is(err, ErrVerificationCodeNotFound) ||
		errors.Is(err, ErrWrongWalletPin) {
		return http.StatusUnauthorized
	} else if errors.Is(err, ErrDataNotFound) ||
		errors.Is(err, ErrMerchantNotFound) ||
		errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrEmailNotFound) ||
		errors.Is(err, ErrVoucherNotFound) ||
		errors.Is(err, ErrPromotionNotFound) ||
		errors.Is(err, ErrProductMerchantPromotionNotFound) ||
		errors.Is(err, ErrOrderNotFound) ||
		errors.Is(err, ErrProductNotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, ErrEmailAddressAlreadyExist) ||
		errors.Is(err, ErrWalletAlreadyActive) ||
		errors.Is(err, ErrUsernameAlreadyExist) ||
		errors.Is(err, ErrPasswordCannotBeTheSameAsPreviousPassword) ||
		errors.Is(err, ErrPhoneNumberAlreadyExist) ||
		errors.Is(err, ErrGoogleIDAlreadyExist) ||
		errors.Is(err, ErrMerchantNameAlreadyExist) ||
		errors.Is(err, ErrSealabsPayCardNumberAlreadyExist) ||
		errors.Is(err, ErrMerchantVoucherCodeAlreadyExist) ||
		errors.Is(err, ErrMerchantPromotionProductConflict) ||
		errors.Is(err, ErrDeleteDefaultAddressConflict) ||
		errors.Is(err, ErrDeleteShopAddressConflict) ||
		errors.Is(err, ErrRefundAlreadyExist) {
		return http.StatusConflict
	} else if errors.Is(err, ErrForbidden) ||
		errors.Is(err, ErrForbiddenMerchant) ||
		errors.Is(err, ErrForbiddenOrder) {
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
		errors.Is(err, ErrExceedMerchandiseStock) ||
		errors.Is(err, ErrWalletIsNotActive) ||
		errors.Is(err, ErrDiscountExceedMinimumSpendAmount) ||
		errors.Is(err, ErrExceedDiscountPercentage) ||
		errors.Is(err, ErrMaximumDiscountAmountExceedMinimumSpendAmount) ||
		errors.Is(err, ErrVoucherStartDatetimeEarlierThanNow) ||
		errors.Is(err, ErrVoucherEndDatetimeEarlierThanVoucherStartDatetime) ||
		errors.Is(err, ErrVoucherPeriodMinimumOneHour) ||
		errors.Is(err, ErrVoucherStartDatetimeMinimumOneDayFromNow) ||
		errors.Is(err, ErrVoucherHasEnded) ||
		errors.Is(err, ErrOngoingStartDatetimeUpdated) ||
		errors.Is(err, ErrMaximumQuotaLowerThanCurrentQuota) ||
		errors.Is(err, ErrMaximumQuotaLowerThanMaximumUserLimit) ||
		errors.Is(err, ErrVoucherEndDatetimeEarlierThanNow) ||
		errors.Is(err, ErrVoucherNotUpcoming) ||
		errors.Is(err, ErrPromotionStartDatetimeEarlierThanNow) ||
		errors.Is(err, ErrPromotionEndDatetimeEarlierThanVoucherStartDatetime) ||
		errors.Is(err, ErrPromotionPeriodMinimumOneHour) ||
		errors.Is(err, ErrPromotionStartDatetimeMinimumOneDayFromNow) ||
		errors.Is(err, ErrPromotionStartDatetimeMinimumOneDayFromNow) ||
		errors.Is(err, ErrPromotionPeriodExceedMaximumDays) ||
		errors.Is(err, ErrPromotionIsOngoing) ||
		errors.Is(err, ErrPromotionHasEnded) ||
		errors.Is(err, ErrPromotionDatetimeUpdatedBeyondRange) ||
		errors.Is(err, ErrPromotionUpdateProductIDListMismatch) ||
		errors.Is(err, ErrStockIsNotAvailable) ||
		errors.Is(err, ErrStatusIsNotPendingSeller) ||
		errors.Is(err, ErrStatusIsNotOnProcess) ||
		errors.Is(err, ErrStatusIsNotDelivered) ||
		errors.Is(err, ErrVoucherNotUpcoming) ||
		errors.Is(err, ErrVoucherPeriodExceedMaximumDays) ||
		errors.Is(err, ErrBlockedWallet) ||
		errors.Is(err, ErrUpdatedVoucherEndDatetimeExceedCurrentEndDatetime) ||
		errors.Is(err, ErrMarketplaceVoucherNotOngoing) ||
		errors.Is(err, ErrMarketplaceVoucherQuotaFull) ||
		errors.Is(err, ErrMerchantVoucherNotOngoing) ||
		errors.Is(err, ErrMerchantVoucherQuotaFull) ||
		errors.Is(err, ErrMerchantVoucherUsedByUser) ||
		errors.Is(err, ErrMerchantPromotionNotOngoing) ||
		errors.Is(err, ErrProductStockNotEnough) ||
		errors.Is(err, ErrPurchasedQuantityExceedRemainingMerchantPromotionQuota) ||
		errors.Is(err, ErrPurchasedQuantityExceedRemainingUserLimit) ||
		errors.Is(err, ErrDifferentTotalPayment) ||
		errors.Is(err, ErrWalletBalanceNotEnough) ||
		errors.Is(err, ErrProductStatusInactive) {
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
