package freesolver

import (
	"container/heap"
	"log"
)

type GameHeap []GameMoment

func (h GameHeap) Len() int {
	return len(h)
}

func (h GameHeap) Less(i, j int) bool {
	return h[i].score > h[j].score
}

func (h GameHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *GameHeap) Push(x interface{}) {
	g := x.(GameMoment)
	*h = append(*h, g)
}

func (h *GameHeap) Pop() interface{} {
	n := len(*h)
	item := (*h)[n-1]
	*h = (*h)[:n-1]
	return item
}

func calcScore(g GameMoment) int32 {
	done := 0
	for _, x := range g.Foundation {
		done += int(x)
	}
	freecells := 0
	for _, x := range g.FreeCells {
		if x.IsEmpty() {
			freecells++
		}
	}
	freeCascades := 0
	sortedPairCount := 0
	for col := range g.Cascades {
		if g.CascadeLens[col] == 0 {
			freeCascades++
			continue
		}
		for row := uint8(1); row < g.CascadeLens[col]; row++ {
			if canMove(g.Cascades[col][row], g.Cascades[col][row-1]) {
				sortedPairCount++
			}
		}
	}
	return int32(done*500 + freeCascades*500 + freecells*500 + sortedPairCount*5 - int(g.moves))
}

type Solver struct {
	heap  *GameHeap
	cache map[GameMoment]bool
}

func sortFreeCells(g *GameMoment) {
	swapped := true
	for swapped {
		swapped = false
		for i := 1; i < len(g.FreeCells); i++ {
			if g.FreeCells[i-1] > g.FreeCells[i] {
				g.FreeCells[i-1], g.FreeCells[i] = g.FreeCells[i], g.FreeCells[i-1]
				swapped = true
			}
		}
	}
}

func sortCascades(g *GameMoment) {
	swapped := true
	for swapped {
		swapped = false
		for i := 1; i < len(g.Cascades); i++ {
			if g.Cascades[i-1][0] > g.Cascades[i][0] {
				g.Cascades[i-1], g.Cascades[i] = g.Cascades[i], g.Cascades[i-1]
				g.CascadeLens[i-1], g.CascadeLens[i] = g.CascadeLens[i], g.CascadeLens[i-1]
				swapped = true
			}
		}
	}
}

func (s *Solver) push(g GameMoment) {
	sortFreeCells(&g)
	sortCascades(&g)
	gBefore := g.before
	gDepth := g.moves
	g.before = nil
	g.moves = 0
	g.score = 0
	g.score = 0
	if s.cache[g] {
		return
	}
	s.cache[g] = true
	g.before = gBefore
	g.moves = gDepth
	g.score = calcScore(g)

	heap.Push(s.heap, g)
}

func NewSolver(g GameMoment) *Solver {
	s := &Solver{
		heap:  &GameHeap{},
		cache: map[GameMoment]bool{},
	}
	s.push(g)
	return s
}

func (s *Solver) Solve() GameMoment {
	for callCount := 1; s.heap.Len() > 0; callCount++ {
		g := heap.Pop(s.heap).(GameMoment)
		if g.isWon() {
			log.Printf("Won! called %dK queue %dK cache %dM moves %d score %d", callCount/1000, s.heap.Len()/1000, len(s.cache)/1000000, g.moves, g.score)
			return g
		}
		if callCount%100000 == 0 {
			log.Printf("called %dK queue %dK cache %dM moves %d score %d game:\n%s", callCount/1000, s.heap.Len()/1000, len(s.cache)/1000000, g.moves, g.score, g.String())
		}
		s.addMoves(g)
	}
	panic("Not solvable! It's usually a bug!!")
}

func canMove(from, to Card) bool {
	if to.IsEmpty() {
		return true
	}
	if from.Rank() != to.Rank()-1 {
		return false
	}
	return from.Suit().IsBlack() != to.Suit().IsBlack()
}

func (s *Solver) addMoves(g GameMoment) {
	gp := &g

	// Cascade to foundation
	for col := range g.Cascades {
		if g.CascadeLens[col] == 0 {
			continue
		}
		card := g.cascadeCard(col)
		if g.Foundation[card.Suit()]+1 == card.Rank() {
			g := g
			g.moves++
			g.before = gp
			g.Foundation[card.Suit()]++
			g.cascadeRemove(col)
			s.push(g)
			return
		}
	}

	// FreeCell to foundation
	for idx, card := range g.FreeCells {
		if g.Foundation[card.Suit()]+1 == card.Rank() {
			g := g
			g.moves++
			g.before = gp
			g.Foundation[card.Suit()]++
			g.FreeCells[idx] = EmptyCard
			s.push(g)
			return
		}
	}

	// FreeCell to cascades
	for idx, card := range g.FreeCells {
		if card.IsEmpty() {
			continue
		}
		for col := range g.Cascades {
			if g.CascadeLens[col] == 0 || canMove(card, g.cascadeCard(col)) {
				g := g
				g.moves++
				g.before = gp
				g.cascadePut(card, col)
				g.FreeCells[idx] = EmptyCard
				s.push(g)
			}
		}
	}

	// Play in cascade
	for colFrom := range g.Cascades {
		if g.CascadeLens[colFrom] == 0 {
			continue
		}
		for colTo := range g.Cascades {
			if colFrom == colTo {
				continue
			}
			if g.CascadeLens[colTo] == 0 || canMove(g.cascadeCard(colFrom), g.cascadeCard(colTo)) {
				g := g
				g.moves++
				g.before = gp
				g.cascadePut(g.cascadeCard(colFrom), colTo)
				g.cascadeRemove(colFrom)
				s.push(g)
			}
		}
	}

	// Cascade to freecell
	emptyFreecellIdx := -1
	for idx, card := range g.FreeCells {
		if card.IsEmpty() {
			emptyFreecellIdx = idx
			break
		}
	}
	if emptyFreecellIdx == -1 {
		return
	}
	for colFrom := range g.Cascades {
		if g.CascadeLens[colFrom] == 0 {
			continue
		}
		g := g
		g.moves++
		g.before = gp
		g.FreeCells[emptyFreecellIdx] = g.cascadeCard(colFrom)
		g.cascadeRemove(colFrom)
		s.push(g)
	}
}
