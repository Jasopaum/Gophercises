package main

import (
	"fmt"
	"gophercises/cardDeck"
	"sort"
)

func main() {
	d := cardDeck.NewDeck(cardDeck.WithMinCard(cardDeck.Seven),
		cardDeck.WithMaxCard(cardDeck.King),
		cardDeck.WithJokers(2),
		cardDeck.WithMultipleDecks(2),
		cardDeck.WithFilter(myFilter))
	for _, c := range d {
		fmt.Printf("%v\n", c)
	}
}

func mySort(d []cardDeck.Card) []cardDeck.Card {
	sort.Slice(d, myLess(d))
	return d
}
func myLess(d []cardDeck.Card) func(i, j int) bool {
	return func(i, j int) bool {
		return d[i].Rank < d[j].Rank
	}
}

func myFilter(d []cardDeck.Card) []cardDeck.Card {
	var res []cardDeck.Card
	for _, c := range d {
		if c.Rank != cardDeck.Nine {
			res = append(res, c)
		}
	}
	return res
}
