package main

import "fmt"

type BitBoard struct {
	// Bit boards for colors
	WhiteBB        uint64
	BlackBB        uint64
	InverseWhiteBB uint64
	InverseBlackBB uint64

	// Bit boards for piece types
	PawnBB   uint64
	KnightBB uint64
	BishopBB uint64
	RookBB   uint64
	QueenBB  uint64
	KingBB   uint64

	// Flags for states
	// Bit     0: Whose turn it is to move:
	//	- 0 White
	//	- 1 Black
	// Bit   1-5: The file that recently did a two-square pawn move (for en passant). 1111 is given if inapplicable.
	// Bit 	   6: Whether or not white can castle king-side.
	// Bit 	   7: Whether or not white can castle queen-side.
	// Bit 	   8: Whether or not black can castle king-side.
	// Bit 	   9: Whether or not black can castle queen-side.
	// Bit 10-16: The count since the last piece capture or pawn move. if this counter passes 100, the game is draw.
	Flags uint32
}

func StartBoard() *BitBoard {
	// @fmt:off
	white 	 := uint64(0x000000000008ffff)
	black 	 := uint64(0xffff000000000000)
	invWhite := uint64(0xfffffffffff70000)
	invBlack := uint64(0x0000ffffffffffff)
	pawns 	 := uint64(0x00ff00000000ff00)
	knights  := uint64(0x4200000000000042)
	bishops  := uint64(0x2400000000000024)
	rooks 	 := uint64(0x8100000000000081)
	queens   := uint64(0x0800000000000008)
	kings    := uint64(0x1000000000000010)
	flags    := uint32(0x000001FE)

	// @fmt:on

	return &BitBoard{
		WhiteBB:        white,
		BlackBB:        black,
		InverseWhiteBB: invWhite,
		InverseBlackBB: invBlack,
		PawnBB:         pawns,
		KnightBB:       knights,
		BishopBB:       bishops,
		RookBB:         rooks,
		QueenBB:        queens,
		KingBB:         kings,
		Flags:          flags,
	}
}

func (b *BitBoard) Turn() Color {
	if value := BitAt32(&b.Flags, 0); value == 0 {
		return White
	} else {
		return Black
	}
}

func (b *BitBoard) EnPassantFile() int {
	if value := b.Flags >> 1 & 0xF; value > 7 {
		return -1
	} else {
		return int(value)
	}
}

func (b *BitBoard) PieceAt(pos int) ColoredPiece {
	isPopulated := func(bb *uint64, pos int) bool {
		return BitAt64(bb, pos) == 1
	}
	var color Color
	var piece Piece

	if isPopulated(&b.WhiteBB, pos) {
		color = White
	} else if isPopulated(&b.BlackBB, pos) {
		color = Black
	} else {
		color = Blank
	}

	if isPopulated(&b.PawnBB, pos) {
		piece = Pawn
	} else if isPopulated(&b.KnightBB, pos) {
		piece = Knight
	} else if isPopulated(&b.BishopBB, pos) {
		piece = Bishop
	} else if isPopulated(&b.RookBB, pos) {
		piece = Rook
	} else if isPopulated(&b.QueenBB, pos) {
		piece = Queen
	} else if isPopulated(&b.KingBB, pos) {
		piece = King
	} else {
		piece = Empty
	}

	return ColoredPiece{piece: piece, color: color}
}

func (b *BitBoard) PrettyBoard() {
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			pos := i*8 + j
			coloredPiece := b.PieceAt(pos)
			fmt.Print(" ")
			fmt.Print(coloredPiece.toUnicode())
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

