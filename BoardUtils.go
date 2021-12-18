package main

import "fmt"

var posToBitBoard = PosToBitBoard()

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
			if bit == 1 {
				fmt.Print("1")
			} else {
				fmt.Print(".")
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

func PrettyMove(baseBoard uint64, move Move) {
	origin := move.Origin()
	destination := move.Destination()
	fmt.Println("------------------------")
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			pos := i*8 + j
			bit := BitAt64(&baseBoard, pos)
			fmt.Print(" ")
			if pos == origin {
				fmt.Print("O")
			} else if pos == destination {
				fmt.Print("X")
			} else if bit == 1 {
				fmt.Print("1")
			} else {
				fmt.Print(".")
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
