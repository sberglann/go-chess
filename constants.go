package main

var FileToLetter = map[int]string{
	1: "a",
	2: "b",
	3: "c",
	4: "d",
	5: "e",
	6: "f",
	7: "g",
	8: "h",
}

var IndexToAlgebraic = map[int]string{
	0: "a1", 1: "b1", 2: "c1", 3: "d1", 4: "e1", 5: "f1", 6: "g1", 7: "h1",
	8: "a2", 9: "b2", 10: "c2", 11: "d2", 12: "e2", 13: "f2", 14: "g2", 15: "h2",
	16: "a3", 17: "b3", 18: "c3", 19: "d3", 20: "e3", 21: "f3", 22: "g3", 23: "h3",
	24: "a4", 25: "b4", 26: "c4", 27: "d4", 28: "e4", 29: "f4", 30: "g4", 31: "h4",
	32: "a5", 33: "b5", 34: "c5", 35: "d5", 36: "e5", 37: "f5", 38: "g5", 39: "h5",
	40: "a6", 41: "b6", 42: "c6", 43: "d6", 44: "e6", 45: "f6", 46: "g6", 47: "h6",
	48: "a7", 49: "b7", 50: "c7", 51: "d7", 52: "e7", 53: "f7", 54: "g7", 55: "h7",
	56: "a8", 57: "b8", 58: "c8", 59: "d8", 60: "e8", 61: "f8", 62: "g8", 63: "h8",
}

