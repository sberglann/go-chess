package main

var kingMoves = KingMasks()
var whitePawnAttackMasks = WhitePawnAttackMasks()
var whitePawnDoubleMasks = WhitePawnDoubleMasks()
var whitePawnStraightMasks = WhitePawnStraightMasks()
var blackPawnAttackMasks = BlackPawnAttackMasks()
var blackPawnDoubleMasks = BlackPawnDoubleMasks()
var blackPawnStraightMasks = BlackPawnStraightMasks()

func KingMoves(bb BitBoard) []Move {
	var validMoves []Move
	kings := bb.KingBB & bb.TurnBoard()
	for kings > 0 {
		pos, newKings := PopFistBit(kings)
		kings = newKings
		moves := kingMoves[pos]
		for _, m := range moves {
			if IsValid(bb, m) {
				validMoves = append(validMoves, m)
			}
		}
	}
	return validMoves
}

func PawnMoves(bb BitBoard) []Move {
	var validMoves []Move
	pawns := bb.PawnBB & bb.TurnBoard()
	for pawns > 0 {
		pos, newPawns := PopFistBit(pawns)
		pawns = newPawns
		var straight, double, attack []Move
		if bb.Turn() == White {
			straight = whitePawnStraightMasks[pos]
			double = whitePawnDoubleMasks[pos]
			attack = whitePawnAttackMasks[pos]
		} else {
			straight = blackPawnStraightMasks[pos]
			double = blackPawnDoubleMasks[pos]
			attack = blackPawnAttackMasks[pos]
		}
		for _, m := range straight {
			if IsValid(bb, m) {
				validMoves = append(validMoves, m)
			}
		}
		for _, m := range double {
			if IsValid(bb, m) {
				validMoves = append(validMoves, m)
			}
		}
		for _, m := range attack {
			if IsValid(bb, m) {
				validMoves = append(validMoves, m)
			}
		}
	}
	return validMoves
}

func IsValid(bb BitBoard, m Move) bool {
	return bb.PieceAt(m.Destination()).color != bb.Turn()
}
