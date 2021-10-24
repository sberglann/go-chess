package main

import "fmt"

type Color int

func main() {

	res := IndexToUint64(7, 2, uint64(4))

	fmt.Println(res)
	board := StartBoard()
	KingMoves(*board)
	mvs := WhitePawnStraightMasks()
	fmt.Println(len(mvs))
	board.PrettyBoard()

}
