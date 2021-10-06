package main

import "fmt"

func main() {
	defer func() {
		str := recover()
		fmt.Println(str)
	}()
	panic("PANIC")
	fmt.Print("Hello")
}
