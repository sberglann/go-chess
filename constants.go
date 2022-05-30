package main

var FileToLetter = map[int]string{
	1: "a",
	2: "b",
	3: "c",
	4: "d",
	5: "e",
	6: "f",
	7: "g",
	8: "h",
}

var PieceToMaterialScore = map[Piece]float64{
	Pawn:   1.0,
	Knight: 2.75,
	Bishop: 3.0,
	Rook:   5.0,
	Queen:  9.0,
	King:   100.0,
}
