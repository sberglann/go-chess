package main

func KingMasks() {
	offsets := [][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}}
	GenerateStaticMasks(offsets, 0, 64)
}

func KnightMasks() {
	offsets := [][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
	GenerateStaticMasks(offsets, 0, 64)
}

func WhitePawnStraightMasks() {
	offsets := [][2]int{{1, 0}}
	GenerateStaticMasks(offsets, 8, 56)
}

func WhitePawnDoubleMasks() {
	offsets := [][2]int{{2, 0}}
	GenerateStaticMasks(offsets, 8, 16)
}

func WhitePawnAttackMasks() {
	offsets := [][2]int{{1, -1}, {1, 1}}
	GenerateStaticMasks(offsets, 8, 56)
}

func BlackPawnDoubleMasks() {
	offsets := [][2]int{{-2, 0}}
	GenerateStaticMasks(offsets, 48, 56)
}

func BlackPawnStraightMasks() {
	offsets := [][2]int{{-1, 0}}
	GenerateStaticMasks(offsets, 8, 56)
}

func BlackPawnAttackMasks() {
	offsets := [][2]int{{-1, -1}, {-1, 1}}
	GenerateStaticMasks(offsets, 8, 56)
}

func GenerateStaticMasks(offsets [][2]int, from int, to int) {
	generate := func(i int) {
		baseRank, baseFile := IndexToCartesian(i)

		var indices []int
		var board uint64

		for _, offset := range offsets {
			rank := baseRank + offset[0]
			file := baseFile + offset[1]
			if rank >= 1 && rank <= 8 && file >= 1 && file <= 8 {
				indices = append(indices, CartesianToIndex(rank, file))
			}
		}

		for _, index := range indices {
			board |= 1 << index
		}
		Pretty64(board)
	}

	for i := from; i < to; i++ {
		generate(i)

	}
}
