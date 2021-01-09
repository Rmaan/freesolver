package freesolver

import (
	"log"
	"os"
)

func Solve(g GameMoment) {
	solve(g, 0)
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
	g.before = nil

	for old, count := gBefore, 0; old != nil && count < 10; old = old.before {
		oldBefore := old.before
		old.before = nil
		isEqual := *g == *old
		old.before = oldBefore
		if isEqual {
			g.before = gBefore
			return true
		}
		count++
	}

	g.before = gBefore
	return false
}

var maxDepth = 20
var callCount = 0
var maxDepthReachedCount = 0

func solve(g GameMoment, depth int) {
	callCount++
	if callCount%10000000 == 0 {
		log.Printf("called %dM depth %d fullgame:\n%s", callCount/1000000, depth, " ")
	}

	if hasRepeats(&g) {
		return
	}

	if depth > maxDepth {
		maxDepthReachedCount++
		if maxDepthReachedCount % 1000000 == 0 {
			log.Printf("max depth reached %dM times", maxDepthReachedCount/1000000)
		}
		return
	}
	if g.isWon() {
		log.Printf("Won!!!")
		os.Exit(0)
	}
	gp := &g

	// Cascade to foundation
	for col := range g.Cascades {
		if g.CascadeLens[col] == 0 {
			continue
		}
		card := g.cascadeCard(col)
		if g.Foundation[card.Suit]+1 == card.Rank {
			g := g
			g.before = gp
			g.Foundation[card.Suit]++
			g.cascadeRemove(col)
			solve(g, depth + 1)
		}
	}

	// FreeCell to foundation
	for idx, card := range g.FreeCells {
		if g.Foundation[card.Suit]+1 == card.Rank {
			g := g
			g.before = gp
			g.Foundation[card.Suit]++
			g.FreeCells[idx] = Card{}
			solve(g, depth + 1)
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
				g.before = gp
				g.cascadePut(card, col)
				g.FreeCells[idx] = Card{}
				solve(g, depth + 1)
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
				g.before = gp
				g.cascadePut(g.cascadeCard(colFrom), colTo)
				g.cascadeRemove(colFrom)
				solve(g, depth + 1)
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
		g.before = gp
		g.FreeCells[emptyFreecellIdx] = g.cascadeCard(colFrom)
		g.cascadeRemove(colFrom)
		solve(g, depth + 1)
	}
}
