package phoneNumber

import (
	"fmt"
	"strconv"
)

func IsValid(phoneNumber string) bool {
	fmt.Println("len checking...")
	if len(phoneNumber) != 11 {
		return false
	}

	//TODO : supporting +98 in phone number
	fmt.Println("09 checking...")
	if phoneNumber[0:2] != "09" {
		return false
	}

	fmt.Println("conv checking...")
	if _, err := strconv.Atoi(phoneNumber[2:]); err != nil {
		return false
	}

	fmt.Println("end of checking...")
	return true
}
