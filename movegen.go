package main

// Magics
var bishopMasksWithEdges = BishopMasks(true)
var magicBishopMasks = BishopMasks(false)
var bishopMagics = BishopMagics()
var bishopBits = BishopBits()
var bishopMoveTable = BishopMoveTable()

var rookMasksWithEdges = RookMasks(true)
var magicRookMasks = RookMasks(false)
var rookMagics = RookMagics()
var rookBits = RookBits()
var rookMoveTable = RookMoveTable()

// Move tables
var kingMovesMapping = KingMovesMapping()
var whitePawnAttackMoves = WhitePawnAttackMovesMapping()
var whitePawnDoubleMoves = WhitePawnDoubleMovesMapping()
var whitePawnStraightMoves = WhitePawnStraightMovesMapping()
var whiteEnPassantMoves = WhiteEnPassantMovesMapping()
var blackPawnAttackMoves = BlackPawnAttackMovesMapping()
var blackPawnDoubleMoves = BlackPawnDoubleMovesMapping()
var blackPawnStraightMoves = BlackPawnStraightMovesMapping()
var blackEnPassantMoves = BlackEnPassantMovesMapping()
var knightMovesMapping = KnightMovesMapping()

var knightMasks = KnightMasks()
var kingMasks = KingMasks()
var whitePawnAttackMasks = WhitePawnAttackMasks()
var blackPawnAttackMasks = BlackPawnAttackMasks()

// Castling empty square checks
var wqSideCastleInMask = posToBitBoard(1) | posToBitBoard(2) | posToBitBoard(3) | posToBitBoard(4)
var wkCastleInMask = posToBitBoard(4) | posToBitBoard(5) | posToBitBoard(6)
var bqSideCastleInMask = posToBitBoard(57) | posToBitBoard(58) | posToBitBoard(59) | posToBitBoard(60)
var bkSideCastleInMask = posToBitBoard(60) |posToBitBoard(61) | posToBitBoard(62)

var wqCastleMove = Move{bits: uint32(0xC102)}
var wkCastleMove = Move{bits: uint32(0xC106)}
var bkCastleMove = Move{bits: uint32(0xCF3E)}
var bqCastleMove = Move{bits: uint32(0xCF3A)}

func GenerateLegalStates(b *BitBoard) ([200]BitBoard, int) {
	// Theoretically, there is posisble to generate a position where white has 218 legal moves. That will never happen though.
	var nextStates [200]BitBoard
	kings := b.KingBB & b.TurnBoard()
	if kings <= 0 {
		return [200]BitBoard{}, 0
	}
	currentKingPos, _ := PopFistBit(kings)

	pawnMoves := pawnMoves(b)
	knightMoves := knightMoves(b)
	bishopMoves := bishopMoves(b)
	rookMoves := rookMoves(b)
	queenMoves := queenMoves(b)
	kingMoves := kingMoves(b)
	castlingMoves := castlingMoves(b)

	i := 0
	j := 0
	m := pawnMoves[j]
	for m.bits > 0 {
		nextState := transition(b, &m, Pawn)
		if !isChecked(&nextState, currentKingPos, nextState.Turn()) {
			nextStates[i] = nextState
			i++
		}
		j++
		m = pawnMoves[j]
	}

	j = 0
	m = knightMoves[j]
	for m.bits > 0 {
		nextState := transition(b, &m, Knight)
		if !isChecked(&nextState, currentKingPos, nextState.Turn()) {
			nextStates[i] = nextState
			i++
		}
		j++
		m = knightMoves[j]
	}
	j = 0
	m = bishopMoves[j]
	for m.bits > 0 {
		nextState := transition(b, &m, Bishop)
		if !isChecked(&nextState, currentKingPos, nextState.Turn()) {
			nextStates[i] = nextState
			i++
		}
		j++
		m = bishopMoves[j]
	}
	j = 0
	m = rookMoves[j]
	for m.bits > 0 {
		nextState := transition(b, &m, Rook)
		if !isChecked(&nextState, currentKingPos, nextState.Turn()) {
			nextStates[i] = nextState
			i++
		}
		j++
		m = rookMoves[j]
	}
	j = 0
	m = queenMoves[j]
	for m.bits > 0 {
		nextState := transition(b, &m, Queen)
		if !isChecked(&nextState, currentKingPos, nextState.Turn()) {
			nextStates[i] = nextState
			i++
		}
		j++
		m = queenMoves[j]
	}
	j = 0
	m = kingMoves[j]
	for m.bits > 0 {
		nextState := transition(b, &m, King)
		// Since the king move, we'll have to use the next position when looking for checks.
		newKingPos := m.Destination()
		if !isChecked(&nextState, newKingPos, nextState.Turn()) {
			nextStates[i] = nextState
			i++
		}
		j++
		if j < len(kingMoves) {
			m = kingMoves[j]
		} else {
			m = Move{0}
		}
	}

	j = 0
	m = castlingMoves[j]
	for m.bits > 0 {
		nextState := transition(b, &m, King)
		newKingPos := m.Destination()
		if !isChecked(&nextState, newKingPos, nextState.Turn()) {
			nextStates[i] = nextState
			i++
		}
		j++
		if j < len(castlingMoves) {
			m = castlingMoves[j]
		} else {
			m = Move{0}
		}
	}

	return nextStates, i
}

