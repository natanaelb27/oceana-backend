package constants

var ROLE_RESOURCES_RULE map[string]map[string]string = map[string]map[string]string{
	"/roles": {
		"GET": "100",
	},
	"/profile/wallet/activate": {
		"PATCH": "011",
	},
	"/profile": {
		"PUT": "011",
		"GET": "011",
	},
	"/my-shop/vouchers": {
		"POST": "010",
	},
}
