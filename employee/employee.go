package main

import (
	"fmt"
)

type Employee struct {
	Name   string
	Age    int
	Salary int
}

func (e *Employee) doubleSalary() {
	e.Salary *= 2
}

func main() {
	employee := Employee{Name: "Huka", Age: 10, Salary: 1000}
	fmt.Println(employee)
	employee.doubleSalary()
	fmt.Println(employee)
}
