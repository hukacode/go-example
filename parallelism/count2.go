package main

import "fmt"

func a(myChan chan string) {
	for i := 0; i < 50; i++ {
		// fmt.Print("a")
		myChan <- "a"
	}
}

func b(myChan chan string) {
	for i := 0; i < 50; i++ {
		// fmt.Print("b")
		myChan <- "b"
	}
}

func main() {
	myChan := make(chan string)
	go a(myChan)
	go b(myChan)
	for i := 0; i < 100; i++ {
		fmt.Print(<-myChan)
	}
	fmt.Println()
	fmt.Println("end main()")
}
