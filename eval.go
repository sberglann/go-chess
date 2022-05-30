package main

import (
	"math/bits"
)

func Eval(board BitBoard) float64 {
	score := material(board)
	if board.OppositeTurn() == Black {
		return -score
	} else {
		return score
	}
}

func material(board BitBoard) float64 {
	return float64(bits.OnesCount64(board.WhiteBB&board.PawnBB))*PieceToMaterialScore[Pawn] +
		float64(bits.OnesCount64(board.WhiteBB&board.KnightBB))*PieceToMaterialScore[Knight] +
		float64(bits.OnesCount64(board.WhiteBB&board.BishopBB))*PieceToMaterialScore[Bishop] +
		float64(bits.OnesCount64(board.WhiteBB&board.RookBB))*PieceToMaterialScore[Rook] +
		float64(bits.OnesCount64(board.WhiteBB&board.QueenBB))*PieceToMaterialScore[Queen] -
		float64(bits.OnesCount64(board.BlackBB&board.PawnBB))*PieceToMaterialScore[Pawn] -
		float64(bits.OnesCount64(board.BlackBB&board.KnightBB))*PieceToMaterialScore[Knight] -
		float64(bits.OnesCount64(board.BlackBB&board.BishopBB))*PieceToMaterialScore[Bishop] -
		float64(bits.OnesCount64(board.BlackBB&board.RookBB))*PieceToMaterialScore[Rook] -
		float64(bits.OnesCount64(board.BlackBB&board.QueenBB))*PieceToMaterialScore[Queen]

}
