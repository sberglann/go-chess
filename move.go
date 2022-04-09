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
