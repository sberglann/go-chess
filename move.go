package main

import "strings"

type Move struct {
	// Bit   0-5: destination square [0, 63]
	// Bit  6-11: origin square [0, 63]
	// Bit 12-13: promotion piece
	//	- 00: Knight
	// 	- 01: Bishop
	//	- 10: Rook
	//  - 11: Queen
	// Bit 14-16: special move type
	//  - 001: Promotion
	//  - 010: Double pawn move, making en passant possible in the next move
	//  - 011: Castling
	//  - 100: en passant
	bits uint32
}

func (m *Move) toAlgebraicNotation() []string {
	if !m.IsCastleMove() {
		o := IndexToAlgebraic[m.Origin()]
		d := IndexToAlgebraic[m.Destination()]
		return []string{o + "-" + d}
	} else {
		if m.Destination() == 2 {
			// White queen side
			return []string{"e1-c1", "a1-d1"}
		} else if m.Destination() == 6 {
			// White king side
			return []string{"e1-g1", "h1-f1"}
		} else if m.Destination() == 56 {
			// Black queen side
			return []string{"e8-c8", "a8-d8"}
		} else if m.Destination() == 62 {
			// White king side
			return []string{"e8-g8", "h8-f8"}
		} else {
			println("No castle match when computing algebraic notation.")
			return []string{}
		}
	}
}

func (m *Move) Destination() int {
	return int(m.bits & 0x3F)
}
func (m *Move) IsDoublePawnMove() bool {
	return m.bits&0xC000 == 0x8000
}

func (m *Move) IsEnPassantMove() bool {
	return m.bits&0x10000 > 0
}

func (m *Move) IsCastleMove() bool {
	return m.bits&0xC000 == 0xC000
}

func (m *Move) Origin() int {
	return int(m.bits & 0xFC0 >> 6)
}

func (m *Move) Promotion() Piece {
	if value := m.bits & 0xF000; value < 12 {
		return Empty
	} else if value == 12 {
		return Knight
	} else if value == 13 {
		return Bishop
	} else if value == 14 {
		return Rook
	} else if value == 15 {
		return Queen
	} else {
		return Empty
	}
}

// ToUCI converts a Move to UCI format (e.g., "e2e4", "e1g1" for castling)
func (m Move) ToUCI() string {
	if m.IsCastleMove() {
		// Castling moves: e1g1 (white kingside), e1c1 (white queenside), e8g8 (black kingside), e8c8 (black queenside)
		if m.Destination() == 6 {
			return "e1g1" // White kingside
		} else if m.Destination() == 2 {
			return "e1c1" // White queenside
		} else if m.Destination() == 62 {
			return "e8g8" // Black kingside
		} else if m.Destination() == 58 {
			return "e8c8" // Black queenside
		}
	}
	
	origin := IndexToAlgebraic[m.Origin()]
	dest := IndexToAlgebraic[m.Destination()]
	
	// UCI format is just the squares concatenated (e.g., "e2e4")
	uci := origin + dest
	
	// Add promotion piece if applicable
	if m.Promotion() != Empty {
		promotionLetter := m.Promotion().Letter()
		uci += promotionLetter
	}
	
	return uci
}

// FindMoveBetweenBoards finds the Move that transforms fromBoard to toBoard
// Returns the move and true if found, or an empty move and false if not found
func FindMoveBetweenBoards(fromBoard *BitBoard, toBoard *BitBoard) (Move, bool) {
	// Get all legal moves with their piece types
	pawnMovesList := pawnMoves(fromBoard)
	knightMovesList := knightMoves(fromBoard)
	bishopMovesList := bishopMoves(fromBoard)
	rookMovesList := rookMoves(fromBoard)
	queenMovesList := queenMoves(fromBoard)
	kingMovesList := kingMoves(fromBoard)
	castlingMovesList := castlingMoves(fromBoard)
	
	// Check pawn moves
	for _, move := range pawnMovesList {
		if move.bits > 0 {
			testBoard := transition(fromBoard, &move, Pawn)
			if boardsEqual(&testBoard, toBoard) {
				return move, true
			}
		}
	}
	
	// Check knight moves
	for _, move := range knightMovesList {
		if move.bits > 0 {
			testBoard := transition(fromBoard, &move, Knight)
			if boardsEqual(&testBoard, toBoard) {
				return move, true
			}
		}
	}
	
	// Check bishop moves
	for _, move := range bishopMovesList {
		if move.bits > 0 {
			testBoard := transition(fromBoard, &move, Bishop)
			if boardsEqual(&testBoard, toBoard) {
				return move, true
			}
		}
	}
	
	// Check rook moves
	for _, move := range rookMovesList {
		if move.bits > 0 {
			testBoard := transition(fromBoard, &move, Rook)
			if boardsEqual(&testBoard, toBoard) {
				return move, true
			}
		}
	}
	
	// Check queen moves
	for _, move := range queenMovesList {
		if move.bits > 0 {
			testBoard := transition(fromBoard, &move, Queen)
			if boardsEqual(&testBoard, toBoard) {
				return move, true
			}
		}
	}
	
	// Check king moves
	for _, move := range kingMovesList {
		if move.bits > 0 {
			testBoard := transition(fromBoard, &move, King)
			if boardsEqual(&testBoard, toBoard) {
				return move, true
			}
		}
	}
	
	// Check castling moves (also use King as piece type)
	for _, move := range castlingMovesList {
		if move.bits > 0 {
			testBoard := transition(fromBoard, &move, King)
			if boardsEqual(&testBoard, toBoard) {
				return move, true
			}
		}
	}
	
	return Move{}, false
}

