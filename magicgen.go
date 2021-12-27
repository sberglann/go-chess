package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

// https://www.chessprogramming.org/index.php?title=Looking_for_Magics
var RookBitsArray = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

var BishopBitsArray = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

func RandomUint64() uint64 {
	return uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
}

func RandomFewBits() uint64 {
	return RandomUint64() & RandomUint64() & RandomUint64()
}

// With includeEdges=false, the masks for the outermost squares will not be considered.
// This is to save key space when generatic magics.
func GenerateRookMask(pos int, includeEdges bool) uint64 {
	result := uint64(0)
	rank := pos / 8
	file := pos % 8
	var upperLimit, lowerLimit int
	if includeEdges {
		upperLimit = 7
		lowerLimit = 0
	} else {
		upperLimit = 6
		lowerLimit = 1
	}

	for r := rank + 1; r <= upperLimit; r++ {
		result |= uint64(1) << (file + r*8)
	}
	for r := rank - 1; r >= lowerLimit; r-- {
		result |= uint64(1) << (file + r*8)
	}
	for f := file + 1; f <= upperLimit; f++ {
		result |= uint64(1) << (f + rank*8)
	}
	for f := file - 1; f >= lowerLimit; f-- {
		result |= uint64(1) << (f + rank*8)
	}
	return result
}

// With includeEdges=false, the masks for the outermost squares will not be considered.
// This is to save key space when generatic magics.
func GenerateBishopMask(pos int, includeEdges bool) uint64 {
	result := uint64(0)
	rank := pos / 8
	file := pos % 8
	var upperLimit, lowerLimit int
	if includeEdges {
		upperLimit = 7
		lowerLimit = 0
	} else {
		upperLimit = 6
		lowerLimit = 1
	}
	for r, f := rank+1, file+1; r <= upperLimit && f <= upperLimit; r, f = r+1, f+1 {
		result |= uint64(1) << (f + r*8)
	}
	for r, f := rank+1, file-1; r <= upperLimit && f >= lowerLimit; r, f = r+1, f-1 {
		result |= uint64(1) << (f + r*8)
	}
	for r, f := rank-1, file+1; r >= lowerLimit && f <= upperLimit; r, f = r-1, f+1 {
		result |= uint64(1) << (f + r*8)
	}
	for r, f := rank-1, file-1; r >= lowerLimit && f >= lowerLimit; r, f = r-1, f-1 {
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

func FindMagic(square int, numBits int, piece Piece) uint64 {
	var mask uint64
	var fail bool
	var blocker, attacks, used [4096]uint64

	existing_entries := make(map[string]bool)

	if piece == Bishop {
		mask = GenerateBishopMask(square, false)
	} else {
		mask = GenerateRookMask(square, false)
	}
	n := NumberOfOnes(mask)

	for i := 0; i < (1 << n); i++ {
		blocker[i] = IndexToUint64(i, n, mask)
		if piece == Bishop {
			attacks[i] = BishopAttacks(square, blocker[i])
		} else {
			attacks[i] = RookAttacks(square, blocker[i])
		}
	}
	for k := 0; k < 100000000; k++ {
		magic := RandomFewBits()
		if NumberOfOnes((mask*magic)&uint64(0xFF00000000000000)) < 6 {
			continue
		}
		for i := 0; i < 4096; i++ {
			used[i] = 0
		}
		fail = false
		for i := 0; !fail && i < (1<<n); i++ {
			j := MagicTransformation(blocker[i], magic, numBits)
			if used[j] == 0 {
				used[j] = attacks[i]
			} else if used[j] != attacks[i] {
				fail = true
			}
		}
		if !fail {
			for _, blocker := range blocker {
				key := (blocker * magic) >> (64 - numBits)
				if piece == Rook {
					attack := RookAttacks(square, blocker)
					line := strconv.Itoa(square) + ";" + strconv.FormatUint(key, 10) + ";" + strconv.FormatUint(attack, 10) + "\n"
					if _, value := existing_entries[line]; !value {
						existing_entries[line] = true
						AppendToFile("move_tables/rook.txt", line)
					}
				} else {
					attack := BishopAttacks(square, blocker)
					line := strconv.Itoa(square) + ";" + strconv.FormatUint(key, 10) + ";" + strconv.FormatUint(attack, 10) + "\n"
					if _, value := existing_entries[line]; !value {
						existing_entries[line] = true
						AppendToFile("move_tables/bishop.txt", line)
					}
				}
			}

			if piece == Rook {
				AppendToFile("move_tables/magics_rook.txt", strconv.Itoa(square)+";"+strconv.FormatUint(magic, 10)+"\n")
			} else {
				AppendToFile("move_tables/magics_bishop.txt", strconv.Itoa(square)+";"+strconv.FormatUint(magic, 10)+"\n")
			}

			return magic
		}
	}
	fmt.Println("Magic generation failed")
	return uint64(0)
}

func GenerateMagics() {
	var bishopMagics [64]uint64
	DeleteFile("move_tables/magics_rook.txt")
	DeleteFile("move_tables/magics_bishop.txt")
	DeleteFile("move_tables/bishop.txt")
	DeleteFile("move_tables/rook.txt")

	pieces := [2]Piece{Rook, Bishop}

	for _, piece := range pieces {
		var bitCounts [64]int
		if piece == Rook {
			bitCounts = RookBitsArray
		} else {
			bitCounts = BishopBitsArray
		}
		for square := 0; square < 64; square++ {
			magic := FindMagic(square, bitCounts[square], piece)
			bishopMagics[square] = magic
			fmt.Println("Finding magic for ", piece, " at square ", square)
		}
	}
}
