package utils

func ValidateLuhn(number string) bool {
	sum := 0
	length := len(number)
	for i := 0; i < length; i++ {
		digit := int(number[i] - '0')
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
