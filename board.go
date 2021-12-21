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

var StartBoard = BitBoard{
	WhiteBB:        uint64(0x000000000000feff),
	BlackBB:        uint64(0xffff000000000000),
	InverseWhiteBB: uint64(0xfffffffffff00000),
	InverseBlackBB: uint64(0x0000ffffffffffff),
	PawnBB:         uint64(0x00ff00000000fe00),
	KnightBB:       uint64(0x4200000000000042),
	BishopBB:       uint64(0x2400000000000024),
	RookBB:         uint64(0x8100000000000081),
	QueenBB:        uint64(0x0800000000000008),
	KingBB:         uint64(0x1000000000000010),
	Flags:          uint32(0x000001FE),
}

func (b BitBoard) TurnBoard() uint64 {
	if b.Turn() == White {
		return b.WhiteBB
	} else {
		return b.BlackBB
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
		return ColoredPiece{color: Blank, piece: Empty}
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

func (b BitBoard) PrettyBoard() {
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

var posToBitBoard = PosToBitBoard()

func IndexToCartesian(pos int) (int, int) {
	rank := pos/8 + 1
	file := pos%8 + 1

	return rank, file
}

func CartesianToIndex(rank int, file int) int {
	return (rank-1)*8 + file - 1
}

func Pretty64(number uint64) {
	fmt.Println("------------------------")
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			pos := i*8 + j
			bit := BitAt64(&number, pos)
			fmt.Print(" ")
			if bit == 1 {
				fmt.Print("1")
			} else {
				fmt.Print(".")
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
