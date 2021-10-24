package main

// https://www.chessprogramming.org/index.php?title=Looking_for_Magics

func RookMask(pos int) uint64 {
	result := uint64(0)
	rank := pos / 8
	file := pos % 8

	for r := rank + 1; r <= 6; r++ {
		result |= uint64(1) << (file + r*8)
	}
	for r := rank - 1; r >= 1; r-- {
		result |= uint64(1) << (file + r*8)
	}
	for f := file + 1; f <= 6; f++ {
		result |= uint64(1) << (f + rank*8)
	}
	for f := file - 1; f >= 1; f-- {
		result |= uint64(1) << (f + rank*8)
	}
	return result
}

func BishopMask(pos int) uint64 {
	result := uint64(0)
	rank := pos / 8
	file := pos % 8

	for r, f := rank+1, file+1; r <= 6 && f <= 6; r, f = r+1, f+1 {
		result |= uint64(1) << (f + r*8)
	}
	for r, f := rank+1, file-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
		result |= uint64(1) << (f + r*8)
	}
	for r, f := rank-1, file+1; r >= 1 && f <= 6; r, f = r-1, f+1 {
		result |= uint64(1) << (f + r*8)
	}
	for r, f := rank-1, file-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
		result |= uint64(1) << (f + r*8)
	}
	return result
}
func RookAttacks(pos int, block uint64) uint64 {
	result := uint64(0)
	rank := pos / 8
	file := pos % 8
	for r := rank + 1; r <= 7; r++ {
		result |= uint64(1) << (file + r*8)
		if block&(uint64(1)<<(file+r*8)) > 0 {
			break
		}
	}
	for r := rank - 1; r >= 0; r-- {
		result |= uint64(1) << (file + r*8)
		if block&(uint64(1)<<(file+r*8)) > 0 {
			break
		}
	}
	for f := file + 1; f <= 7; f++ {
		result |= uint64(1) << (f + rank*8)
		if block&(uint64(1)<<(f+rank*8)) > 0 {
			break
		}
	}
	for f := file - 1; f >= 0; f-- {
		result |= uint64(1) << (f + rank*8)
		if block&(uint64(1)<<(f+rank*8)) > 0 {
			break
		}
	}

	return result
}

func BishopAttacks(pos int, block uint64) uint64 {
	result := uint64(0)
	rank := pos / 8
	file := pos % 8

	for r, f := rank+1, file+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		result |= uint64(1) << (f + r*8)
		if block&(uint64(1)<<(f+r*8)) > 0 {
			break
		}
	}
	for r, f := rank+1, file-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		result |= uint64(1) << (f + r*8)
		if block&(uint64(1)<<(f+r*8)) > 0 {
			break
		}
	}
	for r, f := rank-1, file+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		result |= uint64(1) << (f + r*8)
		if block&(uint64(1)<<(f+r*8)) > 0 {
			break
		}
	}
	for r, f := rank-1, file-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		result |= uint64(1) << (f + r*8)
		if block&(uint64(1)<<(f+r*8)) > 0 {
			break
		}
	}
	return result
}

func MagicTransformation(block uint64, magic uint64, bits int) uint64 {
	return (block * magic) >> (64 - bits)
}

/*
func FindMagic(pos int, bits int, piece Piece) {
	var mask uint64
	var blocker [4096]uint64
	var attacks [4096]uint64
	if piece == Bishop {
		mask = BishopMask(pos)
	} else {
		mask = RookMask(pos)
	}
	n := NumberOfOnes(mask)

	for i := 0; i < (1 << n); i++ {
		blocker[i] = Carte
	}

}
*/
