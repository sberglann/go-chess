package main

import (
	"fmt"
	"strconv"
)

var EnPassantCounter int
var CaptureCounter int
var CheckMateCounter int

func Perft() {
	var previousStates []BitBoard
	//b0 := StartBoard
	b0 := BoardFromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0")
	b0.PrettyBoard()
	previousStates = append(previousStates, b0)
	for depth := 0; depth < 6; depth++ {
		nextStates := perftStep(previousStates)
		fmt.Println("Depth " + strconv.Itoa(depth+1) + ": " + strconv.Itoa(len(nextStates)) + " " + strconv.Itoa(EnPassantCounter) + " " + strconv.Itoa(CheckMateCounter) + " " + strconv.Itoa(CaptureCounter))
		previousStates = nextStates
	}
}

func perftStep(previousStates []BitBoard) []BitBoard {
	var nextStates []BitBoard
	resetCounters()
	for _, b := range previousStates {
		nextStates = append(nextStates, GenerateLegalMoves(b)...)
	}
	return nextStates
}

func resetCounters() {
	EnPassantCounter = 0
	CaptureCounter = 0
	CheckMateCounter = 0
}