func isChecked(board *BitBoard, kingPos int, kingColor Color) bool {
	var isCheckByBishopOrQueen, isCheckByRookOrQueen bool
	var attackingPawnMask uint64
	var turnBoard uint64

	if kingColor == White {
		turnBoard = board.WhiteBB
		attackingPawnMask = blackPawnAttackMasks[kingPos]
	} else {
		turnBoard = board.BlackBB
		attackingPawnMask = whitePawnAttackMasks[kingPos]
	}
	if attackingPawnMask&board.PawnBB&turnBoard > 0 {
		return true
	}

	bishopsAndQueens := board.BishopBB | board.QueenBB

	if bishopMasksWithEdges[kingPos]&bishopsAndQueens&turnBoard > 0 {
		// The move resulted in a piece being moved from the diagonal where an opposing bishop can attack
		// the king. We need to see if it is a discovery check.
		blockers := magicBishopMasks[kingPos] & (board.WhiteBB | board.BlackBB)
		magic := bishopMagics[kingPos]
		key := int((blockers * magic) >> (64 - bishopBits[kingPos]))
		bishopAttackMask := bishopMoveTable[MagicKey{Square: kingPos, Key: key}]
		isCheckByBishopOrQueen = bishopAttackMask&bishopsAndQueens&turnBoard > 0
		if isCheckByBishopOrQueen {
			return true
		}
	}
	rooksAndQueens := board.RookBB | board.QueenBB
	if rookMasksWithEdges[kingPos]&rooksAndQueens&turnBoard > 0 {
		// Same as above, just with rooks instead
		blockers := magicRookMasks[kingPos] & (board.WhiteBB | board.BlackBB)
		magic := rookMagics[kingPos]
		key := int((blockers * magic) >> (64 - rookBits[kingPos]))
		rookAttackMask := rookMoveTable[MagicKey{Square: kingPos, Key: key}]
		isCheckByRookOrQueen = rookAttackMask&rooksAndQueens&turnBoard > 0
		if isCheckByRookOrQueen {
			return true
		}
	}

	isCheckByKnight := knightMasks[kingPos]&board.KnightBB&turnBoard > 0

	return isCheckByKnight
}

