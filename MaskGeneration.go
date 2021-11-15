package main

var promotionFlags = []int{4, 5, 6, 7}

func KingMasks() map[int][]uint16 {
	var mapping map[int][]uint16
	offsets := [][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}}
	for i, mask := range GenerateStaticMasks(offsets, 0, 64) {
		mapping[i] = mask
	}
}

func KnightMasks() []uint16 {
	offsets := [][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
	return GenerateStaticMasks(offsets, 0, 64)
}

func WhitePawnStraightMasks() []uint16 {
	offsets := [][2]int{{1, 0}}
	nonPromotingMoves := GenerateStaticMasks(offsets, 8, 48)
	promotingMoves := GenerateStaticMasksWithFlags(offsets, 48, 56, promotionFlags)
	return append(nonPromotingMoves, promotingMoves...)
}

func WhitePawnDoubleMasks() []uint16 {
	offsets := [][2]int{{2, 0}}
	return GenerateStaticMasks(offsets, 8, 16)
}

func WhitePawnAttackMasks() []uint16 {
	offsets := [][2]int{{1, -1}, {1, 1}}
	nonPromotingMoves := GenerateStaticMasks(offsets, 8, 48)
	promotingMoves := GenerateStaticMasksWithFlags(offsets, 48, 56, promotionFlags)
	return append(nonPromotingMoves, promotingMoves...)
}

func BlackPawnDoubleMasks() []uint16 {
	offsets := [][2]int{{-2, 0}}
	return GenerateStaticMasks(offsets, 48, 56)
}

func BlackPawnStraightMasks() []uint16 {
	offsets := [][2]int{{-1, 0}}
	nonPromotingMoves := GenerateStaticMasks(offsets, 16, 56)
	promotingMoves := GenerateStaticMasks(offsets, 8, 16)
	return append(nonPromotingMoves, promotingMoves...)
}

func BlackPawnAttackMasks() []uint16 {
	offsets := [][2]int{{-1, -1}, {-1, 1}}
	nonPromotingMoves := GenerateStaticMasks(offsets, 16, 56)
	promotingMoves := GenerateStaticMasks(offsets, 8, 16)
	return append(nonPromotingMoves, promotingMoves...)
}

func GenerateStaticMasks(offsets [][2]int, from int, to int) []uint16 {
	standardFlags := []int{0}
	return GenerateStaticMasksWithFlags(offsets, from, to, standardFlags)
}

func GenerateStaticMasksWithFlags(offsets [][2]int, from int, to int, flags []int) []uint16 {
	generateEndPosition := func(i int) []int {
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

	var moves []uint16

	for start := from; start < to; start++ {
		for _, end := range generateEndPosition(start) {
			for _, flag := range flags {
				startBits := uint16(start)
				endBits := uint16(end) << 6
				flagBits := uint16(flag) << 12

				move := startBits | endBits | flagBits
				moves = append(moves, move)
			}
		}
	}

	return moves
}
