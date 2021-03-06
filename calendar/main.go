package main

import (
	"fmt"

	"github.com/hukacode/go-example/calendar"
)

func main() {
	event := calendar.Event{}
	event.Year = 1
	// event.Month = 2 // ambiguous selector event.Month
	event.Date2.Month = 2
	event.Date.Year = 3
	// event.Day2 = 4 // error
	fmt.Println(event)
}
