package main

// Magics
var bishopMasks = BishopMasks()
var bishopMagics = BishopMagics()
var bishopBits = BishopBits()
var bishopMoveTable = BishopMoveTable()

var rookMasks = RookMasks()
var rookMagics = RookMagics()
var rookBits = RookBits()
var rookMoveTable = RookMoveTable()

// Move tables
var kingMoves = KingMasks()
var whitePawnAttackMoves = WhitePawnAttackMasks()
var whitePawnDoubleMoves = WhitePawnDoubleMasks()
var whitePawnStraightMoves = WhitePawnStraightMasks()
var blackPawnAttackMoves = BlackPawnAttackMasks()
var blackPawnDoubleMoves = BlackPawnDoubleMasks()
var blackPawnStraightMoves = BlackPawnStraightMasks()
var knightMoves = KnightMasks()

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
			if IsNotSelfCapture(bb, m) {
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
			straight = whitePawnStraightMoves[pos]
			double = whitePawnDoubleMoves[pos]
			attack = whitePawnAttackMoves[pos]
		} else {
			straight = blackPawnStraightMoves[pos]
			double = blackPawnDoubleMoves[pos]
			attack = blackPawnAttackMoves[pos]
		}

		// These are used to see if a piece is blocking the destinations of the pawn.
		var straightOffset, doubleStraightOffset int
		if bb.Turn() == White {
			straightOffset = 8
			doubleStraightOffset = 16
		} else {
			straightOffset = -8
			doubleStraightOffset = -16
		}
		for _, m := range straight {
			if bb.PieceAt(m.Origin()+straightOffset).piece == Empty {
				validMoves = append(validMoves, m)
			}
		}
		for _, m := range double {
			if bb.PieceAt(m.Origin()+straightOffset).piece == Empty && bb.PieceAt(m.Origin()+doubleStraightOffset).piece == Empty {
				validMoves = append(validMoves, m)
			}
		}
		for _, m := range attack {
			if IsNotSelfCapture(bb, m) && IsCapture(bb, m) {
				validMoves = append(validMoves, m)
			}
		}
	}
	return validMoves
}

func KnightMoves(bb BitBoard) []Move {
	var validMoves []Move
	knights := bb.KnightBB & bb.TurnBoard()
	for knights > 0 {
		pos, newKnights := PopFistBit(knights)
		knights = newKnights
		for _, move := range knightMoves[pos] {
			if IsNotSelfCapture(bb, move) {
				validMoves = append(validMoves, move)
			}
		}
	}
	return validMoves
}

func BishopMoves(bb BitBoard) []Move {
	var validMoves []Move
	bishops := bb.BishopBB & bb.TurnBoard()
	for bishops > 0 {
		pos, newBishops := PopFistBit(bishops)
		bishops = newBishops
		for _, move := range bishopMovesFromPos(bb, pos) {
			if IsNotSelfCapture(bb, move) {
				validMoves = append(validMoves, move)
			}
		}
	}
	return validMoves
}

func bishopMovesFromPos(bb BitBoard, origin int) []Move {
	var moves []Move
	mask := bishopMasks[origin]
	blockers := mask & (bb.WhiteBB | bb.BlackBB)
	magic := bishopMagics[origin]
	key := int((blockers * magic) >> (64 - bishopBits[origin]))
	legalSquares := bishopMoveTable[MagicKey{Square: origin, Key: key}]
	for legalSquares > 0 {
		destination, newSquares := PopFistBit(legalSquares)
		legalSquares = newSquares
		destinationBits := uint16(destination)
		originBits := uint16(origin) << 6
		flagBits := uint16(0) << 12

		move := Move{destinationBits | originBits | flagBits}
		moves = append(moves, move)
	}
	return moves
}

func RookMoves(bb BitBoard) []Move {
	var validMoves []Move
	rooks := bb.RookBB & bb.TurnBoard()
	for rooks > 0 {
		pos, newRooks := PopFistBit(rooks)
		rooks = newRooks
		for _, move := range rookMovesFromPos(bb, pos) {
			if IsNotSelfCapture(bb, move) {
				validMoves = append(validMoves, move)
			}
		}
	}
	return validMoves
}

func rookMovesFromPos(bb BitBoard, origin int) []Move {
	var moves []Move
	mask := rookMasks[origin]
	blockers := mask & (bb.WhiteBB | bb.BlackBB)
	magic := rookMagics[origin]
	key := int((blockers * magic) >> (64 - rookBits[origin]))
	legalSquares := rookMoveTable[MagicKey{Square: origin, Key: key}]
	for legalSquares > 0 {
		destination, newSquares := PopFistBit(legalSquares)
		legalSquares = newSquares
		destinationBits := uint16(destination)
		originBits := uint16(origin) << 6
		flagBits := uint16(0) << 12

		move := Move{destinationBits | originBits | flagBits}
		moves = append(moves, move)
	}
	return moves
}

func QueenMoves(bb BitBoard) []Move {
	var validMoves []Move
	queens := bb.QueenBB & bb.TurnBoard()
	for queens > 0 {
		origin, newQueens := PopFistBit(queens)
		queens = newQueens
		rook := rookMovesFromPos(bb, origin)
		bishop := bishopMovesFromPos(bb, origin)
		for _, move := range append(rook, bishop...) {
			if IsNotSelfCapture(bb, move) {
				validMoves = append(validMoves, move)
			}
		}
	}
	return validMoves
}

func IsNotSelfCapture(bb BitBoard, m Move) bool {
	return bb.PieceAt(m.Destination()).color != bb.Turn()
}

func IsCapture(bb BitBoard, m Move) bool {
	return bb.PieceAt(m.Destination()).color == bb.Turn().Opposite()
}
