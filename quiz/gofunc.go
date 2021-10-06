package main

import "fmt"

func sayHi() {
	fmt.Println("Hi")
}
func sayDear(name string) {
	fmt.Printf("Dear %v\n", name)
}
func add(a int, b int) int {
	return a + b
}
func minus(a int, b int) int {
	return a - b
}
func doSomething(function func(int, int) int) {
	result := function(1, 2)
	fmt.Println(result)
}
func main() {
	var hi func()
	hi = sayHi
	hi() // chú ý dấu ()

	hi2 := sayHi
	hi2()

	var dear func(string)
	dear = sayDear
	dear("World")

	doSomething(add)
	doSomething(minus)
}
