package main

import "fmt"

// Gets the bit at position pos, an integer in [0, 63]
func BitAt64(i *uint64, pos int) int {
	return int((*i >> uint64(pos)) & 1)
}

// Gets the bit at position pos, an integer in [0, 63]
func BitAt32(i *uint32, pos int) int {
	return int((*i >> uint32(pos)) & 1)
}

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
