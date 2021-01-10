package freesolver

import (
	"fmt"
	"math/rand"
	"strings"
)

func GenerateGame() GameMoment {
	var g GameMoment
	var allCards []Card
	for suit := Suit(0); suit < SuitCount; suit++ {
		for rank := Rank(1); rank <= RankCount; rank++ {
			allCards = append(allCards, NewCard(suit, rank))
		}
	}
	rand.Shuffle(len(allCards), func(i, j int) {
		allCards[i], allCards[j] = allCards[j], allCards[i]
	})

	for col := 0; col < CascadesCount; col++ {
		n := uint8(7)
		if col >= 4 {
			n--
		}

		g.CascadeLens[col] = n
		copy(g.Cascades[col][:], allCards[:n])
		allCards = allCards[n:]
	}
	return g
}

func GenerateGameFromString(all string) (*GameMoment, error) {
	all = strings.ToUpper(all)
	allCards := map[Card]bool{}
	var g GameMoment
	colIdx := 0
	for _, colString := range strings.Split(all, "\n") {
		colString = strings.TrimSpace(colString)
		if colString == "" {
			continue
		}
		var cascade []Card
		for _, cardString := range strings.Split(colString, " ") {
			card, err := CardFromString(cardString)
			if err != nil {
				return nil, fmt.Errorf("invalid card %s: %w", cardString, err)
			}
			if allCards[card] {
				return nil, fmt.Errorf("duplicate card %s", cardString)
			}
			allCards[card] = true
			cascade = append(cascade, card)
		}
		if colIdx >= CascadesCount {
			return nil, fmt.Errorf("too many cascades")
		}
		g.CascadeLens[colIdx] = uint8(copy(g.Cascades[colIdx][:], cascade))
		colIdx++
	}
	if len(allCards) != SuitCount*RankCount {
		return nil, fmt.Errorf("not all cards are passed, only %d", len(allCards))
	}
	return &g, nil
}

func CardFromString(s string) (Card, error) {
	if len(s) != 2 {
		return EmptyCard, fmt.Errorf("length should be 2")
	}

	var rank Rank
	switch s[0] {
	case 'A':
		rank = 1
	case 'J':
		rank = Jack
	case 'Q':
		rank = Queen
	case 'K':
		rank = King
	case 'T':
		rank = 10
	default:
		if s[0] >= '1' && s[0] <= '9' {
			rank = Rank(s[0] - '0')
		} else {
			return EmptyCard, fmt.Errorf("rank is wrong")
		}
	}

	var suit Suit
	switch s[1] {
	case 'S':
		suit = Spades
	case 'H':
		suit = Hearts
	case 'C':
		suit = Clubs
	case 'D':
		suit = Diamonds
	default:
		return EmptyCard, fmt.Errorf("suit is wrong")
	}
	return NewCard(suit, rank), nil
}
