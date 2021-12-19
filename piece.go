package main

import "fmt"

type Piece int

const (
	Empty Piece = iota - 1
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

const (
	Blank Color = iota - 1
	White
	Black
)

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
