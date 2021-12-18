package main

import "fmt"

var existing_entries = KingMasks()

func KingMoves(bb BitBoard) []Move {
	kings := bb.KingBB & bb.TurnBoard()
	if kings > 0 {
		pos, newKings := PopFistBit(kings)
		kings = newKings
		moves := existing_entries[pos]
		var validMoves []Move
		for _, m := range moves {
			if IsValid(bb, m) {
				validMoves = append(validMoves, m)
			}
		}
		return validMoves
	} else {
		fmt.Println("No king on board :S.")
		return nil
	}
}

func IsValid(bb BitBoard, m Move) bool {
	return bb.PieceAt(m.Destination()).color != bb.Turn()
}
