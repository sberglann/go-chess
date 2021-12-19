package main

var kingMoves = KingMasks()
var whitePawnAttackMasks = WhitePawnAttackMasks()
var whitePawnDoubleMasks = WhitePawnDoubleMasks()
var whitePawnStraightMasks = WhitePawnStraightMasks()
var blackPawnAttackMasks = BlackPawnAttackMasks()
var blackPawnDoubleMasks = BlackPawnDoubleMasks()
var blackPawnStraightMasks = BlackPawnStraightMasks()

func Transition(b BitBoard, m Move, piece ColoredPiece) BitBoard {

	origin := m.Origin()
	destination := m.Destination()
	originBB := posToBitBoard[origin]
	destinationBB := posToBitBoard[destination]

	makeMove := func(bitboard uint64) uint64 {
		return (bitboard &^ originBB) | destinationBB
	}

	makeMoveInverse := func(bitboard uint64) uint64 {
		return (bitboard | originBB) &^ destinationBB
	}

	moveOrPass := func(currentPiece Piece, pieceBB uint64) uint64 {
		if piece.piece == currentPiece {
			return makeMove(pieceBB)
		} else {
			return pieceBB
		}
	}

	var WhiteBB, BlackBB, InverseWhiteBB, InverseBlackBB, PawnBB, KnightBB, BishopBB, RookBB, QueenBB, KingBB uint64

	if piece.color == White {
		WhiteBB = makeMove(b.WhiteBB)
		InverseWhiteBB = makeMoveInverse(b.InverseWhiteBB)
		BlackBB = b.BlackBB
		InverseBlackBB = b.InverseBlackBB
	} else {
		BlackBB = makeMove(b.BlackBB)
		InverseBlackBB = makeMoveInverse(b.InverseBlackBB)
		WhiteBB = b.WhiteBB
		InverseWhiteBB = b.InverseWhiteBB
	}

	PawnBB = moveOrPass(Pawn, b.PawnBB)
	KnightBB = moveOrPass(Knight, b.KnightBB)
	BishopBB = moveOrPass(Bishop, b.BishopBB)
	RookBB = moveOrPass(Rook, b.RookBB)
	QueenBB = moveOrPass(Queen, b.QueenBB)
	KingBB = moveOrPass(King, b.KingBB)

	return BitBoard{
		WhiteBB:        WhiteBB,
		BlackBB:        BlackBB,
		InverseWhiteBB: InverseWhiteBB,
		InverseBlackBB: InverseBlackBB,
		PawnBB:         PawnBB,
		KnightBB:       KnightBB,
		BishopBB:       BishopBB,
		RookBB:         RookBB,
		QueenBB:        QueenBB,
		KingBB:         KingBB,
		Flags:          b.Flags,
	}
}

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
