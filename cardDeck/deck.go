package cardDeck

import (
	"fmt"
	"math/rand"
	"time"
)

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

var suits = [4]Suit{Spade, Diamond, Club, Heart}

type Rank uint8

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Card struct composed of a Rank that should be Ace, Two, ..., King
// and a Suit (Spade, Diamond, Club, Heart, or Joker)
type Card struct {
	Rank
	Suit
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

func NewDeck(opts ...funcOptions) []Card {
	options := deckOptions{
		minCard:     1,
		maxCard:     13,
		nbDecks:     1,
		nbJokers:    0,
		shuffle:     false,
		sortingFn:   func(d []Card) []Card { return d },
		filteringFn: func(d []Card) []Card { return d },
	}

	for _, f := range opts {
		f(&options)
	}

	nbCards := options.nbDecks*4*(int(options.maxCard)-int(options.minCard)+1) + options.nbJokers
	res := make([]Card, nbCards)
	i := 0
	for n := 0; n < options.nbDecks; n++ {
		for _, s := range suits {
			for v := options.minCard; v <= options.maxCard; v++ {
				res[i] = Card{v, s}
				i++
			}
		}
	}
	for j := 0; j < options.nbJokers; j++ {
		res[i+j] = Card{0, Joker}
	}

	res = options.sortingFn(res)
	res = options.filteringFn(res)
	if options.shuffle {
		res = shuffleDeck(res)
	}

	return res
}

type deckOptions struct {
	minCard     Rank
	maxCard     Rank
	nbDecks     int
	nbJokers    int
	shuffle     bool
	sortingFn   func([]Card) []Card
	filteringFn func([]Card) []Card
}

type funcOptions func(*deckOptions)

// To set the max rank of the cards in the deck
func WithMaxCard(max Rank) func(*deckOptions) {
	return func(opts *deckOptions) {
		opts.maxCard = max
	}
}

// To set the min rank of the cards in the deck
func WithMinCard(min Rank) func(*deckOptions) {
	return func(opts *deckOptions) {
		opts.minCard = min
	}
}

// If several sets of cards are needed in the deck
// n is the number of sets needed
func WithMultipleDecks(n int) func(*deckOptions) {
	return func(opts *deckOptions) {
		opts.nbDecks = n
	}
}

// To add jokers in the deck. Note that n jokers will be added
// and not n times the number of sets of cards
func WithJokers(n int) func(*deckOptions) {
	return func(opts *deckOptions) {
		opts.nbJokers = n
	}
}

// Add this functional option to perform shuffling
func WithShuffle() func(*deckOptions) {
	return func(opts *deckOptions) {
		opts.shuffle = true
	}
}
func shuffleDeck(d []Card) []Card {
	rand.Seed(time.Now().UTC().Unix())
	res := make([]Card, len(d))
	perm := rand.Perm(len(d))
	for i, j := range perm {
		res[i] = d[j]
	}
	return res
}

// Add this functional option to sort the deck with a custom function
func WithSort(f func([]Card) []Card) func(*deckOptions) {
	return func(opts *deckOptions) {
		opts.sortingFn = f
	}
}

// Add this functional option to filter the deck with a custom function
func WithFilter(f func([]Card) []Card) func(*deckOptions) {
	return func(opts *deckOptions) {
		opts.filteringFn = f
	}
}
