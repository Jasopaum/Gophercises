package main

import (
	"fmt"
	"gophercises/cardDeck"
	"strings"
)

type hand []cardDeck.Card

type gameState struct {
	DealerHand hand
	PlayerHand hand
	Deck       []cardDeck.Card
}

func main() {
	gs := initGame()

	var nbRounds int
	fmt.Println("How many rounds do you want to play?")
	fmt.Scanf("%d\n", &nbRounds)

	for r := 0; r < nbRounds; r++ {
		distribute(&gs)
		playerPlays(&gs)
		dealerPlays(&gs)
		showResult(&gs)
		resetHand(&gs)
	}
}

func draw(h hand, d []cardDeck.Card) (hand, []cardDeck.Card) {
	c, d := d[0], d[1:]
	h = append(h, c)
	return h, d
}

func distribute(gs *gameState) {
	for i := 0; i < 2; i++ {
		gs.PlayerHand, gs.Deck = draw(gs.PlayerHand, gs.Deck)
		gs.DealerHand, gs.Deck = draw(gs.DealerHand, gs.Deck)
	}
}

func initGame() gameState {
	deck := cardDeck.NewDeck(cardDeck.WithMultipleDecks(3), cardDeck.WithShuffle())

	gs := gameState{
		Deck: deck,
	}

	return gs
}

func playerPlays(gs *gameState) {
	var input string
	for input != "s" {
		fmt.Println("Dealer's hand: ", gs.DealerHand.StringDealer())
		fmt.Println("Player's hand: ", gs.PlayerHand)
		fmt.Println("Do you (h)it or (s)tand?")
		fmt.Scanf("%s\n", &input)
		if input == "h" {
			gs.PlayerHand, gs.Deck = draw(gs.PlayerHand, gs.Deck)
		}
	}
}

func dealerPlays(gs *gameState) {
	s, ace := score(gs.DealerHand)
	for s < 17 || (s == 17 && ace) {
		gs.DealerHand, gs.Deck = draw(gs.DealerHand, gs.Deck)
		s, ace = score(gs.DealerHand)
	}
}

func showResult(gs *gameState) {
	sd, _ := score(gs.DealerHand)
	sp, _ := score(gs.PlayerHand)

	fmt.Println("Dealer:")
	fmt.Println("\t", gs.DealerHand)
	fmt.Println("\tscore:", sd)

	fmt.Println("Player:")
	fmt.Println("\t", gs.PlayerHand)
	fmt.Println("\tscore:", sp)

	if sp > 21 {
		fmt.Println("You busted")
	} else if sd > 21 {
		fmt.Println("Dealer busted")
	} else if sp > sd {
		fmt.Println("You won")
	} else if sp < sd {
		fmt.Println("Dealer won")
	} else if sp == sd {
		fmt.Println("Tie")
	}
}

func resetHand(gs *gameState) {
	gs.DealerHand = nil
	gs.PlayerHand = nil
}

func score(h hand) (int, bool) {
	s := 0
	acePresent := false
	for _, c := range h {
		s += min(int(c.Rank), 10)
		if c.Rank == cardDeck.Ace {
			acePresent = true
		}
	}
	if acePresent && s < 12 {
		s += 10
	}
	return s, acePresent
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func (h hand) String() string {
	str := make([]string, len(h))
	for i, c := range h {
		str[i] = c.String()
	}
	return strings.Join(str, ", ")
}

func (h hand) StringDealer() string {
	return fmt.Sprintf("%s, ***Hidden***", h[0])
}
