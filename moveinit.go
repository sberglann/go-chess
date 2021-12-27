package main

import (
	"math"
	"strconv"
	"strings"
)

var kingOffsets = [][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}}
var knightOffsets = [][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
var whitePawnStraightOffsets = [][2]int{{1, 0}}
var whitePawnDoubleOffsets = [][2]int{{2, 0}}
var whitePawnAttackOffsets = [][2]int{{1, -1}, {1, 1}}
var blackPawnStraightOffsets = [][2]int{{-1, 0}}
var blackPawnDoubleOffsets = [][2]int{{-2, 0}}
var blackPawnAttackOffsets = [][2]int{{-1, -1}, {-1, 1}}

// Magic move gen

type MagicKey struct {
	Square, Key int
}

var promotionFlags = []int{4, 5, 6, 7}

func BishopMasks(includeEdges bool) map[int]uint64 {
	var mapping = make(map[int]uint64)
	for i := 0; i < 64; i++ {
		mapping[i] = GenerateBishopMask(i, includeEdges)
	}
	return mapping
}

func BishopMagics() map[int]uint64 {
	var mapping = make(map[int]uint64)
	lines, _ := ReadLines("move_tables/magics_bishop.txt")
	for _, line := range lines {
		split := strings.Split(line, ";")
		pos, _ := strconv.Atoi(split[0])
		magic, _ := strconv.ParseUint(split[1], 10, 64)
		mapping[pos] = magic
	}
	return mapping
}

func BishopBits() map[int]int {
	var mapping = make(map[int]int)
	for i, bits := range BishopBitsArray {
		mapping[i] = bits
	}
	return mapping
}

func BishopMoveTable() map[MagicKey]uint64 {
	var mapping = make(map[MagicKey]uint64)
	lines, _ := ReadLines("move_tables/bishop.txt")
	for _, line := range lines {
		split := strings.Split(line, ";")
		pos, _ := strconv.Atoi(split[0])
		key, _ := strconv.Atoi(split[1])
		move, _ := strconv.ParseUint(split[2], 10, 64)
		mapping[MagicKey{Square: pos, Key: key}] = move
	}
	return mapping
}

func RookMasks(includeEdges bool) map[int]uint64 {
	var mapping = make(map[int]uint64)
	for i := 0; i < 64; i++ {
		mapping[i] = GenerateRookMask(i, includeEdges)
	}
	return mapping
}

func RookMagics() map[int]uint64 {
	var mapping = make(map[int]uint64)
	lines, _ := ReadLines("move_tables/magics_rook.txt")
	for _, line := range lines {
		split := strings.Split(line, ";")
		pos, _ := strconv.Atoi(split[0])
		magic, _ := strconv.ParseUint(split[1], 10, 64)
		mapping[pos] = magic
	}
	return mapping
}

func RookBits() map[int]int {
	var mapping = make(map[int]int)
	for i, bits := range RookBitsArray {
		mapping[i] = bits
	}
	return mapping
}

func RookMoveTable() map[MagicKey]uint64 {
	var mapping = make(map[MagicKey]uint64)
	lines, _ := ReadLines("move_tables/rook.txt")
	for _, line := range lines {
		split := strings.Split(line, ";")
		pos, _ := strconv.Atoi(split[0])
		key, _ := strconv.Atoi(split[1])
		move, _ := strconv.ParseUint(split[2], 10, 64)
		mapping[MagicKey{Square: pos, Key: key}] = move
	}
	return mapping
}

// Static moves

func KingMovesMapping() map[int][]Move {
	return generateStaticMoves(kingOffsets, 0, 64)
}

func KnightMovesMapping() map[int][]Move {
	return generateStaticMoves(knightOffsets, 0, 64)
}

func WhitePawnStraightMovesMapping() map[int][]Move {
	nonPromotingMoves := generateStaticMoves(whitePawnStraightOffsets, 8, 48)
	promotingMoves := generateStaticMovesWithFlags(whitePawnStraightOffsets, 48, 56, promotionFlags)
	return mergeMoveMaps(nonPromotingMoves, promotingMoves, 8, 48)
}

func WhitePawnDoubleMovesMapping() map[int][]Move {
	return generateStaticMoves(whitePawnDoubleOffsets, 8, 16)
}

func WhitePawnAttackMovesMapping() map[int][]Move {
	nonPromotingMoves := generateStaticMoves(whitePawnAttackOffsets, 8, 48)
	promotingMoves := generateStaticMovesWithFlags(whitePawnAttackOffsets, 48, 56, promotionFlags)
	return mergeMoveMaps(nonPromotingMoves, promotingMoves, 8, 56)
}

func BlackPawnDoubleMovesMapping() map[int][]Move {
	return generateStaticMoves(blackPawnDoubleOffsets, 48, 56)
}

func BlackPawnStraightMovesMapping() map[int][]Move {
	nonPromotingMoves := generateStaticMoves(blackPawnStraightOffsets, 16, 56)
	promotingMoves := generateStaticMoves(blackPawnStraightOffsets, 8, 16)
	return mergeMoveMaps(nonPromotingMoves, promotingMoves, 8, 56)
}

func BlackPawnAttackMovesMapping() map[int][]Move {
	nonPromotingMoves := generateStaticMoves(blackPawnAttackOffsets, 16, 56)
	promotingMoves := generateStaticMoves(blackPawnAttackOffsets, 8, 16)
	return mergeMoveMaps(nonPromotingMoves, promotingMoves, 8, 56)
}

func generateStaticMoves(offsets [][2]int, from int, to int) map[int][]Move {
	standardFlags := []int{0}
	return generateStaticMovesWithFlags(offsets, from, to, standardFlags)
}

func generateStaticMovesWithFlags(offsets [][2]int, from int, to int, flags []int) map[int][]Move {
	moves := make(map[int][]Move)
	var currentMoves []Move

	for origin := from; origin < to; origin++ {
		for _, destination := range generateDestinations(origin, offsets) {
			for _, flag := range flags {
				destinationBits := uint16(destination)
				originBits := uint16(origin) << 6
				flagBits := uint16(flag) << 12

				move := Move{destinationBits | originBits | flagBits}
				currentMoves = append(currentMoves, move)
			}
		}
		moves[origin] = currentMoves
		currentMoves = nil
	}

	return moves
}

func generateDestinations(i int, offsets [][2]int) []int {
	baseRank, baseFile := IndexToCartesian(i)

	var indices []int

	for _, offset := range offsets {
		rank := baseRank + offset[0]
		file := baseFile + offset[1]
		if rank >= 1 && rank <= 8 && file >= 1 && file <= 8 {
			indices = append(indices, CartesianToIndex(rank, file))
		}
	}

	return indices
}

// Creates a mapping from index to the uint64 one-hot representation of the number.
func PosToBitBoard() map[int]uint64 {
	var mapping = make(map[int]uint64)
	for i := 0; i < 64; i++ {
		mapping[i] = uint64(math.Pow(2, float64(i)))
	}
	return mapping
}

func mergeMoveMaps(a map[int][]Move, b map[int][]Move, from int, to int) map[int][]Move {
	merged := make(map[int][]Move)
	for i := from; i < to; i++ {
		merged[i] = append(a[i], b[i]...)
	}
	return merged
}

// Attack masks
func KnightMasks() map[int]uint64 {
	return generateMaskMapping(knightOffsets)
}

func WhitePawnAttackMasks() map[int]uint64 {
	return generateMaskMapping(whitePawnAttackOffsets)
}

func BlackPawnAttackMasks() map[int]uint64 {
	return generateMaskMapping(blackPawnAttackOffsets)
}

func KingMasks() map[int]uint64 {
	return generateMaskMapping(kingOffsets)
}

func generateMaskMapping(offsets [][2]int) map[int]uint64 {
	mapping := make(map[int]uint64)
	for i := 0; i < 64; i++ {
		var mask uint64
		for _, destination := range generateDestinations(i, offsets) {
			mask |= uint64(1) << destination
		}
		mapping[i] = mask
	}
	return mapping
}
