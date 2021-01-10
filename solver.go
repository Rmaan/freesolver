package freesolver

import (
	"container/heap"
	"log"
	"os"
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

func calcScore(g GameMoment) int {
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
	return done*1000 + freeCascades*200 + freecells*100 + sortedPairCount*5 - g.depth
}

type Solver struct {
	heap  *GameHeap
	cache map[GameMoment]bool
}

func (s *Solver) push(g GameMoment) {
	if hasRepeats(&g) {
		return
	}
	gBefore := g.before
	gDepth := g.depth
	gScore := g.score
	g.before = nil
	g.depth = 0
	g.score = 0
	if s.cache[g] {
		return
	}
	s.cache[g] = true
	g.before = gBefore
	g.depth = gDepth
	g.score = gScore

	g.score = calcScore(g)
	heap.Push(s.heap, g)
}

func NewSolver(g GameMoment) *Solver {
	return &Solver{
		heap:  &GameHeap{g},
		cache: map[GameMoment]bool{},
	}
}

func (s *Solver) Solve() {
	for callCount := 1; s.heap.Len() > 0; callCount++ {
		g := heap.Pop(s.heap).(GameMoment)
		if g.isWon() {
			log.Printf("called %dK queue %dK depth %d score %d game:\n%s", callCount/1000, s.heap.Len()/1000, g.depth, g.score, g.FullGameString())
			log.Printf("Won!!!")
			os.Exit(0)
		}
		if callCount%100000 == 0 {
			log.Printf("called %dK queue %dK depth %d score %d game:\n%s", callCount/1000, s.heap.Len()/1000, g.depth, g.score, g.String())
		}
		s.addMoves(g)
		//if callCount > 20 {
		//	return
		//}
	}
}

func canMove(from, to Card) bool {
	if to.IsEmpty() {
		return true
	}
	if from.Rank != to.Rank-1 {
		return false
	}
	return from.Suit.IsBlack() != to.Suit.IsBlack()
}

func hasRepeats(g *GameMoment) bool {
	gBefore := g.before
	gDepth := g.depth
	gScore := g.score
	g.before = nil
	g.depth = 0
	g.score = 0

	for old, count := gBefore, 0; old != nil && count < 5; old = old.before {
		oldBefore := old.before
		oldDepth := old.depth
		oldScore := old.score
		old.before = nil
		old.depth = 0
		old.score = 0
		isEqual := *g == *old
		old.before = oldBefore
		old.depth = oldDepth
		old.score = oldScore
		if isEqual {
			g.before = gBefore
			g.depth = gDepth
			g.score = gScore
			return true
		}
		count++
	}

	g.before = gBefore
	g.depth = gDepth
	g.score = gScore
	return false
}

func (s *Solver) addMoves(g GameMoment) {
	gp := &g

	// Cascade to foundation
	for col := range g.Cascades {
		if g.CascadeLens[col] == 0 {
			continue
		}
		card := g.cascadeCard(col)
		if g.Foundation[card.Suit]+1 == card.Rank {
			g := g
			g.depth++
			g.before = gp
			g.Foundation[card.Suit]++
			g.cascadeRemove(col)
			s.push(g)
		}
	}

	// FreeCell to foundation
	for idx, card := range g.FreeCells {
		if g.Foundation[card.Suit]+1 == card.Rank {
			g := g
			g.depth++
			g.before = gp
			g.Foundation[card.Suit]++
			g.FreeCells[idx] = Card{}
			s.push(g)
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
				g.depth++
				g.before = gp
				g.cascadePut(card, col)
				g.FreeCells[idx] = Card{}
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
				g.depth++
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
		g.depth++
		g.before = gp
		g.FreeCells[emptyFreecellIdx] = g.cascadeCard(colFrom)
		g.cascadeRemove(colFrom)
		s.push(g)
	}
}
