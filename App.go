package main

type Color int

func main() {


	board := StartBoard()
	KingMoves(*board)
	BlackPawnStraightMasks()
	board.PrettyBoard()
}


