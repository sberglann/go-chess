package gochess

import "fmt"

type Piece int
type Color int

const (
	Empty Piece = iota - 1
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

func (p Piece) Letter() string {
	switch {
	case p == Empty:
		return "-"
	case p == Pawn:
		return "p"
	case p == Knight:
		return "n"
	case p == Bishop:
		return "b"
	case p == Rook:
		return "r"
	case p == Queen:
		return "q"
	case p == King:
		return "k"
	default:
		return "-"
	}
}

const (
	Blank Color = iota - 1
	White
	Black
)

func (c Color) Letter() string {
	switch {
	case c == Black:
		return "b"
	case c == White:
		return "w"
	default:
		return "-"
	}
}

func (c Color) Opposite() Color {
	switch c {
	case Black:
		return White
	case White:
		return Black
	case Blank:
		return Blank
	}
	return Blank
}

type ColoredPiece struct {
	piece Piece
	color Color
}

func (p ColoredPiece) toUnicode() string {
	if p.piece == Empty || p.color == Blank {
		return "·"
	} else if p.piece == Pawn && p.color == White {
		return "♙"
	} else if p.piece == Pawn && p.color == Black {
		return "♟"
	} else if p.piece == Knight && p.color == White {
		return "♘"
	} else if p.piece == Knight && p.color == Black {
		return "♞"
	} else if p.piece == Bishop && p.color == White {
		return "♗"
	} else if p.piece == Bishop && p.color == Black {
		return "♝"
	} else if p.piece == Rook && p.color == White {
		return "♖"
	} else if p.piece == Rook && p.color == Black {
		return "♜"
	} else if p.piece == Queen && p.color == White {
		return "♕"
	} else if p.piece == Queen && p.color == Black {
		return "♛"
	} else if p.piece == King && p.color == White {
		return "♔"
	} else if p.piece == King && p.color == Black {
		return "♚"
	} else {
		fmt.Printf("No piece found for %d %d", p.piece, p.color)
		return " "
	}
}
