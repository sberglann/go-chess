package main

func Eval(board BitBoard, move Move) float64 {
	return board.Eval + (psq(board, move) / 100.0)
}

const (
	pawnWeight   = 1.0
	knightWeigh  = 2.75
	bishopWeight = 3.0
	rookWeight   = 5.0
	queenWeight  = 9.0
	kingWeight   = 20.0
)

var psqScores = map[ColoredPiece][128]float64{
	ColoredPiece{Pawn, White}:   extractPsqScores(psqPawn, true, pawnWeight*100.0),
	ColoredPiece{Pawn, Black}:   extractPsqScores(psqPawn, false, pawnWeight*100.0),
	ColoredPiece{Knight, White}: extractPsqScores(psqKnight, true, knightWeigh*100.0),
	ColoredPiece{Knight, Black}: extractPsqScores(psqKnight, false, knightWeigh*100.0),
	ColoredPiece{Bishop, White}: extractPsqScores(psqBishop, true, bishopWeight*100.0),
	ColoredPiece{Bishop, Black}: extractPsqScores(psqBishop, false, bishopWeight*100.0),
	ColoredPiece{Rook, White}:   extractPsqScores(psqRook, true, rookWeight*100.0),
	ColoredPiece{Rook, Black}:   extractPsqScores(psqRook, false, rookWeight*100.0),
	ColoredPiece{Queen, White}:  extractPsqScores(psqQueen, true, queenWeight*100.0),
	ColoredPiece{Queen, Black}:  extractPsqScores(psqQueen, false, queenWeight*100.0),
	ColoredPiece{King, White}:   extractPsqScores(psqKing, true, kingWeight*100.0),
	ColoredPiece{King, Black}:   extractPsqScores(psqKing, false, kingWeight*100.0),
}

func psq(board BitBoard, move Move) float64 {
	// In the end game we would like to use other psq tables, which are appended to the early game psq tables
	scoreOffset := 0
	if board.TurnCount() > 40 {
		scoreOffset = 64
	}

	origin := move.Origin()
	destination := move.Destination()
	originPiece := board.PieceAt(origin)
	capturedPiece := board.PieceAt(move.Destination())

	captureScore := 0.0
	if capturedPiece.piece != Empty {
		captureScore = psqScores[capturedPiece][destination+scoreOffset]
	}

	return psqScores[originPiece][destination+scoreOffset] - psqScores[originPiece][origin+scoreOffset] - captureScore
}

func extractPsqScores(psq [][]int, white bool, pieceValue float64) [128]float64 {
	// create a map based on psq and use the first element if earlyGame is true
	var psqMapping [128]float64
	var scoreMultiplier float64
	if white {
		scoreMultiplier = 1.0
	} else {
		scoreMultiplier = -1.0
	}
	addPsqScores := func(earlyGame bool) {
		for i, values := range psq {
			var square int
			if white {
				square = i
			} else {
				rank := i / 8
				file := i % 8
				sq := (8-rank)*8 - (8 - file)
				square = sq
			}
			if earlyGame {
				psqMapping[square] = scoreMultiplier * (pieceValue + float64(values[0])/100)
			} else {
				psqMapping[square+64] = scoreMultiplier * (pieceValue + float64(values[1])/100)
			}
		}
	}

	addPsqScores(true)
	addPsqScores(false)

	return psqMapping
}
