package main

import (
	"fmt"
	"strconv"
)

var EnPassantCounter int
var CheckMateCounter int

func Perft() {
	var previousStates []BitBoard
	b := BoardFromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0")
	previousStates = append(previousStates, b)
	for depth := 0; depth < 5; depth++ {
		nextStates := perftStep(previousStates)
		fmt.Println("Depth " + strconv.Itoa(depth+1) + ": " + strconv.Itoa(len(nextStates)) + " " + strconv.Itoa(EnPassantCounter) + " " + strconv.Itoa(CheckMateCounter))
		previousStates = nextStates
	}
}

func perftStep(previousStates []BitBoard) []BitBoard {
	var nextStates []BitBoard
	for _, b := range previousStates {
		//b.PrettyBoard()
		//fmt.Println("--------------------------")
		nextStates = append(nextStates, GenerateLegalMoves(b)...)
	}
	return nextStates
}
