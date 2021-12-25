package main

import (
	"fmt"
	"strconv"
)

func Perft() {
	var previousStates []BitBoard
	previousStates = append(previousStates, StartBoard)
	for depth := 0; depth < 6; depth++ {
		nextStates := perftStep(previousStates)
		fmt.Println("Depth " + strconv.Itoa(depth+1) + ": " + strconv.Itoa(len(nextStates)))
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
