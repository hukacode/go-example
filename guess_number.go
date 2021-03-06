package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	maxRange := 100
	randomSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randomSource)
	target := random.Intn(maxRange)
	maxNumberOfTries := 10
	isCorrect := false

	for maxNumberOfTries > 0 {
		fmt.Printf("Pick a number between 0 and %v. You have %d turns left.\n", maxRange, maxNumberOfTries)

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		input = strings.TrimSpace(input)
		inputInt, err := strconv.Atoi(input)

		if err != nil {
			log.Fatal(err)
		}

		if inputInt == target {
			fmt.Printf("Your guess is correct")
			isCorrect = true
			break
		} else if inputInt > target {
			fmt.Println("Your guess is greater than target")
		} else {
			fmt.Println("Your guess is lower than target")
		}

		fmt.Println()
		maxNumberOfTries--
	}

	if !isCorrect {
		fmt.Printf("Target is: %v", target)
	}
}
