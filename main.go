package main

import "fmt"

func main() {
	b := BoardFromFEN("8/5k2/3p4/1p1Pp2p/pP2Pp1P/P4P1K/8/8 b - - 99 50")
	b.PrettyBoard()
	ms := RookMoves(b)
	for _, m := range ms {
		Transition(b, m, ColoredPiece{piece: Rook, color: White}).PrettyBoard()
		fmt.Println()
		fmt.Println("-----------------------")
		fmt.Println()
	}
}
