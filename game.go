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

type Rank uint8

const RankCount = 13
const King Rank = 13

func (r Rank) String() string {
	switch r {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
		return strconv.Itoa(int(r))
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	}
	return fmt.Sprintf("<Invalid rank %d>", r)
}

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	zero := Card{}
	if c == zero {
		return " "
	}
	if c.Rank == 10 {
		return c.Suit.String() + c.Rank.String()
	}
	return c.Suit.String() + c.Rank.String() + " "
}

const FreeCellCount = 4
const CascadesCount = 8
const MaxCascadeLen = 15

type GameMoment struct {
	FreeCells [FreeCellCount]Card
	Cascades  [CascadesCount][MaxCascadeLen]Card
	CascadeLens [CascadesCount]uint8
	Foundation [SuitCount]Rank
}

func (g GameMoment) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b,"--------------------------------------------------------------------------\n")
	fmt.Fprintf(b,"Foundation: %s %s %s %s\n", g.Foundation[0], g.Foundation[1], g.Foundation[2], g.Foundation[3])
	fmt.Fprintf(b,"FreeCells: %s %s %s %s\n", g.FreeCells[0], g.FreeCells[1], g.FreeCells[2], g.FreeCells[3])

	dirty := true
	for row := 0; dirty; row++ {
		dirty = false
		for col := 0; col < CascadesCount; col++ {
			if row > int(g.CascadeLens[col]) {
				fmt.Fprintf(b, "   ")
			} else {
				fmt.Fprintf(b, "%s  ", g.Cascades[col][row])
				dirty = true
			}
		}
		fmt.Fprintf(b, "\n")
	}
	fmt.Fprintf(b,"--------------------------------------------------------------------------\n")
	return b.String()
}

func (g GameMoment) isWon() bool {
	return g.Foundation[0] == King && g.Foundation[1] == King && g.Foundation[2] == King && g.Foundation[3] == King
}
