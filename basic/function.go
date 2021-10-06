package main

import "fmt"

func main() {
	xs := []float64{1, 2, 4}
	fmt.Println(avg(xs))
	fmt.Println(sum("x", 1, 2, 4))

	fmt.Println(half(1))
	fmt.Println(half(2))

	nextEven := makeEvenGenerator()
	fmt.Println(nextEven()) // 0
	fmt.Println(nextEven()) // 2
	fmt.Println(nextEven()) // 4
}

func avg(xs []float64) float64 {
	total := 0.0

	for _, v := range xs {
		total += v
	}

	return total / float64(len(xs))
}

func sum(x string, args ...int) int {
	fmt.Println(x)
	total := 0
	for _, v := range args {
		total += v
	}
	return total
}

func makeEvenGenerator() func() uint {
	i := uint(0)
	return func() uint {
		i += 2
		return i
	}
}

func half(i int) (int, bool) {
	return i / 2, i%2 == 0
}