// boardsEqual checks if two BitBoards are equal
func boardsEqual(b1 *BitBoard, b2 *BitBoard) bool {
	return b1.WhiteBB == b2.WhiteBB &&
		b1.BlackBB == b2.BlackBB &&
		b1.PawnBB == b2.PawnBB &&
		b1.KnightBB == b2.KnightBB &&
		b1.BishopBB == b2.BishopBB &&
		b1.RookBB == b2.RookBB &&
		b1.QueenBB == b2.QueenBB &&
		b1.KingBB == b2.KingBB &&
		b1.Flags == b2.Flags
}

// UCIToMove converts a UCI move string (e.g., "e2e4", "e1g1") to a Move
// Returns the move and true if found, or empty move and false if not found
func UCIToMove(board BitBoard, uciMove string) (Move, bool) {
	if len(uciMove) < 4 {
		return Move{}, false
	}
	
	// Parse UCI string: "e2e4" or "e2e4q" (with promotion)
	originStr := uciMove[0:2]
	destStr := uciMove[2:4]
	var promotion Piece = Empty
	if len(uciMove) > 4 {
		switch strings.ToLower(uciMove[4:5]) {
		case "q": promotion = Queen
		case "r": promotion = Rook
		case "b": promotion = Bishop
		case "n": promotion = Knight
		}
	}
	
	// Convert algebraic notation to indices
	var originIdx, destIdx int = -1, -1
	for idx, square := range IndexToAlgebraic {
		if square == originStr {
			originIdx = idx
		}
		if square == destStr {
			destIdx = idx
		}
	}
	if originIdx == -1 || destIdx == -1 {
		return Move{}, false
	}
	
	// Get piece at origin and generate moves only from that square
	pieceInfo := board.PieceAt(originIdx)
	if pieceInfo.piece == Empty {
		return Move{}, false
	}
	
	// Special handling for castling moves: detect by king origin and castling destination
	// This is needed because castling rights might be lost when reconstructing from moves
	if pieceInfo.piece == King {
		// Check if this looks like a castling move
		// White: e1->g1 (kingside) or e1->c1 (queenside)
		// Black: e8->g8 (kingside) or e8->c8 (queenside)
		if (originIdx == 4 && destIdx == 6) || // White kingside
		   (originIdx == 4 && destIdx == 2) || // White queenside
		   (originIdx == 60 && destIdx == 62) || // Black kingside
		   (originIdx == 60 && destIdx == 58) { // Black queenside
			// This is a castling move - return the appropriate castling move
			if originIdx == 4 && destIdx == 6 {
				return wkCastleMove, true
			} else if originIdx == 4 && destIdx == 2 {
				return wqCastleMove, true
			} else if originIdx == 60 && destIdx == 62 {
				return bkCastleMove, true
			} else if originIdx == 60 && destIdx == 58 {
				return bqCastleMove, true
			}
		}
	}
	
	// Generate moves from the specific origin square based on piece type
	var moves []Move
	switch pieceInfo.piece {
	case Pawn:
		pm := pawnMovesFromPos(&board, originIdx)
		moves = pm[:]
	case Knight:
		km := knightMovesFromPos(&board, originIdx)
		moves = km[:]
	case Bishop:
		bm := bishopMovesFromPos(&board, originIdx)
		moves = bm[:]
	case Rook:
		rm := rookMovesFromPos(&board, originIdx)
		moves = rm[:]
	case Queen:
		// Queen moves are bishop + rook moves from the same position
		bm := bishopMovesFromPos(&board, originIdx)
		rm := rookMovesFromPos(&board, originIdx)
		moves = append(bm[:], rm[:]...)
	case King:
		// Check regular king moves first
		km := kingMoves(&board)
		moves = km[:]
		// Also check castling moves (in case castling rights still exist)
		cm := castlingMoves(&board)
		moves = append(moves, cm[:]...)
	}
	
	// Find the matching move
	for _, move := range moves {
		if move.bits > 0 && move.Origin() == originIdx && move.Destination() == destIdx {
			if promotion == Empty || move.Promotion() == promotion {
				return move, true
			}
		}
	}
	
	return Move{}, false
}

// ApplyUCIMove applies a UCI move string to a board and returns the new board
func ApplyUCIMove(board BitBoard, uciMove string) (BitBoard, bool) {
	move, found := UCIToMove(board, uciMove)
	if !found {
		return board, false
	}
	
	// Determine piece type from the origin square
	pieceAtOrigin := board.PieceAt(move.Origin())
	if pieceAtOrigin.piece == Empty {
		return board, false
	}
	
	// Apply the move
	return transition(&board, &move, pieceAtOrigin.piece), true
}
