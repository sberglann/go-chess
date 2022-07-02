package main

import (
	"fmt"
	"strconv"
	"strings"
)

type BitBoard struct {
	// Bit boards for colors
	WhiteBB uint64
	BlackBB uint64

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
	// Bit 17-24: The move counter. Incremented after black has moved.
	Flags uint32
}

var StartBoard = BitBoard{
	WhiteBB:  uint64(0x000000000000ffff),
	BlackBB:  uint64(0xffff000000000000),
	PawnBB:   uint64(0x00ff00000000ff00),
	KnightBB: uint64(0x4200000000000042),
	BishopBB: uint64(0x2400000000000024),
	RookBB:   uint64(0x8100000000000081),
	QueenBB:  uint64(0x0800000000000008),
	KingBB:   uint64(0x1000000000000010),
	Flags:    uint32(0x000001FE),
}

func (b BitBoard) TurnCount() int {
	return int(b.Flags >> 17)
}

func (b BitBoard) TurnBoard() uint64 {
	if b.Turn() == White {
		return b.WhiteBB
	} else {
		return b.BlackBB
	}
}

func (b BitBoard) OppositeTurnBoard() uint64 {
	if b.Turn() == White {
		return b.BlackBB
	} else {
		return b.WhiteBB
	}
}

func (b BitBoard) Turn() Color {
	if BitAt32(b.Flags, 0) == 0 {
		return White
	} else {
		return Black
	}
}

func (b BitBoard) OppositeTurn() Color {
	if b.Turn() == White {
		return Black
	} else {
		return White
	}
}

func (b BitBoard) DoublePawnMoveFile() int {
	if value := b.Flags >> 1 & 0xF; value > 8 || value < 1 {
		return -1
	} else {
		return int(value)
	}
}

func (b BitBoard) WhiteCanCastleKingSite() bool {
	return b.Flags&(uint32(1)<<6) > 0
}

func (b BitBoard) WhiteCanCastleQueenSite() bool {
	return b.Flags&(uint32(1)<<7) > 0
}

func (b BitBoard) BlackCanCastleKingSite() bool {
	return b.Flags&(uint32(1)<<8) > 0
}

func (b BitBoard) BlackCanCastleQueenSite() bool {
	return b.Flags&(uint32(1)<<9) > 0
}

func (b BitBoard) IsEmpty(pos int) bool {
	res := posToBitBoard(pos)&(b.WhiteBB|b.BlackBB) == 0
	return res
}

