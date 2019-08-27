package blackjack

import (
	"fmt"
	"gophercises/cardDeck"
)

type Player interface {
	Bet(bool) string
	Play(Hand, cardDeck.Card) action
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

func (ai humanPlayer) Play(player Hand, dealer cardDeck.Card) action {
	var input string
	for {
		fmt.Printf("Dealer's hand: %s, ***Hidden***\n", dealer)
		fmt.Println("Player's hand: ", player)
		if len(player) == 2 {
			fmt.Println("Do you (h)it, (s)tand, or (d)ouble?")
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
		}
	}
}
