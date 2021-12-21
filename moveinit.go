package main

import (
	"math"
	"strconv"
	"strings"
)

type MagicKey struct {
	Square, Key int
}

var promotionFlags = []int{4, 5, 6, 7}

func BishopMasks() map[int]uint64 {
	var mapping = make(map[int]uint64)
	for i := 0; i < 64; i++ {
		mapping[i] = GenerateBishopMask(i)
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

func RookMasks() map[int]uint64 {
	var mapping = make(map[int]uint64)
	for i := 0; i < 64; i++ {
		mapping[i] = GenerateRookMask(i)
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

func KingMasks() map[int][]Move {
	offsets := [][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}}
	return GenerateStaticMasks(offsets, 0, 64)
}

func KnightMasks() map[int][]Move {
	offsets := [][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
	return GenerateStaticMasks(offsets, 0, 64)
}

func WhitePawnStraightMasks() map[int][]Move {
	offsets := [][2]int{{1, 0}}
	nonPromotingMoves := GenerateStaticMasks(offsets, 8, 48)
	promotingMoves := GenerateStaticMasksWithFlags(offsets, 48, 56, promotionFlags)
	return mergeMoveMaps(nonPromotingMoves, promotingMoves, 8, 48)
}

func WhitePawnDoubleMasks() map[int][]Move {
	offsets := [][2]int{{2, 0}}
	return GenerateStaticMasks(offsets, 8, 16)
}

func WhitePawnAttackMasks() map[int][]Move {
	offsets := [][2]int{{1, -1}, {1, 1}}
	nonPromotingMoves := GenerateStaticMasks(offsets, 8, 48)
	promotingMoves := GenerateStaticMasksWithFlags(offsets, 48, 56, promotionFlags)
	return mergeMoveMaps(nonPromotingMoves, promotingMoves, 8, 56)
}

func BlackPawnDoubleMasks() map[int][]Move {
	offsets := [][2]int{{-2, 0}}
	return GenerateStaticMasks(offsets, 48, 56)
}

func BlackPawnStraightMasks() map[int][]Move {
	offsets := [][2]int{{-1, 0}}
	nonPromotingMoves := GenerateStaticMasks(offsets, 16, 56)
	promotingMoves := GenerateStaticMasks(offsets, 8, 16)
	return mergeMoveMaps(nonPromotingMoves, promotingMoves, 8, 56)
}

func BlackPawnAttackMasks() map[int][]Move {
	offsets := [][2]int{{-1, -1}, {-1, 1}}
	nonPromotingMoves := GenerateStaticMasks(offsets, 16, 56)
	promotingMoves := GenerateStaticMasks(offsets, 8, 16)
	return mergeMoveMaps(nonPromotingMoves, promotingMoves, 8, 56)
}

func GenerateStaticMasks(offsets [][2]int, from int, to int) map[int][]Move {
	standardFlags := []int{0}
	return GenerateStaticMasksWithFlags(offsets, from, to, standardFlags)
}

func GenerateStaticMasksWithFlags(offsets [][2]int, from int, to int, flags []int) map[int][]Move {
	generateDestinations := func(i int) []int {
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

	moves := make(map[int][]Move)
	var currentMoves []Move

	for origin := from; origin < to; origin++ {
		for _, destination := range generateDestinations(origin) {
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
