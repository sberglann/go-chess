package main

type Color int

func main() {
	b := StartBoard()
	ms := KingMoves(*b)
	for _, m := range ms {
		PrettyMove(b.WhiteBB, m)
	}
}
