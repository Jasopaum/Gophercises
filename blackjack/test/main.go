package main

import (
	"fmt"
	"gophercises/blackjack"
	"strconv"
)

// Implement the blackjack AI interface
type myAI struct{}

func main() {
	ai := blackjack.HumanPlayer()

	var (
		h       string
		nbHands int
		err     error
	)
	for {
		fmt.Println("How many hands do you want to play?")
		fmt.Scanf("%s\n", &h)
		if nbHands, err = strconv.Atoi(h); err == nil {
			break
		}
	}

	opts := blackjack.Options{
		Hands: nbHands,
		Decks: 3,
	}
	game := blackjack.New(opts)
	winnings := game.Play(ai)
	fmt.Println("Our AI won/lost:", winnings)
}
