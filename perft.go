package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var EnPassantCounter int
var CaptureCounter int
var CheckMateCounter int

func Perft() {
	start := time.Now()
	// TODO: Investigate why black can't castle in step two of case 10.
	//TestSingle(1)
	TestAll()
	elapsed := time.Since(start)
	fmt.Printf("Time: %s\n\n", elapsed)
}

func NumMoves(fen string) {
	board := BoardFromFEN(fen)
	states1 := GenerateLegalMoves(board)
	states2 := perftStep(states1)
	println(len(states2) + 1)
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

	maxRoutines := 16
	guard := make(chan struct{}, maxRoutines)

	var wg sync.WaitGroup
	wg.Add(len(previousStates))
	resetCounters()
	for i, b := range previousStates {
		guard <- struct{}{}
		go func(j int, b BitBoard) {
			defer wg.Done()
			nextStates[j] = GenerateLegalMoves(b)
			<-guard
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

type PerftTestCase struct {
	id             int
	board          BitBoard
	expectedCounts []int
}

func TestSingle(id int) {
	testCase := extractTestCases()[id]
	testCase.board.PrettyBoard()
	assertTestCase(testCase, true)
}

func TestAll() {
	var passed int
	var failed int

	for _, testCase := range extractTestCases() {
		if assertTestCase(testCase, false) {
			passed += 1
		} else {
			failed += 1
		}
	}

	s := fmt.Sprintf("%.2f", 100*float64(passed)/float64(passed+failed))
	fmt.Println("Test run finished. Passed ", s, "%.")
}

func extractTestCases() []PerftTestCase {
	lines, _ := ReadLines("resources/perft_answers.csv")
	var cases []PerftTestCase
	for i, line := range lines {
		cases = append(cases, parseLine(i, line))
	}
	return cases
}

func parseLine(id int, line string) PerftTestCase {
	split := strings.Split(line, ",")
	fen, expectedCountsString := split[0], split[1:]

	board := BoardFromFEN(fen)
	var expectedCounts []int
	for _, c := range expectedCountsString {
		count, _ := strconv.Atoi(c)
		if count < 1000000 {
			expectedCounts = append(expectedCounts, count)
		}
	}
	return PerftTestCase{id, board, expectedCounts}
}

func assertTestCase(testCase PerftTestCase, printCounts bool) bool {
	var boards []BitBoard
	boards = append(boards, testCase.board)
	for step, expected := range testCase.expectedCounts {
		boards = perftStepPar(boards)
		if printCounts {
			fmt.Println(testCase.id, "-", step, ":", len(boards), "vs expected", expected)
		}
		if len(boards) != expected {
			fmt.Printf("Failed perft no. %03d at step %d\n", testCase.id, step)
			return false
		}
	}
	return true
}
