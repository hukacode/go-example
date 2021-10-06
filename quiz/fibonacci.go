package main

import "fmt"

func main() {
	n := 4
	recursion(n)

	fmt.Println()

	c := closure()
	for i := 0; i < n; i++ {
		fmt.Print(c())
	}
}

func recursion(n int) int {
	if n <= 1 {
		fmt.Println(n)
		return n
	}

	return recursion(n-1) + recursion(n-2)
}

func closure() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}