// https://github.com/official-stockfish/Stockfish/blob/master/src/psqt.cpp
var psqPawn = [][]int{
	{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
	{2, -8}, {4, -6}, {11, 9}, {18, 5}, {16, 16}, {21, 6}, {9, -6}, {-3, -18},
	{-9, -9}, {-15, -7}, {11, -10}, {15, 5}, {31, 2}, {23, 3}, {6, -8}, {-20, -5},
	{-3, 7}, {-20, 1}, {8, -8}, {19, -2}, {39, -14}, {17, -13}, {2, -11}, {-5, -6},
	{11, 12}, {-4, 6}, {-11, 2}, {2, -6}, {11, -5}, {0, -4}, {-12, 14}, {5, 9},
	{3, 27}, {-11, 18}, {-6, 19}, {22, 29}, {-8, 30}, {-5, 9}, {-14, 8}, {-11, 14},
	{-7, -1}, {6, -14}, {-2, 13}, {-11, 22}, {4, 24}, {-14, 17}, {10, 7}, {-9, 7},
	{900, 900}, {900, 900}, {900, 900}, {900, 900}, {900, 900}, {900, 900}, {900, 900}, {900, 900},
}

var psqKnight = [][]int{
	{-175, -96}, {-92, -65}, {-74, -49}, {-73, -21}, {-73, -21}, {-74, -49}, {-92, -65}, {-175, -96},
	{-77, -67}, {-41, -54}, {-27, -18}, {-15, 8}, {-15, 8}, {-27, -18}, {-41, -54}, {-77, -67},
	{-61, -40}, {-17, -27}, {6, -8}, {12, 29}, {12, 29}, {6, -8}, {-17, -27}, {-61, -40},
	{-35, -35}, {8, -2}, {40, 13}, {49, 28}, {49, 28}, {40, 13}, {8, -2}, {-35, -35},
	{-34, -45}, {13, -16}, {44, 9}, {51, 39}, {51, 39}, {44, 9}, {13, -16}, {-34, -45},
	{-9, -51}, {22, -44}, {58, -16}, {53, 17}, {53, 17}, {58, -16}, {22, -44}, {-9, -51},
	{-67, -69}, {-27, -50}, {4, -51}, {37, 12}, {37, 12}, {4, -51}, {-27, -50}, {-67, -69},
	{-201, -100}, {-83, -88}, {-56, -56}, {-26, -17}, {-26, -17}, {-56, -56}, {-83, -88}, {-201, -100},
}

var psqBishop = [][]int{
	{-37, -40}, {-4, -21}, {-6, -26}, {-16, -8}, {-16, -8}, {-6, -26}, {-4, -21}, {-37, -40},
	{-11, -26}, {6, -9}, {13, -12}, {3, 1}, {3, 1}, {13, -12}, {6, -9}, {-11, -26},
	{-5, -11}, {15, -1}, {-4, -1}, {12, 7}, {12, 7}, {-4, -1}, {15, -1}, {-5, -11},
	{-4, -14}, {8, -4}, {18, 0}, {27, 12}, {27, 12}, {18, 0}, {8, -4}, {-4, -14},
	{-8, -12}, {20, -1}, {15, -10}, {22, 11}, {22, 11}, {15, -10}, {20, -1}, {-8, -12},
	{-11, -21}, {4, 4}, {1, 3}, {8, 4}, {8, 4}, {1, 3}, {4, 4}, {-11, -21},
	{-12, -22}, {-10, -14}, {4, -1}, {0, 1}, {0, 1}, {4, -1}, {-10, -14}, {-12, -22},
	{-34, -32}, {1, -29}, {-10, -26}, {-16, -17}, {16, -17}, {-10, -26}, {1, -29}, {-34, -32},
}

var psqRook = [][]int{
	{-31, -9}, {-20, -13}, {-14, -10}, {-5, -9}, {-5, -9}, {-14, -10}, {-20, -13}, {-31, -9},
	{-21, -12}, {-13, -9}, {-8, -1}, {6, -2}, {6, -2}, {-8, -1}, {-13, -9}, {-21, -12},
	{-25, 6}, {-11, -8}, {-1, -2}, {3, -6}, {3, -6}, {-1, -2}, {-11, -8}, {-25, 6},
	{-13, -6}, {-5, 1}, {-4, -9}, {-6, 7}, {-6, 7}, {-4, -9}, {-5, 1}, {-13, -6},
	{-27, -5}, {-15, 8}, {-4, 7}, {3, -6}, {3, -6}, {-4, 7}, {-15, 8}, {-27, -5},
	{-22, 6}, {-2, 1}, {6, -7}, {12, 10}, {12, 10}, {6, -7}, {-2, 1}, {-22, 6},
	{-2, 4}, {12, 5}, {16, 20}, {18, -5}, {18, -5}, {16, 20}, {12, 5}, {-2, 4},
	{-17, 18}, {-19, 0}, {-1, 19}, {9, 13}, {9, 13}, {-1, 19}, {-19, 0}, {-17, 18},
}

var psqQueen = [][]int{
	{3, -69}, {-5, -57}, {-5, -47}, {4, -26}, {4, -26}, {-5, -47}, {-5, -57}, {3, -69},
	{-3, -54}, {5, -31}, {8, -22}, {12, -4}, {12, -4}, {8, -22}, {5, -31}, {-3, -54},
	{-3, -39}, {6, -18}, {13, -9}, {7, 3}, {7, 3}, {13, -9}, {6, -18}, {-3, -39},
	{4, -23}, {5, -3}, {9, 13}, {8, 24}, {8, 24}, {9, 13}, {5, -3}, {4, -23},
	{0, -29}, {14, -6}, {12, 9}, {5, 21}, {5, 21}, {12, 9}, {14, -6}, {0, -29},
	{-4, -38}, {10, -18}, {6, -11}, {8, 1}, {8, 1}, {6, -11}, {10, -18}, {-4, -38},
	{-5, -50}, {6, -27}, {10, -24}, {8, -8}, {8, -8}, {10, -24}, {6, -27}, {-5, -50},
	{-2, -74}, {-2, -52}, {1, -43}, {-2, -34}, {-2, -34}, {1, -43}, {-2, -52}, {-2, -74},
}

var psqKing = [][]int{
	{271, 1}, {327, 45}, {271, 85}, {198, 76}, {198, 76}, {271, 85}, {327, 45}, {271, 1},
	{278, 53}, {303, 100}, {234, 133}, {179, 135}, {179, 135}, {234, 133}, {303, 100}, {278, 53},
	{195, 88}, {258, 130}, {169, 169}, {120, 175}, {120, 175}, {169, 169}, {258, 130}, {195, 88},
	{164, 103}, {190, 156}, {138, 172}, {98, 172}, {98, 172}, {138, 172}, {190, 156}, {164, 103},
	{154, 96}, {179, 166}, {105, 199}, {70, 199}, {70, 199}, {105, 199}, {179, 166}, {154, 96},
	{123, 92}, {145, 172}, {81, 184}, {31, 191}, {31, 191}, {81, 184}, {145, 172}, {123, 92},
	{88, 47}, {120, 121}, {65, 116}, {33, 131}, {33, 131}, {65, 116}, {120, 121}, {88, 47},
	{59, 11}, {89, 59}, {45, 73}, {-1, 78}, {-1, 78}, {45, 73}, {89, 59}, {59, 11},
}

var OneHot32 = [32]uint32{
	0b00000000000000000000000000000001,
	0b00000000000000000000000000000010,
	0b00000000000000000000000000000100,
	0b00000000000000000000000000001000,
	0b00000000000000000000000000010000,
	0b00000000000000000000000000100000,
	0b00000000000000000000000001000000,
	0b00000000000000000000000010000000,
	0b00000000000000000000000100000000,
	0b00000000000000000000001000000000,
	0b00000000000000000000010000000000,
	0b00000000000000000000100000000000,
	0b00000000000000000001000000000000,
	0b00000000000000000010000000000000,
	0b00000000000000000100000000000000,
	0b00000000000000001000000000000000,
	0b00000000000000010000000000000000,
	0b00000000000000100000000000000000,
	0b00000000000001000000000000000000,
	0b00000000000010000000000000000000,
	0b00000000000100000000000000000000,
	0b00000000001000000000000000000000,
	0b00000000010000000000000000000000,
	0b00000000100000000000000000000000,
	0b00000001000000000000000000000000,
	0b00000010000000000000000000000000,
	0b00000100000000000000000000000000,
	0b00001000000000000000000000000000,
	0b00010000000000000000000000000000,
	0b00100000000000000000000000000000,
	0b01000000000000000000000000000000,
	0b10000000000000000000000000000000,
}
