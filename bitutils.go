package gochess

import "math/bits"

// Gets the bit at position pos, an integer in [0, 63]
func BitAt64(i uint64, pos int) int {
	return int((i >> uint64(pos)) & 1)
}

// Gets the bit at position pos, an integer in [0, 63]
func BitAt32(i uint32, pos int) bool {
	return OneHot32[pos]&i != 0
}

func LSB(i uint64) int {
	return bits.TrailingZeros64(i)
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
		if pos&(1<<i) > 0 && j >= 0 {
			result |= 1 << j
		}
	}
	return result
}

func MurmurHash(x uint64) uint64 {
	//https: //en.wikipedia.org/wiki/MurmurHash
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	x *= 0xc4ceb9fe1a85ec53
	x ^= x >> 33
	return x
}
