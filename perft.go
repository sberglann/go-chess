package main

import (
	"fmt"
	"os"
	"runtime/pprof"
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

func PerformanceTest(cpuprofile string) {
	var f *os.File
	if cpuprofile != "" {
		var err error
		f, err = os.Create(cpuprofile)
		if err != nil {
			return
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	positions := []string{
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		"rnbqkbnr/ppp1pppp/8/3p4/4P3/2N5/PPPP1PPP/R1BQKBNR b KQkq - 1 2",
		"rnbqkbnr/ppp1pppp/8/8/4p3/2N5/PPPP1PPP/R1BQKBNR w KQkq - 0 3",
		"rnbqkbnr/p1p1pppp/1p6/8/4N3/5N2/PPPP1PPP/R1BQKB1R b KQkq - 1 4",
		"rn1qkbnr/pbp1pppp/1p6/8/4N3/3B1N2/PPPP1PPP/R1BQK2R b KQkq - 3 5",
		"rn1qkbnr/p1p1pppp/1p6/8/4B3/5N2/PPPP1PPP/R1BQK2R b KQkq - 0 6",
		"rn1qkbnr/p1p1pp1p/1p4p1/8/4B3/5N2/PPPP1PPP/R1BQ1RK1 b kq - 1 7",
		"rn1qk1nr/p1p1ppbp/1p4p1/8/3PB3/5N2/PPP2PPP/R1BQ1RK1 b kq - 0 8",
		"Bn1qk2r/p1p1ppbp/1p3np1/8/3P4/5N2/PPP2PPP/R1BQ1RK1 b k - 0 9",
		"Bn1qk2r/p3ppbp/1pp2np1/8/3P1B2/5N2/PPP2PPP/R2Q1RK1 b k - 1 10",
		"Bn1q1rk1/p3ppbp/1pp2np1/8/3P1B2/5N2/PPP2PPP/R2QR1K1 b - - 3 11",
		"Bn1q1rk1/4ppbp/1pp2np1/p7/3P1B2/2P2N2/PP3PPP/R2QR1K1 b - - 0 12",
		"BB1q1rk1/5pbp/1pp1pnp1/p7/3P4/2P2N2/PP3PPP/R2QR1K1 b - - 0 13",
		"1q3rk1/5pbp/1pB1pnp1/p7/3P4/2P2N2/PP3PPP/R2QR1K1 b - - 0 14",
		"1q3rk1/5pbp/2B1pnp1/pp2N3/3P4/2P5/PP3PPP/R2QR1K1 b - - 1 15",
		"1q3rk1/5pb1/2B1pnp1/pp2N2p/3P4/2P4P/PP3PP1/R2QR1K1 b - - 0 16",
		"1q3rk1/5pb1/2B1pnp1/p3N2p/1pPP4/7P/PP3PP1/R2QR1K1 b - - 0 17",
		"5rk1/5pb1/2Bqpnp1/p3N2p/1pPP4/7P/PP1Q1PP1/R3R1K1 b - - 2 18",
		"5rk1/5pb1/2N1pnp1/p6p/1pPP4/7P/PP1Q1PP1/R3R1K1 b - - 0 19",
		"r5k1/5pb1/2N1pnp1/p5Qp/1pPP4/7P/PP3PP1/R3R1K1 b - - 2 20",
		"r5k1/3nQpb1/2N1p1p1/p6p/1pPP4/7P/PP3PP1/R3R1K1 b - - 4 21",
		"r5k1/3nQp2/4p1p1/p6p/1pPN4/7P/PP3PP1/R3R1K1 b - - 0 22",
		"6k1/r2nQp2/4p1p1/p1P4p/1p1N4/7P/PP3PP1/R3R1K1 b - - 0 23",
		"6k1/3Q1p2/r3p1p1/p1P4p/1p1N4/7P/PP3PP1/R3R1K1 b - - 0 24",
		"6k1/3Q1p2/r5p1/p1P1R2p/1p1N4/7P/PP3PP1/R5K1 b - - 0 25",
		"4R1k1/3Q1p2/r5p1/p1P4p/3N4/1p5P/PP3PP1/R5K1 b - - 1 26",
		"4R3/4Qpk1/r5p1/p1P4p/3N4/1p5P/PP3PP1/R5K1 b - - 3 27",
		"4R3/5pk1/r7/p1P3Qp/3N4/1p5P/PP3PP1/R5K1 b - - 0 28",
		"4R1Q1/5p1k/r7/p1P4p/3N4/1p5P/PP3PP1/R5K1 b - - 2 29",
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	start := time.Now()
	for i, position := range positions {
		moveStart := time.Now()
		board := BoardFromFEN(position)
		BestMove(board)

		moveElapsed := time.Since(moveStart)
		println("Position", i, "took", moveElapsed.Milliseconds(), "ms")
	}
	totalTime := time.Since(start)
	println("Total time:", totalTime.Milliseconds(), "ms")
}

func perftStep(previousStates []BitBoard) []BitBoard {
	var nextStates []BitBoard
	resetCounters()
	for _, b := range previousStates {
		states, _ := GenerateLegalStates(b)
		nextStates = append(nextStates, states[:]...)
	}
	return nextStates
}

func perftStepPar(previousStates []BitBoard) []BitBoard {
	var nextStates = make([][200]BitBoard, len(previousStates))

	maxRoutines := 16
	guard := make(chan struct{}, maxRoutines)

	var wg sync.WaitGroup
	wg.Add(len(previousStates))
	resetCounters()
	for i, b := range previousStates {
		guard <- struct{}{}
		go func(j int, b BitBoard) {
			defer wg.Done()
			nextStates[j], i = GenerateLegalStates(b)
			<-guard
		}(i, b)
	}

	wg.Wait()
	var nextStatesFlat []BitBoard
	for _, ns := range nextStates {
		nextStatesFlat = append(nextStatesFlat, ns[:]...)
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
