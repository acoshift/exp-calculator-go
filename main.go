package main

import (
	"bufio"
	"fmt"
	"os"

	"calculator/calculator"
)

func main() {
	fmt.Print("Input: ")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()

	result := calculator.CalcString(string(line))
	fmt.Printf("Result: %f\n", result)
}
