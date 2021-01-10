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

// Adding ace to foundation has much more value than a jack
var foundationScoreTable = [...]int{0, 13, 25, 36, 46, 55, 63, 70, 76, 81, 85, 88, 90, 91}

func calcScore(g GameMoment) int32 {
	foundationScore := 0
	maxFoundation := g.Foundation[0]
	minFoundation := g.Foundation[0]
	for _, x := range g.Foundation {
		foundationScore += foundationScoreTable[x]
		if x > maxFoundation {
			maxFoundation = x
		}
		if x < minFoundation {
			minFoundation = x
		}
	}

	if maxFoundation-minFoundation > 3 {
		foundationScore -= 10
	}
	if maxFoundation-minFoundation > 6 {
		foundationScore -= 20
	}
	if maxFoundation-minFoundation > 9 {
		foundationScore -= 30
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
		sequence := 0
		for row := g.CascadeLens[col] - 1; row >= 1; row-- {
			if !canMove(g.Cascades[col][row], g.Cascades[col][row-1]) {
				break
			}
			sequence++
		}

		cof := 1
		if sequence == int(g.CascadeLens[col]-1) {
			cof = 10
			if g.Cascades[col][0].Rank() == King {
				cof = 50
			} else if g.Cascades[col][0].Rank() == Queen {
				cof = 30
			} else if g.Cascades[col][0].Rank() == Jack {
				cof = 15
			}
		}
		sortedPairCount += cof * sequence
	}
	return int32(foundationScore*50 + freeCascades*500 + freecells*500 + sortedPairCount*5 - int(g.moves)*10)
}

func calcMinMovesNeeded(g *GameMoment) int32 {
	return int32(4*King - g.Foundation[0] - g.Foundation[1] - g.Foundation[2] - g.Foundation[3])
}

type Solver struct {
	Debug bool
	heap  *GameHeap
	cache map[GameMoment]bool

	shortestWin int32
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
		heap:        &GameHeap{},
		cache:       map[GameMoment]bool{},
		shortestWin: 999,
	}
	s.push(g)
	return s
}

func (s *Solver) Solve() GameMoment {
	for tries := 1; s.heap.Len() > 0; tries++ {
		g := heap.Pop(s.heap).(GameMoment)
		if g.moves+calcMinMovesNeeded(&g) >= s.shortestWin {
			continue
		}
		if g.isWon() {
			if s.Debug {
				log.Printf("WON!!! tried=%dK queue=%dK cache=%dM moves=%d score=%d", tries/1000, s.heap.Len()/1000, len(s.cache)/1000000, g.moves, g.score)
			}
			s.shortestWin = g.moves
			return g
		}
		if s.Debug && tries%100000 == 0 {
			log.Printf("tried=%dK queue=%dK cache=%dM moves=%d score=%d", tries/1000, s.heap.Len()/1000, len(s.cache)/1000000, g.moves, g.score)
		}
		s.addMoves(g)
	}
	panic("Not solvable! Probably it's a bug!!")
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
