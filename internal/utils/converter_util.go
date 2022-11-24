package utils

import "strconv"

// Convert string to uint
// If err when converting then set target to 0
func StringToUint(source string) (target uint) {
	result, err := strconv.Atoi(source)
	if err != nil {
		result = 0
	} else if result <= 0 {
		target = 0
	} else {
		target = uint(result)
	}
	return
}
