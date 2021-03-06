package main

import "fmt"

func a(myChan chan string) {
	myChan <- "hello from chan"
}

func main() {
	myChan := make(chan string)
	go a(myChan)
	output := <-myChan
	fmt.Print(output)
}
