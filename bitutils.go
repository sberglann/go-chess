package main

var debruijnIndices = [64]int{
	0, 1, 48, 2, 57, 49, 28, 3,
	61, 58, 50, 42, 38, 29, 17, 4,
	62, 55, 59, 36, 53, 51, 43, 22,
	45, 39, 33, 30, 24, 18, 12, 5,
	63, 47, 56, 27, 60, 41, 37, 16,
	54, 35, 52, 21, 44, 32, 23, 11,
	46, 26, 40, 15, 34, 20, 31, 10,
	25, 14, 19, 9, 13, 8, 7, 6,
}
var debruijn64 = uint64(0x03f79d71b4cb0a89)

// Gets the bit at position pos, an integer in [0, 63]
func BitAt64(i *uint64, pos int) int {
	return int((*i >> uint64(pos)) & 1)
}

// Gets the bit at position pos, an integer in [0, 63]
func BitAt32(i *uint32, pos int) int {
	return int((*i >> uint32(pos)) & 1)
}

func LSB(i uint64) int {
	if i > 0 {
		return debruijnIndices[((i&-i)*debruijn64)>>58]
	} else {
		return -1
	}
}

func PopFistBit(i uint64) (int, uint64) {
	if i > 0 {
		pos := LSB(i)
		i &^= 1 << pos
		return pos, i
	} else {
		return -1, i
	}
}

func NumberOfOnes(n uint64) int {
	count := uint64(0)
	for n != 0 {
		count += n & uint64(1)
		n >>= 1
	}
	return int(count)
}

func IndexToUint64(pos int, bits int, m uint64) uint64 {
	result := uint64(0)
	for i := 0; i < bits; i++ {
		j, popped := PopFistBit(m)
		m = popped
		if pos & (1 << i) > 0 && j >= 0 {
			result |= 1 << j
		}
	}
	return result
}
