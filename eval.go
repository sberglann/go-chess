package main

import (
	"math/bits"
)

func Eval(board BitBoard) float64 {
	p := psq(board) / 100
	m := material(board)
	score := p + m
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

func psq(board BitBoard) float64 {
	score := func(pieceBoard uint64, psq map[int]float64, isBlack bool) float64 {
		var sum float64
		for pieceBoard > 0 {
			bit, pb := PopFistBit(pieceBoard)
			pieceBoard = pb
			if isBlack {
				sum += psq[63-bit]
			} else {
				sum += psq[bit]
			}
		}
		return sum
	}
	return score(board.WhiteBB&board.PawnBB, PsqPawn, false) +
		score(board.WhiteBB&board.KnightBB, PsqKnight, false) +
		score(board.WhiteBB&board.BishopBB, PsqBishop, false) +
		score(board.WhiteBB&board.RookBB, PsqRook, false) +
		score(board.WhiteBB&board.QueenBB, PsqQueen, false) +
		score(board.WhiteBB&board.KingBB, PsqKing, false) +
		score(board.BlackBB&board.PawnBB, PsqPawn, true) +
		score(board.BlackBB&board.KnightBB, PsqKnight, true) +
		score(board.BlackBB&board.BishopBB, PsqBishop, true) +
		score(board.BlackBB&board.RookBB, PsqRook, true) +
		score(board.BlackBB&board.QueenBB, PsqQueen, true) +
		score(board.BlackBB&board.KingBB, PsqKing, true)

}
