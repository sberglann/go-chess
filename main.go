package main

import "fmt"

func main() {
	originalFEN := "8/5k2/3p4/1p1Pp2p/pP2Pp1P/P4P1K/8/8 b - - 99 50"
	b := BoardFromFEN(originalFEN)
	b.PrettyBoard()
	fenAagain := b.ToFEN()
	isEqual := originalFEN == fenAagain
	_ = isEqual

	ms := KingMoves(b)
	for _, m := range ms {
		Transition(b, m, ColoredPiece{piece: King, color: White}).PrettyBoard()
		fmt.Println()
		fmt.Println("-----------------------")
		fmt.Println()
	}
}
