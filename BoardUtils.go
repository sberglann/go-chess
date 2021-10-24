package main

import "fmt"


func IndexToCartesian(pos int) (int, int) {
	rank := pos/8 + 1
	file := pos%8 + 1

	return rank, file
}

func CartesianToIndex(rank int, file int) int {
	return (rank-1)*8 + file - 1
}

func OppositeColor(color Color) Color {
	switch color {
	case White:
		return Black
	case Black:
		return White
	default:
		return Blank
	}
}

func Inverse(i uint64) uint64 {
	return i ^ 0xFFFFFFFFFFFFFFFF
}

func Pretty64(number uint64) {
	fmt.Println("------------------------")
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			pos := i*8 + j
			bit := BitAt64(&number, pos)
			fmt.Print(" ")
			fmt.Print(bit)
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
