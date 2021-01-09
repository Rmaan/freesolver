package freesolver

import "math/rand"

func GenerateGame() GameMoment {
	var g GameMoment
	var allCards []Card
	for suit := Suit(0); suit < SuitCount; suit++ {
		for rank := Rank(1); rank <= RankCount; rank++ {
			allCards = append(allCards, Card{Suit: suit, Rank: rank})
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
