package quiz

import (
	"testing"
)

type testData struct {
	number int
	expect bool
}

func TestArmstrongNumber(t *testing.T) {
	testDataList := []testData{
		{number: 1, expect: true},
		{number: 371, expect: true},
		{number: 0, expect: true},
		{number: 2, expect: false},
	}

	for _, testData := range testDataList {
		actual := IsArmstrongNumber(testData.number)
		if actual != testData.expect {
			t.Errorf("IsArmstrongNumber (%#v) = \"%v\", actual \"%v\"", testData.number, actual, testData.expect)
		}
	}
}
