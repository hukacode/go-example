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
	myChan1 := make(chan string)
	myChan2 := make(chan string)
	go a(myChan1)
	go b(myChan2)
	for i := 0; i < 50; i++ {
		fmt.Print(<-myChan1)
		fmt.Print(<-myChan2)
	}
	fmt.Println()
	fmt.Println("end main()")
}
