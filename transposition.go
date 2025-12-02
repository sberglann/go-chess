package main

import "sync"

type TranspositionTable struct {
	lower sync.Map
	upper sync.Map
	exact sync.Map
}

func newTranspositionTable() *TranspositionTable {
	return &TranspositionTable{
		lower: sync.Map{},
		upper: sync.Map{},
		exact: sync.Map{},
	}

}

func (tt *TranspositionTable) Clear() {
	tt.lower.Clear()
	tt.upper.Clear()
	tt.exact.Clear()
}

type TTEntry struct {
	score float64
	depth int8
}

type TTResult struct {
	lower float64
	upper float64
}

func (tt *TranspositionTable) getExact(board *BitBoard, depth int8) (bool, float64) {
	entry, found := tt.exact.Load(board.Hash())
	if !found {
		return false, 0
	}
	ttEntry := entry.(TTEntry)
	if ttEntry.depth >= depth {
		return true, ttEntry.score
	}
	return false, 0
}

func (tt *TranspositionTable) getUpperAndLower(board *BitBoard, depth int8) TTResult {
	hash := board.Hash()
	
	lower := -1000.0
	lowerEntry, found := tt.lower.Load(hash)
	if found {
		ttEntry := lowerEntry.(TTEntry)
		if ttEntry.depth >= depth {
			lower = ttEntry.score
		}
	}
	
	upper := 1000.0
	upperEntry, found := tt.upper.Load(hash)
	if found {
		ttEntry := upperEntry.(TTEntry)
		if ttEntry.depth >= depth {
			upper = ttEntry.score
		}
	}

	return TTResult{lower: lower, upper: upper}
}

func (tt *TranspositionTable) setExact(board *BitBoard, score float64, depth int8) {
	hash := board.Hash()
	existingEntry, existingFound := tt.exact.Load(hash)
	if existingFound {
		ttEntry := existingEntry.(TTEntry)
		if ttEntry.depth >= depth {
			return
		}
	}
	tt.exact.Store(hash, TTEntry{score: score, depth: depth})
}
func (tt *TranspositionTable) setLower(board *BitBoard, score float64, depth int8) {
	tt.lower.Store(board.Hash(), TTEntry{score: score, depth: depth})
}

func (tt *TranspositionTable) setUpper(board *BitBoard, score float64, depth int8) {
	tt.upper.Store(board.Hash(), TTEntry{score: score, depth: depth})
}
