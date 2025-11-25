package main

import (
	"math"
	"math/rand"
	"sync"
)

type EvaluatedBoard struct {
	board BitBoard
	eval  float64
}

const maxDepth = 5
const deterministic = true
const randomRange = 0

func BestMove(board *BitBoard) EvaluatedBoard {
	var bestMove BitBoard
	var bestEval float64

	legalMoves, numMoves := GenerateLegalStates(board)
	maxRoutines := 16
	guard := make(chan struct{}, maxRoutines)
	var evals = make([]EvaluatedBoard, numMoves)
	var wg sync.WaitGroup
	wg.Add(numMoves)

	if board.Turn() == White {
		for i := 0; i < numMoves; i++ {
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				eval := minimax(&b, 1, false, -1000.0, 1000.0)
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, legalMoves[i])
		}
		wg.Wait()
		maxEval := -1000.0
		for _, evaledBoard := range evals {
			if evaledBoard.eval > maxEval {
				maxEval = evaledBoard.eval
				bestEval = maxEval
				bestMove = evaledBoard.board
			}
		}

	} else {
		for i := 0; i < numMoves; i++ {
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				eval := minimax(&b, 1, true, -1000.0, 1000.0)
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, legalMoves[i])
		}
		wg.Wait()
		minEval := 1000.0
		for _, evaledBoard := range evals {
			if evaledBoard.eval < minEval {
				minEval = evaledBoard.eval
				bestEval = minEval
				bestMove = evaledBoard.board
			}
		}
	}

	return EvaluatedBoard{bestMove, bestEval}
}

func minimax(board *BitBoard, depth int, isWhite bool, alpha float64, beta float64) float64 {
	if depth == maxDepth {
		return Eval(board)
	}
	children, numChildren := GenerateLegalStates(board)
	var bestEval float64
	if numChildren == 0 {
		return Eval(board)
	} else if isWhite {
		bestEval = -1000.0
		i := 0
		child := children[i]
		for child.WhiteBB > 0 {
			currentEval := minimax(&child, depth+1, false, alpha, beta)
			bestEval = math.Max(bestEval, currentEval)
			alpha = math.Max(alpha, bestEval)
			if beta <= alpha {
				break
			}
			i++
			child = children[i]
		}
	} else {
		bestEval = 1000.0
		i := 0
		child := children[i]
		for child.WhiteBB > 0 {
			currentEval := minimax(&child, depth+1, true, alpha, beta)
			bestEval = math.Min(bestEval, currentEval)
			beta = math.Min(beta, bestEval)
			if beta <= alpha {
				break
			}
			i++
			child = children[i]
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
