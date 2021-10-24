package main

import "fmt"

var debruijnIndices = [64]int{
	0, 1, 48, 2, 57, 49, 28, 3,
	61, 58, 50, 42, 38, 29, 17, 4,
	62, 55, 59, 36, 53, 51, 43, 22,
	45, 39, 33, 30, 24, 18, 12, 5,
	63, 47, 56, 27, 60, 41, 37, 16,
	54, 35, 52, 21, 44, 32, 23, 11,
	46, 26, 40, 15, 34, 20, 31, 10,
	25, 14, 19, 9, 13, 8, 7, 6,
}

// Gets the bit at position pos, an integer in [0, 63]
func BitAt64(i *uint64, pos int) int {
	return int((*i >> uint64(pos)) & 1)
}

// Gets the bit at position pos, an integer in [0, 63]
func BitAt32(i *uint32, pos int) int {
	return int((*i >> uint32(pos)) & 1)
}

func LSB(i uint64) int {
	if i > 0 {
		debruijn64 := uint64(0x03f79d71b4cb0a89)
		return debruijnIndices[((i&-i)*debruijn64)>>58]
	} else {
		return -1
	}
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
