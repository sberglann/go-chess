package main

import (
	"fmt"
	"math"
)

type Color int

func main() {

	magic := uint64(72101608930526208)
	blockers := uint64(math.Pow(2, 17) + math.Pow(2, 53))

	key := blockers * magic >> 57
	Pretty64(key)
	fmt.Println(key)

	Pretty64(uint64(9025933902745632))

	board := StartBoard()
	KingMoves(*board)
	board.PrettyBoard()

}
