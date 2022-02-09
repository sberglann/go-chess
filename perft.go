package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var EnPassantCounter int
var CaptureCounter int
var CheckMateCounter int

func Perft() {
	start := time.Now()
	var previousStates []BitBoard
	b0 := StartBoard
	//b0 := BoardFromFEN("8/8/8/2k5/2pP4/8/B7/4K3 b - d3 0 3")
	b0.PrettyBoard()
	previousStates = append(previousStates, b0)
	for depth := 0; depth < 8; depth++ {
		nextStates := perftStepPar(previousStates)
		fmt.Println("Depth " + strconv.Itoa(depth+1) + ": " + strconv.Itoa(len(nextStates)) + " " + strconv.Itoa(EnPassantCounter) + " " + strconv.Itoa(CheckMateCounter) + " " + strconv.Itoa(CaptureCounter))
		previousStates = nextStates
	}
	elapsed := time.Since(start)
	fmt.Printf("Time: %s", elapsed)
}

func perftStep(previousStates []BitBoard) []BitBoard {
	var nextStates []BitBoard
	resetCounters()
	for _, b := range previousStates {
		nextStates = append(nextStates, GenerateLegalMoves(b)...)
	}
	return nextStates
}

func perftStepPar(previousStates []BitBoard) []BitBoard {
	var nextStates = make([][]BitBoard, len(previousStates))

	var wg sync.WaitGroup
	wg.Add(len(previousStates))
	resetCounters()
	for i, b := range previousStates {
		go func(j int, b BitBoard) {
			defer wg.Done()
			nextStates[j] = GenerateLegalMoves(b)
		}(i, b)
	}

	wg.Wait()
	var nextStatesFlat []BitBoard
	for _, ns := range nextStates {
		nextStatesFlat = append(nextStatesFlat, ns...)
	}

	return nextStatesFlat
}

func resetCounters() {
	EnPassantCounter = 0
	CaptureCounter = 0
	CheckMateCounter = 0
}
