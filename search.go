package main

import (
	"math"
	"sync"
)

type EvaluatedBoard struct {
	board BitBoard
	eval  float64
}

const maxDepth = 5

func BestMove(board BitBoard) BitBoard {
	var bestMove BitBoard

	legalMoves := GenerateLegalMoves(board)
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
				evals[j] = EvaluatedBoard{b, eval}
				<-guard
			}(i, b)
		}
		wg.Wait()
		maxEval := math.Inf(-1)
		for _, evaledBoard := range evals {
			if evaledBoard.eval > maxEval {
				maxEval = evaledBoard.eval
				bestMove = evaledBoard.board
			}
		}

	} else {
		for i, b := range legalMoves {
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				eval := minimax(b, 1, true, math.Inf(-1), math.Inf(1))
				evals[j] = EvaluatedBoard{b, eval}
				<-guard
			}(i, b)
		}
		wg.Wait()
		minEval := math.Inf(1)
		for _, evaledBoard := range evals {
			if evaledBoard.eval < minEval {
				minEval = evaledBoard.eval
				bestMove = evaledBoard.board
			}
		}
	}

	return bestMove
}

func minimax(board BitBoard, depth int, isWhite bool, alpha float64, beta float64) float64 {
	children := GenerateLegalMoves(board)
	var sign int
	if isWhite {
		sign = 1
	} else {
		sign = -1
	}
	if len(children) <= 0 {
		// Check mate
		return math.Inf(-sign)
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
