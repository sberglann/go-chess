package main

import "fmt"

/*
func GenerateLegalMoves(bb BitBoard) []Move {

}

 */

func KingMoves(bb BitBoard) {
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
