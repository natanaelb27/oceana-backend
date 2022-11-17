package constants

var ROLE_RESOURCES_RULE map[string]map[string]string = map[string]map[string]string{
	"/roles": {
		"GET": "100",
	},
	"/profile/wallet/activate": {
		"PATCH": "011",
	},
	"/profile/register-shop": {
		"POST": "001",
	},
	"/profile/wallet/pin": {
		"PATCH": "011",
	},
	"/auth/step-up/password": {
		"POST": "011",
	},
	"/profile": {
		"PUT": "011",
		"GET": "011",
	},
	"/profile/sealabs-pay": {
		"POST": "011",
	},
	"/merchant/vouchers": {
		"POST":   "010",
		"PUT":    "010",
		"DELETE": "010",
	},
}
