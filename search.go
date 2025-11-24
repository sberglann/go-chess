package main

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"

	cache "github.com/sberglann/uint64gocache"
)

type EvaluatedBoard struct {
	board BitBoard
	eval  float64
}

const defaultMaxDepth = 5
const deterministic = true
const randomRange = 0
const useCache = false

var upperMemo = cache.New()
var lowerMemo = cache.New()

// BestMove finds the best move using iterative deepening with optional time constraint.
// If timeLimit is 0 or negative, it uses defaultMaxDepth. Otherwise, it performs
// iterative deepening, increasing depth until time runs out, and returns the best
// move found so far.
func BestMove(board BitBoard, timeLimit time.Duration) EvaluatedBoard {
	legalMoves, numMoves := GenerateLegalStates(board)
	if numMoves == 0 {
		return EvaluatedBoard{board, Eval(board)}
	}

	// Convert array to slice
	legalMovesSlice := legalMoves[:numMoves]

	// If no time limit, use fixed depth
	if timeLimit <= 0 {
		return bestMoveAtDepth(board, legalMovesSlice, numMoves, defaultMaxDepth, nil)
	}

	// Create context that will be cancelled when time runs out
	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	// Shared state for the best move found so far
	var bestMoveMutex sync.Mutex
	var currentBestMove BitBoard
	var currentBestEval float64
	hasBestMove := false

	// Start iterative deepening in a goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		maxDepth := 1
		for {
			// Try to find best move at current depth
			// Pass context so it can be cancelled
			result := bestMoveAtDepth(board, legalMovesSlice, numMoves, maxDepth, ctx)
			
			// Only update best move if we got a valid result (not cancelled mid-search)
			// Check if the result is valid by ensuring we have a proper board
			if result.board.WhiteBB > 0 || result.board.BlackBB > 0 {
				bestMoveMutex.Lock()
				currentBestMove = result.board
				currentBestEval = result.eval
				hasBestMove = true
				bestMoveMutex.Unlock()
			}

			// Check if cancelled after completing this depth
			select {
			case <-ctx.Done():
				return
			default:
			}

			maxDepth++
		}
	}()

	// Wait for context to be cancelled (timeout)
	<-ctx.Done()
	
	// Give the goroutine a moment to finish updating the best move
	// This ensures we get the result from the last completed depth
	select {
	case <-done:
		// Search completed naturally
	case <-time.After(50 * time.Millisecond):
		// Give it a bit more time to finish the current depth
	}

	// Return the best move found, or fallback to first legal move
	bestMoveMutex.Lock()
	defer bestMoveMutex.Unlock()
	if hasBestMove {
		return EvaluatedBoard{currentBestMove, currentBestEval}
	}
	
	// Fallback: evaluate all moves at depth 1 to find the best one
	// This ensures we always return a proper evaluation, not just the first move
	if numMoves > 0 {
		// Quick evaluation of all moves to find the best
		bestFallback := legalMoves[0]
		bestFallbackEval := Eval(legalMoves[0])
		for i := 1; i < numMoves; i++ {
			eval := Eval(legalMoves[i])
			if board.Turn() == White {
				if eval > bestFallbackEval {
					bestFallback = legalMoves[i]
					bestFallbackEval = eval
				}
			} else {
				if eval < bestFallbackEval {
					bestFallback = legalMoves[i]
					bestFallbackEval = eval
				}
			}
		}
		return EvaluatedBoard{bestFallback, bestFallbackEval}
	}
	
	return EvaluatedBoard{board, Eval(board)}
}

// bestMoveAtDepth finds the best move at a specific depth.
// Note: minimax will recursively generate moves at each level of the search tree,
// so we don't need to generate moves here - we just evaluate each root-level move
// by calling minimax, which will traverse the tree to the specified depth.
// ctx can be used to cancel the search early.
func bestMoveAtDepth(board BitBoard, legalMoves []BitBoard, numMoves int, depth int, ctx context.Context) EvaluatedBoard {
	var bestMove BitBoard
	var bestEval float64

	maxRoutines := 16
	guard := make(chan struct{}, maxRoutines)
	var evals = make([]EvaluatedBoard, numMoves)
	var wg sync.WaitGroup
	wg.Add(numMoves)

	if board.Turn() == White {
		for i := 0; i < numMoves; i++ {
			// Check for cancellation before starting new goroutine
			if ctx != nil {
				select {
				case <-ctx.Done():
					// If cancelled, return empty result - caller should not use incomplete results
					return EvaluatedBoard{}
				default:
				}
			}
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				// minimax will recursively generate and explore moves at each level
				// until it reaches the specified depth
				eval := minimax(b, 1, false, -1000.0, 1000.0, depth, ctx)
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, legalMoves[i])
		}
		wg.Wait()
		
		// Check if we were cancelled during the search
		if ctx != nil {
			select {
			case <-ctx.Done():
				// Return empty result if cancelled - incomplete search
				return EvaluatedBoard{}
			default:
			}
		}
		
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
			// Check for cancellation before starting new goroutine
			if ctx != nil {
				select {
				case <-ctx.Done():
					// If cancelled, return empty result - caller should not use incomplete results
					return EvaluatedBoard{}
				default:
				}
			}
			guard <- struct{}{}
			go func(j int, b BitBoard) {
				defer wg.Done()
				eval := minimax(b, 1, true, -1000.0, 1000.0, depth, ctx)
				evals[j] = EvaluatedBoard{b, eval * randomFactor()}
				<-guard
			}(i, legalMoves[i])
		}
		wg.Wait()
		
		// Check if we were cancelled during the search
		if ctx != nil {
			select {
			case <-ctx.Done():
				// Return empty result if cancelled - incomplete search
				return EvaluatedBoard{}
			default:
			}
		}
		
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

func minimax(board BitBoard, depth int, isWhite bool, alpha float64, beta float64, maxDepth int, ctx context.Context) float64 {
	// Check for cancellation periodically
	if ctx != nil {
		select {
		case <-ctx.Done():
			// Return current evaluation if cancelled
			return Eval(board)
		default:
		}
	}

	if depth == maxDepth {
		return Eval(board)
	}
	// Generate legal moves at this level of the search tree
	// This is called recursively, so moves are generated at each level
	children, numChildren := GenerateLegalStates(board)
	var bestEval float64
	if numChildren == 0 {
		return Eval(board)
	} else if isWhite {
		bestEval = -1000.0
		i := 0
		child := children[i]
		for child.WhiteBB > 0 {
			// Check for cancellation before each recursive call
			if ctx != nil {
				select {
				case <-ctx.Done():
					return bestEval
				default:
				}
			}
			currentEval := minimax(child, depth+1, false, alpha, beta, maxDepth, ctx)
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
			// Check for cancellation before each recursive call
			if ctx != nil {
				select {
				case <-ctx.Done():
					return bestEval
				default:
				}
			}
			currentEval := minimax(child, depth+1, true, alpha, beta, maxDepth, ctx)
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
