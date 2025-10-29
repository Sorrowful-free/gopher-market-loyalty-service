package utils

import "strconv"

func ValidateLuhn(number int) bool {

	numberString := strconv.Itoa(number)

	sum := 0
	length := len(numberString)
	for i := 0; i < length; i++ {
		digit := int(numberString[i] - '0')
		if (length-i)%2 == 0 {
			digit = digit * 2
			if digit > 9 {
				digit = digit - 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