func transition(b *BitBoard, m *Move, piece Piece) BitBoard {
	var capturedPiece Piece
	var enPassantFile int

	origin := m.Origin()
	destination := m.Destination()
	originBB := posToBitBoard(origin)
	destinationBB := posToBitBoard(destination)
	if m.IsEnPassantMove() {
		enPassantFile = (m.Destination() % 8) + 1
	} else {
		enPassantFile = 0
	}
	oppositeTurnBoard := b.OppositeTurnBoard()

	switch {
	case oppositeTurnBoard&destinationBB <= 0:
		// Not a capture
		capturedPiece = Empty
	case destinationBB&oppositeTurnBoard&b.PawnBB > 0:
		capturedPiece = Pawn
	case destinationBB&oppositeTurnBoard&b.KnightBB > 0:
		capturedPiece = Knight
	case destinationBB&oppositeTurnBoard&b.BishopBB > 0:
		capturedPiece = Bishop
	case destinationBB&oppositeTurnBoard&b.RookBB > 0:
		capturedPiece = Rook
	case destinationBB&oppositeTurnBoard&b.QueenBB > 0:
		capturedPiece = Queen
	case destinationBB&oppositeTurnBoard&b.KingBB > 0:
		capturedPiece = King
	}

	if capturedPiece != Empty {
		CaptureCounter += 1
	}

	makeMove := func(bitboard uint64) uint64 {
		return (bitboard &^ originBB) | destinationBB
	}

	moveOrPass := func(currentPiece Piece, pieceBB uint64) uint64 {
		if currentPiece == piece {
			return makeMove(pieceBB)
		} else if capturedPiece == currentPiece {
			return pieceBB &^ destinationBB
		} else {
			return pieceBB
		}
	}

	// Invert the first bit to change turns
	flags := b.Flags ^ uint32(1)
	// Clear previous double pawn move flag
	flags &^= uint32(0b11110)

	flags &^= uint32(0b111111111111111) << 17
	flags |= uint32(b.TurnCount()+1) << 17

	if m.IsDoublePawnMove() {
		dpfile := (m.Destination() % 8) + 1
		flags |= (uint32(dpfile << 1))
	}

	whiteBB := b.WhiteBB
	blackBB := b.BlackBB
	pawnBB := b.PawnBB
	knightBB := b.KnightBB
	bishopBB := b.BishopBB
	rookBB := b.RookBB
	queenBB := b.QueenBB
	kingBB := b.KingBB

	if enPassantFile > 0 {
		if b.Turn() == White {
			whiteBB = (b.WhiteBB | posToBitBoard(40+enPassantFile-1)) &^ originBB
			blackBB = b.BlackBB &^ posToBitBoard(32+enPassantFile-1)
		} else {
			blackBB = (b.BlackBB | posToBitBoard(16+enPassantFile-1)) &^ originBB
			whiteBB = b.WhiteBB &^ posToBitBoard(24+enPassantFile-1)
		}
	} else {
		if b.Turn() == White {
			whiteBB = makeMove(b.WhiteBB)
			blackBB = b.BlackBB &^ destinationBB
		} else {
			blackBB = makeMove(b.BlackBB)
			whiteBB = b.WhiteBB &^ destinationBB
		}
	}

	if enPassantFile > 0 {
		if b.Turn() == White {
			pawnBB = b.PawnBB | posToBitBoard(40+enPassantFile-1)&^originBB&^posToBitBoard(32+enPassantFile-1)
		} else {
			pawnBB = b.PawnBB | posToBitBoard(16+enPassantFile-1)&^originBB&^posToBitBoard(24+enPassantFile-1)
		}
	} else {
		pawnBB = moveOrPass(Pawn, b.PawnBB)
	}

	if m.IsCastleMove() {
		if m.Destination() == 2 {
			// White queen side.
			rookBB &^= posToBitBoard(0)
			rookBB |= posToBitBoard(3)
			whiteBB &^= posToBitBoard(0)
			whiteBB |= posToBitBoard(3)
			flags &^= uint32(0x60)
		} else if m.Destination() == 6 {
			// White king side
			rookBB &^= posToBitBoard(7)
			rookBB |= posToBitBoard(5)
			whiteBB &^= posToBitBoard(7)
			whiteBB |= posToBitBoard(5)
			flags &^= uint32(0x60)
		} else if m.Destination() == 58 {
			// Black queen side
			rookBB &^= posToBitBoard(56)
			rookBB |= posToBitBoard(59)
			blackBB &^= posToBitBoard(56)
			blackBB |= posToBitBoard(59)
			flags &^= uint32(0x180)
		} else {
			// Black king side
			rookBB &^= posToBitBoard(63)
			rookBB |= posToBitBoard(61)
			blackBB &^= posToBitBoard(63)
			blackBB |= posToBitBoard(61)
			flags &^= uint32(0x180)
		}
	}

	// Disable castling flags if origin is either corner. This means a rook may have moved.
	// Shouldn't matter if we disable these redundantly.
	if m.Origin() == 0 {
		flags &^= uint32(1) << 7
	} else if m.Origin() == 7 {
		flags &^= uint32(1) << 6
	} else if m.Origin() == 56 {
		flags &^= uint32(1) << 9
	} else if m.Origin() == 63 {
		flags &^= uint32(1) << 8
	}

	knightBB = moveOrPass(Knight, knightBB)
	bishopBB = moveOrPass(Bishop, bishopBB)
	rookBB = moveOrPass(Rook, rookBB)
	queenBB = moveOrPass(Queen, queenBB)
	kingBB = moveOrPass(King, kingBB)

	res := BitBoard{
		WhiteBB:  whiteBB,
		BlackBB:  blackBB,
		PawnBB:   pawnBB,
		KnightBB: knightBB,
		BishopBB: bishopBB,
		RookBB:   rookBB,
		QueenBB:  queenBB,
		KingBB:   kingBB,
		Flags:    flags,
	}

	return res
}

