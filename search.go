package main

import (
	cache "github.com/sberglann/uint64gocache"
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
const useCache = false

var upperMemo = cache.New()
var lowerMemo = cache.New()

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
				eval := minimax(b, 1, false, -1000.0, 1000.0)
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, b)
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
		for i, b := range legalMoves {
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				eval := minimax(b, 1, true, -1000.0, 1000.0)
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, b)
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

func minimax(board BitBoard, depth int, isWhite bool, alpha float64, beta float64) float64 {
	key := board.Hash(depth)
	if useCache {
		lowerBound, lowerCached := lowerMemo.Get(key)
		if lowerCached {
			if lowerBound.(float64) >= beta {
				return lowerBound.(float64)
			}
			alpha = math.Max(alpha, lowerBound.(float64))
		}
		upperBound, upperCached := upperMemo.Get(key)
		if upperCached {
			if upperBound.(float64) <= alpha {
				return upperBound.(float64)
			}
			beta = math.Min(beta, upperBound.(float64))
		}
	}
	if depth == maxDepth {
		return Eval(board)
	}
	children := GenerateLegalStates(board)
	var bestEval float64
	if len(children) == 0 {
		return Eval(board)
	} else if isWhite {
		bestEval = -1000.0
		for _, child := range children {
			currentEval := minimax(child, depth+1, false, alpha, beta)
			bestEval = math.Max(bestEval, currentEval)
			alpha = math.Max(alpha, bestEval)
			if beta <= alpha {
				break
			}
		}
	} else {
		bestEval = 1000.0
		for _, child := range children {
			currentEval := minimax(child, depth+1, true, alpha, beta)
			bestEval = math.Min(bestEval, currentEval)
			beta = math.Min(beta, bestEval)
			if beta <= alpha {
				break
			}
		}
	}
	if bestEval <= alpha {
		upperMemo.Set(key, bestEval)
	}
	if bestEval > alpha && bestEval < beta {
		upperMemo.Set(key, bestEval)
		lowerMemo.Set(key, bestEval)
	}
	if bestEval >= beta {
		lowerMemo.Set(key, bestEval)
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
