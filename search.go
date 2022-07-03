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

func BestMove(board BitBoard) EvaluatedBoard {
	var bestMove BitBoard
	var bestEval float64

	legalMoves := GenerateLegalStates(board)
	maxRoutines := 16
	guard := make(chan struct{}, maxRoutines)
	var evals = make([]EvaluatedBoard, len(legalMoves))
	var wg sync.WaitGroup
	wg.Add(len(legalMoves))

	if board.Turn() == White {
		for i, b := range legalMoves {
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				eval := minimax(b, 1, false, math.Inf(-1), math.Inf(1))
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, b)
		}
		wg.Wait()
		maxEval := math.Inf(-1)
		for _, evaledBoard := range evals {
			if evaledBoard.eval > maxEval {
				maxEval = evaledBoard.eval
				bestEval = maxEval
				bestMove = evaledBoard.board
			}
		}

	} else {
		for i, b := range legalMoves {
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				eval := minimax(b, 1, true, math.Inf(-1), math.Inf(1))
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, b)
		}
		wg.Wait()
		minEval := math.Inf(1)
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

func minimax(board BitBoard, depth int, isWhite bool, alpha float64, beta float64) float64 {
	children := GenerateLegalStates(board)
	var sign int
	if isWhite {
		sign = 1
	} else {
		sign = -1
	}
	if len(children) <= 0 {
		// Check mate
		return float64(sign * -1000.0)
	} else if depth >= maxDepth {
		maxEval := math.Inf(-1)
		for _, child := range children {
			eval := Eval(child)
			if eval > maxEval {
				maxEval = eval
			}
		}
		return maxEval
	}

	if isWhite {
		best := math.Inf(-1)
		for _, child := range children {
			value := minimax(child, depth+1, false, alpha, beta)
			best = math.Max(best, value)
			alpha = math.Max(alpha, best)
			if beta <= alpha {
				break
			}
		}
		return best
	} else {
		best := math.Inf(1)
		for _, child := range children {
			value := minimax(child, depth+1, true, alpha, beta)
			best = math.Min(best, value)
			beta = math.Min(beta, best)
			if beta <= alpha {
				break
			}
		}
		return best
	}
}

func randomFactor() float64 {
	if deterministic {
		return 1.0
	} else {
		return (1.0 - randomRange) + rand.Float64()*randomRange
	}
}