func kingMoves(bb *BitBoard) [8]Move {
	var validMoves [8]Move
	kings := bb.KingBB & bb.TurnBoard()
	i := 0
	for kings > 0 {
		pos, newKings := PopFistBit(kings)
		kings = newKings
		moves := kingMovesMapping[pos]
		for _, m := range moves {
			if isNotSelfCapture(bb, &m) {
				validMoves[i] = m
				i++
			}
		}
	}
	return validMoves
}

func pawnMoves(bb *BitBoard) [32]Move {
	var validMoves [32]Move
	pawns := bb.PawnBB & bb.TurnBoard()
	i := 0
	for pawns > 0 {
		pos, newPawns := PopFistBit(pawns)
		pawns = newPawns
		for _, m := range pawnMovesFromPos(bb, pos) {
			if m.bits > 0 {
				validMoves[i] = m
				i++
			}
		}
	}
	return validMoves
}

func pawnMovesFromPos(bb *BitBoard, origin int) [12]Move {
	// Should never be possible to generate more than 12 pawn moves (including different promotion options in all directions) from one single pos.
	var validMoves [12]Move
	var straight, double, attack, enPassant []Move

	if bb.Turn() == White {
		straight = whitePawnStraightMoves[origin]
		double = whitePawnDoubleMoves[origin]
		attack = whitePawnAttackMoves[origin]
		if file := bb.DoublePawnMoveFile(); file > 0 && file <= 8 {
			enPassant = whiteEnPassantMoves[origin]
		}
	} else {
		straight = blackPawnStraightMoves[origin]
		double = blackPawnDoubleMoves[origin]
		attack = blackPawnAttackMoves[origin]
		if file := bb.DoublePawnMoveFile(); file > 0 && file <= 8 {
			enPassant = blackEnPassantMoves[origin]
		}
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

	nextSquareIsEmpty := bb.IsEmpty(origin + straightOffset)

	i := 0
	for _, m := range straight {
		if nextSquareIsEmpty {
			validMoves[i] = m
			i++
		}
	}
	for _, m := range double {
		if nextSquareIsEmpty && bb.IsEmpty(origin+doubleStraightOffset) {
			validMoves[i] = m
			i++
		}
	}
	for _, m := range attack {
		if isNotSelfCapture(bb, &m) && isCapture(bb, &m) {
			validMoves[i] = m
			i++
		}
	}
	for _, m := range enPassant {
		if isValidEnPassantCapture(bb, &m) {
			validMoves[i] = m
			i++
		}
	}

	return validMoves
}

func knightMoves(bb *BitBoard) [32]Move {
	var validMoves [32]Move
	i := 0
	knights := bb.KnightBB & bb.TurnBoard()
	for knights > 0 {
		pos, newKnights := PopFistBit(knights)
		knights = newKnights
		for _, m := range knightMovesFromPos(bb, pos) {
			if m.bits > 0 {
				validMoves[i] = m
				i++
			}
		}
	}
	return validMoves
}

func knightMovesFromPos(bb *BitBoard, origin int) [8]Move {
	var validMoves [8]Move
	i := 0
	for _, move := range knightMovesMapping[origin] {
		if isNotSelfCapture(bb, &move) {
			validMoves[i] = move
			i++
		}
	}
	return validMoves
}

func bishopMoves(bb *BitBoard) [50]Move {
	var validMoves [50]Move
	bishops := bb.BishopBB & bb.TurnBoard()
	i := 0
	for bishops > 0 {
		pos, newBishops := PopFistBit(bishops)
		bishops = newBishops
		for _, move := range bishopMovesFromPos(bb, pos) {
			if isNotSelfCapture(bb, &move) {
				if move.bits > 0 {
					validMoves[i] = move
					i++
				}
			}
		}
	}
	return validMoves
}

func bishopMovesFromPos(bb *BitBoard, origin int) [13]Move {
	var moves [13]Move
	mask := magicBishopMasks[origin]
	blockers := mask & (bb.WhiteBB | bb.BlackBB)
	magic := bishopMagics[origin]
	key := int((blockers * magic) >> (64 - bishopBits[origin]))
	legalSquares := bishopMoveTable[MagicKey{Square: origin, Key: key}]
	i := 0
	for legalSquares > 0 {
		destination, newSquares := PopFistBit(legalSquares)
		legalSquares = newSquares
		destinationBits := uint32(destination)
		originBits := uint32(origin) << 6
		flagBits := uint32(0) << 12

		move := Move{destinationBits | originBits | flagBits}
		moves[i] = move
		i++
	}
	return moves
}

func rookMoves(bb *BitBoard) [50]Move {
	var validMoves [50]Move
	rooks := bb.RookBB & bb.TurnBoard()
	i := 0
	for rooks > 0 {
		pos, newRooks := PopFistBit(rooks)
		rooks = newRooks
		for _, move := range rookMovesFromPos(bb, pos) {
			if isNotSelfCapture(bb, &move) {
				if move.bits > 0 {
					validMoves[i] = move
					i++
				}
			}
		}
	}
	return validMoves
}

func rookMovesFromPos(bb *BitBoard, origin int) [14]Move {
	var moves [14]Move
	mask := magicRookMasks[origin]
	blockers := mask & (bb.WhiteBB | bb.BlackBB)
	magic := rookMagics[origin]
	key := int((blockers * magic) >> (64 - rookBits[origin]))
	legalSquares := rookMoveTable[MagicKey{Square: origin, Key: key}]
	i := 0
	for legalSquares > 0 {
		destination, newSquares := PopFistBit(legalSquares)
		legalSquares = newSquares
		destinationBits := uint32(destination)
		originBits := uint32(origin) << 6
		flagBits := uint32(0) << 12

		move := Move{destinationBits | originBits | flagBits}
		moves[i] = move
		i += 1
	}
	return moves
}

func queenMoves(bb *BitBoard) [100]Move {
	var validMoves [100]Move
	queens := bb.QueenBB & bb.TurnBoard()
	i := 0
	for queens > 0 {
		origin, newQueens := PopFistBit(queens)
		queens = newQueens
		rook := rookMovesFromPos(bb, origin)
		bishop := bishopMovesFromPos(bb, origin)
		for _, move := range rook {
			if isNotSelfCapture(bb, &move) {
				if move.bits > 0 {
					validMoves[i] = move
					i++
				}
			}
		}
		for _, move := range bishop {
			if isNotSelfCapture(bb, &move) {
				validMoves[i] = move
				i++
			}
		}
	}
	return validMoves
}

func castlingMoves(bb *BitBoard) [2]Move {
	canCastle := func(inBetweenSquares uint64) bool {
		if (bb.WhiteBB|bb.BlackBB)&inBetweenSquares == 0 {
			for inBetweenSquares > 0 {
				square, sqs := PopFistBit(inBetweenSquares)
				inBetweenSquares = sqs
				if isChecked(bb, square, bb.OppositeTurn()) {
					return false
				}
			}
			return true
		} else {
			return false
		}
	}
	var validMoves [2]Move
	i := 0
	if bb.Turn() == White {
		correctKingPos := bb.KingBB&bb.WhiteBB&posToBitBoard(4) > 0
		if correctKingPos && bb.WhiteCanCastleKingSite() && canCastle(wkCastleInMask) {
			validMoves[i] = wkCastleMove
			i++
		}
		if correctKingPos && bb.WhiteCanCastleQueenSite() && canCastle(wqSideCastleInMask) {
			validMoves[i] = wqCastleMove
			i++
		}
	} else {
		correctKingPos := bb.KingBB&bb.BlackBB&posToBitBoard(60) > 0
		if correctKingPos && bb.BlackCanCastleKingSite() && canCastle(bkSideCastleInMask) {
			validMoves[i] = bkCastleMove
			i++
		}
		if correctKingPos && bb.BlackCanCastleQueenSite() && canCastle(bqSideCastleInMask) {
			validMoves[i] = bqCastleMove
			i++
		}
	}
	return validMoves
}

func isNotSelfCapture(bb *BitBoard, m *Move) bool {
	return bb.TurnBoard()&posToBitBoard(m.Destination()) == 0
}

func isCapture(bb *BitBoard, m *Move) bool {
	return bb.OppositeTurnBoard()&posToBitBoard(m.Destination()) > 0
}

func isValidEnPassantCapture(bb *BitBoard, m *Move) bool {
	destinationFile := (m.Destination() % 8) + 1
	if bb.DoublePawnMoveFile() == destinationFile {
		EnPassantCounter += 1
		return true
	} else {
		return false
	}
}
