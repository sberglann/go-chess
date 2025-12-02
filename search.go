package main

import (
	"math/rand"
	"sync"
	"time"
)

type EvaluatedBoard struct {
	board BitBoard
	eval  float64
}

const maxDepth = 5
const deterministic = true
const randomRange = 0
const maxRoutines = 16

var transpositionTable = newTranspositionTable()

func BestMoveWithoutTimeLimit(board *BitBoard) EvaluatedBoard {
	return BestMove(board, 5*time.Second)
}

func BestMove(board *BitBoard, timeLimit time.Duration) EvaluatedBoard {
	var bestMove BitBoard
	var bestEval float64

	transpositionTable.Clear()

	legalMoves, numMoves := GenerateLegalStates(board)
	if numMoves == 0 {
		return EvaluatedBoard{bestMove, bestEval}
	}

	isWhite := board.Turn() == White
	isWhiteTurn := !isWhite // The child boards have the opposite turn
	
	
	// Initialize move evaluations
	moveEvals := make([]EvaluatedBoard, numMoves)
	for i := range numMoves {
		moveEvals[i] = EvaluatedBoard{legalMoves[i], 0.0}
	}
	
	// Create persistent worker pool - workers live for all iterations
	workQueue := make(chan EvaluatedBoard, numMoves)
	resultQueue := make(chan EvaluatedBoard, numMoves)
	depthQueue := make(chan int8, numMoves)
	
	var workerWg sync.WaitGroup
	workerWg.Add(maxRoutines)
	for range maxRoutines {
		go func() {
			defer workerWg.Done()
			for board := range workQueue {
				depth := <-depthQueue
				eval := minimax(&board.board, 1, isWhiteTurn, -1000.0, 1000.0, depth)
				resultQueue <- EvaluatedBoard{
					board: board.board,
					eval: eval * randomFactor(),
				}
			}
		}()
	}
	
	start := time.Now()
	currentMaxDepth := int8(1)

	for time.Since(start) < timeLimit && currentMaxDepth <= maxDepth {
		var evals = make([]EvaluatedBoard, numMoves)
		
		// Send all work items to queue
		for i := range numMoves {
			workQueue <- moveEvals[i]
			depthQueue <- currentMaxDepth
		}
		
		// Collect all results and match by board hash
		results := make(map[uint64]float64)
		for range numMoves {
			result := <-resultQueue
			results[result.board.Hash()] = result.eval
		}
		
		// Map results back to evals array
		for i := range numMoves {
			boardHash := moveEvals[i].board.Hash()
			evals[i] = EvaluatedBoard{
				board: moveEvals[i].board,
				eval: results[boardHash],
			}
		}
		
		// Update move evaluations and find best move
		moveEvals = evals
		bestEval = -1000.0
		if !isWhite {
			bestEval = 1000.0
		}
		
		for _, evaledBoard := range evals {
			isBetter := (isWhite && evaledBoard.eval > bestEval) || (!isWhite && evaledBoard.eval < bestEval)
			if isBetter {
				bestEval = evaledBoard.eval
				bestMove = evaledBoard.board
			}
		}
		
		currentMaxDepth++
	}
	
	// Cleanup: close work queue and wait for workers to finish
	close(workQueue)
	close(depthQueue)
	workerWg.Wait()
	close(resultQueue)

	return EvaluatedBoard{bestMove, bestEval}
}

func minimax(board *BitBoard, depth int8, isWhite bool, alpha float64, beta float64, currentMaxDepth int8) float64 {
	if depth >= currentMaxDepth {
		eval := Eval(board)
		return eval
	}

	ttResult := transpositionTable.getUpperAndLower(board, depth)
	if ttResult.upper < alpha {
		return ttResult.upper
	}
	if ttResult.lower >= beta {
		return ttResult.lower
	}

	childrenArray, numChildren := GenerateLegalStates(board)
	if numChildren == 0 {
		kingPos, _ := PopFistBit(board.KingBB)
		isChecked := isChecked(board, kingPos, board.Turn())
		eval := 0.0
		if isChecked && isWhite {
			eval = 1000.0
		} else if isChecked && !isWhite {
			eval = -1000.0
		} else {
			eval = 0.0
		}
		return eval
	}

	children := childrenArray[:numChildren]
	
	var bestEval float64
	
	if isWhite {
		bestEval = -1000.0
		for _, child := range children {
			currentEval := minimax(&child, depth+1, false, alpha, beta, currentMaxDepth)
			if currentEval > bestEval {
				bestEval = currentEval
			}
			if bestEval > alpha {
				alpha = bestEval
			}
			if beta <= alpha {
				transpositionTable.setLower(board, bestEval, depth)
				break
			}
		}
	} else {
		bestEval = 1000.0
		for _, child := range children {
			currentEval := minimax(&child, depth+1, true, alpha, beta, currentMaxDepth)
			if currentEval < bestEval {
				bestEval = currentEval
			}
			if bestEval < beta {
				beta = bestEval
			}
			if beta <= alpha {
				transpositionTable.setUpper(board, bestEval, depth)
				break
			}
		}
	}
	return bestEval
}

func randomFactor() float64 {
	if deterministic {
		return 1.0
	} else {
		return (1.0 - randomRange) + rand.Float64()*randomRange
	}
}
