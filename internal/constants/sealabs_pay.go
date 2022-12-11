package constants

type SealabsPayAction int

const (
	SealabsPayRegister SealabsPayAction = 0
	SealabsPayPayment  SealabsPayAction = 1
	SealabsPayTopup    SealabsPayAction = 2
)

const (
	SealabsPayCardNumberLength = 16
)
