package phoneNumber

import (
	"strconv"
)

func IsValid(phoneNumber string) bool {
	if len(phoneNumber) != 11 {
		return false
	}

	//TODO : supporting +98 in phone number
	if phoneNumber[0:2] != "09" {
		return false
	}

	if _, err := strconv.Atoi(phoneNumber[2:]); err != nil {
		return false
	}

	return true
}
