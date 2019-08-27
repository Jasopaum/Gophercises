package main

import (
	"fmt"
	"gophercises/blackjack"
)

// Implement the blackjack AI interface
type myAI struct{}

func main() {
	ai := blackjack.HumanPlayer()

	opts := blackjack.Options{
		Hands: 1,
		Decks: 3,
	}
	game := blackjack.New(opts)
	winnings := game.Play(ai)
	fmt.Println("Our AI won/lost:", winnings)
}
