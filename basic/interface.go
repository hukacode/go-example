package main

import (
	"fmt"
	"math"
)

type Shape interface {
	area() float64
}

type Square struct {
	a float64
}

type Circle struct {
	r float64
}

func (s Square) area() float64 {
	return s.a * s.a
}

func (c Circle) area() float64 {
	return math.Pi * c.r * c.r
}

func main() {
	var square Shape = Square{2}
	fmt.Println(square.area())

	var circle Shape = Square{2}
	fmt.Println(circle.area())
}
