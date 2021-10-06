package iteration

func Repeat(char string) string {
	var repeated string
	for i := 0; i < 4; i++ {
		repeated += char
	}
	return repeated
}
