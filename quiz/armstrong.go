package quiz

/* */
func IsArmstrongNumber(input int) bool {
	output := 0
	temp := input

	for temp > 0 {
		c := temp % 10
		output += c * c * c
		temp = temp / 10
	}

	return input == output
}
