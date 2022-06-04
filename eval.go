package main

import (
	"math/bits"
)

func Eval(board BitBoard) float64 {
	score := material(board)
	return score
}

const (
	pawnWeight   = 1.0
	knightWeigh  = 2.75
	bishopWeight = 3.0
	rookWeight   = 5.0
	queenWeight  = 9.0
	kingWeight   = 100.0
)

func material(board BitBoard) float64 {
	return float64(bits.OnesCount64(board.WhiteBB&board.PawnBB))*pawnWeight +
		float64(bits.OnesCount64(board.WhiteBB&board.KnightBB))*knightWeigh +
		float64(bits.OnesCount64(board.WhiteBB&board.BishopBB))*bishopWeight +
		float64(bits.OnesCount64(board.WhiteBB&board.RookBB))*rookWeight +
		float64(bits.OnesCount64(board.WhiteBB&board.QueenBB))*queenWeight -
		float64(bits.OnesCount64(board.BlackBB&board.PawnBB))*pawnWeight -
		float64(bits.OnesCount64(board.BlackBB&board.KnightBB))*knightWeigh -
		float64(bits.OnesCount64(board.BlackBB&board.BishopBB))*bishopWeight -
		float64(bits.OnesCount64(board.BlackBB&board.RookBB))*rookWeight -
		float64(bits.OnesCount64(board.BlackBB&board.QueenBB))*queenWeight

}
