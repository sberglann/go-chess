package main

import "fmt"

func main() {
	b := StartBoard
	b.PrettyBoard()
	ms := RookMoves(b)
	for _, m := range ms {
		Transition(b, m, ColoredPiece{piece: Rook, color: White}).PrettyBoard()
		fmt.Println()
		fmt.Println("-----------------------")
		fmt.Println()
	}
}
