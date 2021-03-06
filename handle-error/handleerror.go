package main

import "fmt"

func reportPanic() {
	p := recover()

	if p == nil {
		return
	}

	err, ok := p.(error)

	if ok {
		fmt.Println(err)
	} else {
		panic(p)
	}
}

func panicTest() {
	defer reportPanic()
	panic("I'm a panic")
	fmt.Println("I won't be run")
}

func main() {
	panicTest()
	fmt.Println("I can be run")
}
