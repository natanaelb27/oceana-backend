package constants

var ROLE_RESOURCES_RULE map[string]map[string]string = map[string]map[string]string{
	"/roles": {
		"GET": "100",
	},
	"/profile/wallet/activate": {
		"PATCH": "011",
	},
	"/profile/addresses": {
		"POST": "001",
		"GET":  "001",
	},
}
