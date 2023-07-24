package utils

import "fmt"

func BinaryConverter(number, bits int) []int {
	factor := number
	result := make([]int, bits)

	fmt.Printf(`number: %v`, number)
	fmt.Printf(`factor: %v`, factor)
	for factor >= 0 && number > 0 {
		factor = number % 2
		number /= 2
		fmt.Printf(`numberF: %v`, number)
		fmt.Printf(`factorF: %v`, factor)
		result[bits-1] = factor
		bits--
	}

	return result
}