func (b BitBoard) PieceAt(pos int) ColoredPiece {
	isPopulated := func(bb uint64, pos int) bool {
		return BitAt64(bb, pos) == 1
	}
	var color Color
	var piece Piece

	if isPopulated(b.WhiteBB, pos) {
		color = White
	} else if isPopulated(b.BlackBB, pos) {
		color = Black
	} else {
		return ColoredPiece{color: Blank, piece: Empty}
	}

	if isPopulated(b.PawnBB, pos) {
		piece = Pawn
	} else if isPopulated(b.KnightBB, pos) {
		piece = Knight
	} else if isPopulated(b.BishopBB, pos) {
		piece = Bishop
	} else if isPopulated(b.RookBB, pos) {
		piece = Rook
	} else if isPopulated(b.QueenBB, pos) {
		piece = Queen
	} else if isPopulated(b.KingBB, pos) {
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
	println("------------------------")
}

func (b BitBoard) DebugBoard() string {
	var result string
	print := func(board uint64) {
		for i := 7; i >= 0; i-- {
			for j := 0; j < 8; j++ {
				pos := i*8 + j
				if BitAt64(board, pos) == 1 {
					result += " 1 "
				} else {
					result += " 0 "
				}
			}
			result += "\n"
		}
		result += "------------------------\n"
	}

	result += "White\n"
	print(b.WhiteBB)
	result += "Black\n"
	print(b.BlackBB)
	result += "Pawns\n"
	print(b.PawnBB)
	result += "Knights\n"
	print(b.KnightBB)
	result += "Bishops\n"
	print(b.BishopBB)
	result += "Rooks\n"
	print(b.RookBB)
	result += "Queens\n"
	print(b.QueenBB)
	result += "Kings\n"
	print(b.KingBB)

	result += "Flags\n"
	for i := 0; i < 32; i++ {
		if BitAt32(b.Flags, i) == 1 {
			result += "1"
		} else {
			result += "0"
		}
	}
	fmt.Println(result)
	return result
}

func (b BitBoard) ToFEN() string {
	var rankStrings []string
	var rankLetters string
	for rank := 8; rank > 0; rank-- {
		emptyCounter := 0
		rankLetters = ""
		for file := 1; file <= 8; file++ {
			pos := CartesianToIndex(rank, file)
			piece := b.PieceAt(pos)
			if piece.color == Blank {
				emptyCounter += 1
			} else {
				if emptyCounter > 0 {
					rankLetters += strconv.Itoa(emptyCounter)
					emptyCounter = 0
				}
				if piece.color == Black {
					rankLetters += piece.piece.Letter()
				} else {
					rankLetters += strings.ToUpper(piece.piece.Letter())
				}
			}
			if emptyCounter > 0 && file == 8 {
				rankLetters += strconv.Itoa(emptyCounter)
				emptyCounter = 0
			}
		}
		rankStrings = append(rankStrings, rankLetters)
	}
	pieceString := strings.Join(rankStrings, "/")
	activeColor := b.Turn().Letter()

	var castlingRightLetters []string
	if BitAt32(b.Flags, 6) == 1 {
		castlingRightLetters = append(castlingRightLetters, "K")
	}
	if BitAt32(b.Flags, 7) == 1 {
		castlingRightLetters = append(castlingRightLetters, "Q")
	}
	if BitAt32(b.Flags, 8) == 1 {
		castlingRightLetters = append(castlingRightLetters, "k")
	}
	if BitAt32(b.Flags, 9) == 1 {
		castlingRightLetters = append(castlingRightLetters, "q")
	}

	var castlingRights string
	if len(castlingRightLetters) > 0 {
		castlingRights = strings.Join(castlingRightLetters, "")
	} else {
		castlingRights = "-"
	}

	var enPassentString string
	enPassantFile := int(b.Flags >> 1 & uint32(15))
	if enPassantFile > 7 {
		enPassentString = "-"
	} else {
		// The file number is 0-indexed, while the mapping is 1-indexed.
		file := FileToLetter[enPassantFile+1]
		if b.Turn() == White {
			// Black moved previous turn, and the pawn rank is 6
			enPassentString = file + "6"
		} else {
			enPassentString = file + "3"
		}
	}

	halfMoveClock := int(b.Flags >> 10 & uint32(127))
	moveCounter := int(b.Flags >> 17)

	values := []string{
		pieceString,
		activeColor,
		castlingRights,
		enPassentString,
		strconv.Itoa(halfMoveClock),
		strconv.Itoa(moveCounter),
	}

	return strings.Join(values, " ")
}

func posToBitBoard(i int) uint64 {
	return uint64(1) << i
}

var posToBitBoardLol = PosToBitBoard()

func IndexToCartesian(pos int) (int, int) {
	rank := pos/8 + 1
	file := pos%8 + 1

	return rank, file
}

func CartesianToIndex(rank int, file int) int {
	return (rank-1)*8 + file - 1
}

func BoardFromFEN(fen string) BitBoard {
	var whiteBB, blackBB, pawnBB, knightBB, bishopBB, rookBB, queenBB, kingBB uint64
	var flags uint32
	var castling, enPassant, halfMoveClockRaw, fullMoveNumberRaw string
	split := strings.Split(fen, " ")
	pieces := split[0]
	// Not all FEN strings define meta data. If only pieces are given, fallback to default values.
	if len(split) >= 6 {
		castling = split[2]
		enPassant = split[3]
		halfMoveClockRaw = split[4]
		fullMoveNumberRaw = split[5]
	} else {
		castling = "KQkq"
		enPassant = "-"
		halfMoveClockRaw = "0"
		fullMoveNumberRaw = "0"
	}

	halfMoveClock, _ := strconv.Atoi(halfMoveClockRaw)
	fullMoveNumber, _ := strconv.Atoi(fullMoveNumberRaw)
	flags |= uint32(halfMoveClock) << 10
	flags |= uint32(fullMoveNumber) << 17

	// If turn isn't defined, fallback to white
	if len(split) >= 2 {
		if activeColor := split[1]; activeColor == "b" {
			flags |= uint32(1)
		}
	}
	switch {
	case strings.Contains(enPassant, "a"):
		flags |= uint32(1) << 1
	case strings.Contains(enPassant, "b"):
		flags |= uint32(2) << 1
	case strings.Contains(enPassant, "c"):
		flags |= uint32(3) << 1
	case strings.Contains(enPassant, "d"):
		flags |= uint32(4) << 1
	case strings.Contains(enPassant, "e"):
		flags |= uint32(5) << 1
	case strings.Contains(enPassant, "f"):
		flags |= uint32(6) << 1
	case strings.Contains(enPassant, "g"):
		flags |= uint32(7) << 1
	case strings.Contains(enPassant, "h"):
		flags |= uint32(8) << 1
	default:
		flags |= uint32(0) << 1
	}

	if strings.Contains(castling, "K") {
		flags |= uint32(1) << 6
	}
	if strings.Contains(castling, "Q") {
		flags |= uint32(1) << 7
	}
	if strings.Contains(castling, "k") {
		flags |= uint32(1) << 8
	}
	if strings.Contains(castling, "q") {
		flags |= uint32(1) << 9
	}

	piecesByRank := strings.Split(pieces, "/")
	for rankIndex, rankLetters := range piecesByRank {

		// FEN is given with rank 8 first. Subtract the rankIndex to obtain the correct rank.
		rank := 8 - rankIndex
		letterIndex := 0
		file := 1
		for file <= 8 {
			pos := CartesianToIndex(rank, file)
			value := rankLetters[letterIndex]
			letterIndex += 1
			file += 1
			switch {
			case value == 'p':
				blackBB |= (1 << pos)
				pawnBB |= (1 << pos)
			case value == 'n':
				blackBB |= (1 << pos)
				knightBB |= (1 << pos)
			case value == 'b':
				blackBB |= (1 << pos)
				bishopBB |= (1 << pos)
			case value == 'r':
				blackBB |= (1 << pos)
				rookBB |= (1 << pos)
			case value == 'q':
				blackBB |= (1 << pos)
				queenBB |= (1 << pos)
			case value == 'k':
				blackBB |= (1 << pos)
				kingBB |= (1 << pos)
			case value == 'P':
				whiteBB |= (1 << pos)
				pawnBB |= (1 << pos)
			case value == 'N':
				whiteBB |= (1 << pos)
				knightBB |= (1 << pos)
			case value == 'B':
				whiteBB |= (1 << pos)
				bishopBB |= (1 << pos)
			case value == 'R':
				whiteBB |= (1 << pos)
				rookBB |= (1 << pos)
			case value == 'Q':
				whiteBB |= (1 << pos)
				queenBB |= (1 << pos)
			case value == 'K':
				whiteBB |= (1 << pos)
				kingBB |= (1 << pos)
			default:
				numEmpty, _ := strconv.Atoi(string(value))
				file += numEmpty - 1
			}
		}
	}

	board := BitBoard{
		WhiteBB:  whiteBB,
		BlackBB:  blackBB,
		PawnBB:   pawnBB,
		KnightBB: knightBB,
		BishopBB: bishopBB,
		RookBB:   rookBB,
		QueenBB:  queenBB,
		KingBB:   kingBB,
		Flags:    flags,
	}

	return board
}

func Pretty64(number uint64) {
	fmt.Println("------------------------")
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			pos := i*8 + j
			bit := BitAt64(number, pos)
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
