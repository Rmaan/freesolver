package freesolver

import (
	"fmt"
	"strconv"
	"strings"
)

type Suit uint8

const SuitCount = 4

func (s Suit) String() string {
	switch s {
	case 0:
		return "♤"
	case 1:
		return "♡"
	case 2:
		return "♢"
	case 3:
		return "♧"
	}
	return fmt.Sprintf("<Invalid suit %d>", s)
}

func (s Suit) IsBlack() bool {
	switch s {
	case 0, 3:
		return true
	}
	return false
}

type Rank uint8

const RankCount = 13
const King Rank = 13
const Queen Rank = 12
const Jack Rank = 11

func (r Rank) String() string {
	switch r {
	case 0:
		return "-"
	case 1:
		return "A"
	case 2, 3, 4, 5, 6, 7, 8, 9, 10:
		return strconv.Itoa(int(r))
	case 11:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	}
	return fmt.Sprintf("<Invalid rank %d>", r)
}

type Card uint8
var EmptyCard Card

func NewCard(s Suit, r Rank) Card {
	return Card(uint8(s & 3) | uint8(r << 2))
}

func (c Card) Suit() Suit {
	return Suit(c & 3)
}

func (c Card) Rank() Rank {
	return Rank(c >> 2)
}

func (c Card) String() string {
	if c.IsEmpty() {
		return "   "
	}
	if c.Rank() == 10 {
		return c.Suit().String() + c.Rank().String()
	}
	return c.Suit().String() + c.Rank().String() + " "
}

func (c Card) IsEmpty() bool {
	return c == 0
}

const FreeCellCount = 4
const CascadesCount = 8
const MaxCascadeLen = 18

type GameMoment struct {
	FreeCells   [FreeCellCount]Card
	Cascades    [CascadesCount][MaxCascadeLen]Card
	CascadeLens [CascadesCount]uint8
	Foundation  [SuitCount]Rank

	moves  int32
	score  int32
	before *GameMoment
}

func (g GameMoment) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "Foundation: %s%s %s%s %s%s %s%s\n", Suit(0), g.Foundation[0], Suit(1), g.Foundation[1], Suit(2), g.Foundation[2], Suit(3), g.Foundation[3])
	fmt.Fprintf(b, "FreeCells: %s %s %s %s\n", g.FreeCells[0], g.FreeCells[1], g.FreeCells[2], g.FreeCells[3])

	dirty := true
	for row := 0; dirty; row++ {
		dirty = false
		for col := 0; col < CascadesCount; col++ {
			if row >= int(g.CascadeLens[col]) {
				fmt.Fprintf(b, "%s  ", EmptyCard)
			} else {
				fmt.Fprintf(b, "%s  ", g.Cascades[col][row])
				dirty = true
			}
		}
		fmt.Fprintf(b, "\n")
	}
	return b.String()
}

func (g GameMoment) FullGameString() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "============================================================================\n")
	for p := &g; p != nil; p = p.before {
		b.WriteString(p.String())
		fmt.Fprintf(b, "--------------------------------------------------------------------------\n")
	}
	fmt.Fprintf(b, "============================================================================\n")
	return b.String()
}

func (g GameMoment) isWon() bool {
	return g.Foundation[0] == King && g.Foundation[1] == King && g.Foundation[2] == King && g.Foundation[3] == King
}

func (g *GameMoment) cascadePut(card Card, col int) {
	g.Cascades[col][g.CascadeLens[col]] = card
	g.CascadeLens[col]++
}

func (g *GameMoment) cascadeRemove(col int) {
	g.Cascades[col][g.CascadeLens[col]-1] = EmptyCard
	g.CascadeLens[col]--
}

func (g *GameMoment) cascadeCard(col int) Card {
	return g.Cascades[col][g.CascadeLens[col]-1]
}
