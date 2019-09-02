package blackjack

import (
	"fmt"
	"gophercises/cardDeck"
)

type Player interface {
	Bet(bool) string
	Play([]cardDeck.Card, cardDeck.Card) Action
	EndRound([]hand, []cardDeck.Card)
}

func HumanPlayer() humanPlayer {
	return humanPlayer{}
}

type humanPlayer struct{}

func (ai humanPlayer) Bet(shuffled bool) string {
	var bet string
	fmt.Println("How much do you want to bet?")
	fmt.Scanf("%s\n", &bet)
	return bet
}

func (ai humanPlayer) Play(player []cardDeck.Card, dealer cardDeck.Card) Action {
	var input string
	for {
		fmt.Printf("Dealer's hand: %s, ***Hidden***\n", dealer)
		fmt.Println("Player's hand: ", player)
		if len(player) == 2 {
			fmt.Println("Do you (h)it, (s)tand, s(p)lit, or (d)ouble?")
		} else {
			fmt.Println("Do you (h)it or (s)tand?")
		}
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return Hit
		case "s":
			return Stand
		case "d":
			return Double
		case "p":
			return Split
		}
	}
}

func (ai humanPlayer) EndRound(player []hand, dealer []cardDeck.Card) {
	sd, _ := Score(dealer)
	bjd := Blackjack(dealer)

	fmt.Println("Dealer:")
	fmt.Println("\t", dealer)
	fmt.Println("\tscore:", sd)

	for _, h := range player {
		cards := h.cards
		sp, _ := Score(cards)
		bjp := Blackjack(cards)

		fmt.Println("Player:")
		fmt.Println("\t", h)
		fmt.Println("\twith bet:", h.bet)
		fmt.Println("\tscore:", sp)

		switch {
		case bjp && bjd:
			fmt.Println("Tie")
		case bjd:
			fmt.Println("Dealer won")
		case bjp:
			fmt.Println("You won")
		case sp > 21:
			fmt.Println("You busted")
		case sd > 21:
			fmt.Println("Dealer busted")
		case sp > sd:
			fmt.Println("You won")
		case sp < sd:
			fmt.Println("Dealer won")
		case sp == sd:
			fmt.Println("Tie")
		}
	}
}
