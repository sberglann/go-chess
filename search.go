package gochess

import (
	"math"
	"math/rand"
	"sync"
)

type EvaluatedBoard struct {
	Board BitBoard
	Eval  float64
}

type SearchConfig struct {
	MaxDepth    int
	Randomness  float64
	Parallelism int
}

var maxDepth int
var randomness float64
var parallelism int

func EvaluatedMoves(board BitBoard, config SearchConfig) []EvaluatedBoard {
	maxDepth = config.MaxDepth
	randomness = config.Randomness
	parallelism = config.Parallelism

	legalMoves, numMoves := GenerateLegalStates(board)
	guard := make(chan struct{}, parallelism)
	var evals = make([]EvaluatedBoard, numMoves)
	var wg sync.WaitGroup
	wg.Add(numMoves)
	turn := board.Turn()

	if turn == White {
		for i := 0; i < numMoves; i++ {
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				eval := minimax(b, 1, false, -1000.0, 1000.0)
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, legalMoves[i])
		}
	} else {
		for i := 0; i < numMoves; i++ {
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				eval := minimax(b, 1, true, -1000.0, 1000.0)
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, legalMoves[i])
		}
		wg.Wait()
	}

	return evals
}

func SelectBest(turn Color, boards []EvaluatedBoard) EvaluatedBoard {
	var bestMove BitBoard
	var bestEval float64
	if turn == White {
		bestEval = -1000.0
		for _, evaledBoard := range boards {
			if evaledBoard.Eval > bestEval {
				bestEval = evaledBoard.Eval
				bestMove = evaledBoard.Board
			}
		}
	} else {
		bestEval = 1000.0
		for _, evaledBoard := range boards {
			if evaledBoard.Eval < bestEval {
				bestEval = evaledBoard.Eval
				bestMove = evaledBoard.Board
			}
		}
	}
	return EvaluatedBoard{bestMove, bestEval}
}

func minimax(board BitBoard, depth int, isWhite bool, alpha float64, beta float64) float64 {
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
			currentEval := minimax(child, depth+1, false, alpha, beta)
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
			currentEval := minimax(child, depth+1, true, alpha, beta)
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
	if randomness == 0.0 {
		return 1.0
	} else {
		return (1.0 - randomness) + rand.Float64()*randomness
	}
}
