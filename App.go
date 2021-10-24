package main

import "fmt"

type Color int

func main() {

	board := StartBoard()
	KingMoves(*board)
	mvs := WhitePawnStraightMasks()
	fmt.Println(len(mvs))
	board.PrettyBoard()

}
