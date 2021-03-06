package main

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

func (m Move) toAlgebraicNotation() []string {
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

func (m Move) Destination() int {
	return int(m.bits & 0x3F)
}
func (m Move) IsDoublePawnMove() bool {
	return m.bits&0xC000 == 0x8000
}

func (m Move) IsEnPassantMove() bool {
	return m.bits&0x10000 > 0
}

func (m Move) IsCastleMove() bool {
	return m.bits&0xC000 == 0xC000
}

func (m Move) Origin() int {
	return int(m.bits & 0xFC0 >> 6)
}

func (m Move) Promotion() Piece {
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
