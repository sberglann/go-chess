package main

import "fmt"

type Color int

func main() {
	b := StartBoard
	b.PrettyBoard()
	ms := PawnMoves(b)
	for _, m := range ms {
		Transition(b, m, ColoredPiece{piece: Pawn, color: White}).PrettyBoard()
		fmt.Println()
		fmt.Println("-----------------------")
		fmt.Println()
	}
}
