package main

import "fmt"

var existing_entries = KingMasks()

func GenerateMoves(bb BitBoard) []Move {
	var turnBoard uint64
	turn := bb.Turn()
	if turn == White {
		turnBoard = bb.WhiteBB
	} else {
		turnBoard = bb.BlackBB
	}

}

func KingMoves(bb BitBoard) {
	kings := bb.KingBB & bb.TurnBoard()
	for kings > 0 {
		pos, newKings := PopFistBit(kings)
		kings = newKings
		KingMasks
	}
	if cp := bb.PieceAt(19); cp.piece == King {
		toAnd := uint64(0x000000001C141C00)
		Pretty64(toAnd)
		fmt.Println()
		fmt.Println()
		switch cp.color {
		case White:
			Pretty64(bb.InverseWhiteBB & toAnd)
		case Black:
			Pretty64(bb.InverseBlackBB & toAnd)
		}
	}
}
