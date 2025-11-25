package main

import (
	"math/bits"
)

func Eval(board *BitBoard) float64 {
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

var (
	PsqPawnWhiteEarly   = extractPsqScores(psqPawn, true, true)
	PsqPawnWhiteLate    = extractPsqScores(psqPawn, false, true)
	PsqPawnBlackEarly   = extractPsqScores(psqPawn, true, false)
	PsqPawnBlackLate    = extractPsqScores(psqPawn, false, false)
	PsqKnightWhiteEarly = extractPsqScores(psqKnight, true, true)
	PsqKnightWhiteLate  = extractPsqScores(psqKnight, false, true)
	PsqKnightBlackEarly = extractPsqScores(psqKnight, true, false)
	PsqKnightBlackLate  = extractPsqScores(psqKnight, false, false)
	PsqBishopWhiteEarly = extractPsqScores(psqBishop, true, true)
	PsqBishopWhiteLate  = extractPsqScores(psqBishop, false, true)
	PsqBishopBlackEarly = extractPsqScores(psqBishop, true, false)
	PsqBishopBlackLate  = extractPsqScores(psqBishop, false, false)
	PsqRookWhiteEarly   = extractPsqScores(psqRook, true, true)
	PsqRookWhiteLate    = extractPsqScores(psqRook, false, true)
	PsqRookBlackEarly   = extractPsqScores(psqRook, true, false)
	PsqRookBlackLate    = extractPsqScores(psqRook, false, false)
	PsqQueenWhiteEarly  = extractPsqScores(psqQueen, true, true)
	PsqQueenWhiteLate   = extractPsqScores(psqQueen, false, true)
	PsqQueenBlackEarly  = extractPsqScores(psqQueen, true, false)
	PsqQueenBlackLate   = extractPsqScores(psqQueen, false, false)
	PsqKingWhiteEarly   = extractPsqScores(psqKing, true, true)
	PsqKingWhiteLate    = extractPsqScores(psqKing, false, true)
	PsqKingBlackLate    = extractPsqScores(psqKing, false, false)
	PsqKingBlackEarly   = extractPsqScores(psqKing, true, false)
)

func material(board *BitBoard) float64 {
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

func psq(board *BitBoard) float64 {
	score := func(pieceBoard uint64, psq *[64]float64) float64 {
		var sum float64
		for pieceBoard > 0 {
			bit, pb := PopFistBit(pieceBoard)
			pieceBoard = pb
			sum += psq[bit]
		}
		return sum
	}
	if board.TurnCount() <= 40 {
		return score(board.WhiteBB&board.PawnBB, &PsqPawnWhiteEarly) +
			score(board.WhiteBB&board.KnightBB, &PsqKnightWhiteEarly) +
			score(board.WhiteBB&board.BishopBB, &PsqBishopWhiteEarly) +
			score(board.WhiteBB&board.RookBB, &PsqRookWhiteEarly) +
			score(board.WhiteBB&board.QueenBB, &PsqQueenWhiteEarly) +
			score(board.WhiteBB&board.KingBB, &PsqKingWhiteEarly) -
			score(board.BlackBB&board.PawnBB, &PsqPawnBlackEarly) -
			score(board.BlackBB&board.KnightBB, &PsqKnightBlackEarly) -
			score(board.BlackBB&board.BishopBB, &PsqBishopBlackEarly) -
			score(board.BlackBB&board.RookBB, &PsqRookBlackEarly) -
			score(board.BlackBB&board.QueenBB, &PsqQueenBlackEarly) -
			score(board.BlackBB&board.KingBB, &PsqKingBlackEarly)
	} else {
		return score(board.WhiteBB&board.PawnBB, &PsqPawnWhiteLate) +
			score(board.WhiteBB&board.KnightBB, &PsqKnightWhiteLate) +
			score(board.WhiteBB&board.BishopBB, &PsqBishopWhiteLate) +
			score(board.WhiteBB&board.RookBB, &PsqRookWhiteLate) +
			score(board.WhiteBB&board.QueenBB, &PsqQueenWhiteLate) +
			score(board.WhiteBB&board.KingBB, &PsqKingWhiteLate) -
			score(board.BlackBB&board.PawnBB, &PsqPawnBlackLate) -
			score(board.BlackBB&board.KnightBB, &PsqKnightBlackLate) -
			score(board.BlackBB&board.BishopBB, &PsqBishopBlackLate) -
			score(board.BlackBB&board.RookBB, &PsqRookBlackLate) -
			score(board.BlackBB&board.QueenBB, &PsqQueenBlackLate) -
			score(board.BlackBB&board.KingBB, &PsqKingBlackLate)
	}
}

func extractPsqScores(psq [][]int, earlyGame bool, white bool) [64]float64 {
	// create a map based on psq and use the first element if earlyGame is true
	var psqMapping [64]float64
	for i, values := range psq {
		var square int
		if white {
			square = i
		} else {
			rank := i / 8
			file := i % 8
			sq := (8-rank)*8 - (8 - file)
			square = sq
		}
		if earlyGame {
			psqMapping[square] = float64(values[0]) / 100
		} else {
			psqMapping[square] = float64(values[1]) / 100
		}
	}

	return psqMapping
}
