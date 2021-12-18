package main

type Move struct {
	// Bit   0-5: destination square [0, 63]
	// Bit  6-11: origin square [0, 63]
	// Bit 12-13: promotion piece
	//	- 00: Knight
	// 	- 01: Bishop
	//	- 10: Rook
	//  - 11: Queen
	// Bit 14-15: special move type
	//  - 01: Promotion
	//  - 10: En passant
	//  - 11: Castling
	bits uint16
}

func (m Move) Destination() int {
	return int(m.bits & 0x003F)
}

func (m Move) Origin() int {
	return int(m.bits & 0x0FC0 >> 6)
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
