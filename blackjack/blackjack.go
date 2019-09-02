package blackjack

import (
	"errors"
	"fmt"
	"gophercises/cardDeck"
	"strconv"
	"strings"
)

type hand struct {
	cards []cardDeck.Card
	bet   int
}

type gameState struct {
	nbDecks  int
	deck     []cardDeck.Card
	nbRounds int

	dealerHand  []cardDeck.Card
	playerHands []hand
	idxHand     int

	initBet     int
	cumWinnings int

	playersTurn bool

	ai Player
}

type Options struct {
	Hands int
	Decks int
}

func New(opts Options) gameState {
	deck := cardDeck.NewDeck(cardDeck.WithMultipleDecks(opts.Decks), cardDeck.WithShuffle())

	gs := gameState{
		nbDecks:     opts.Decks,
		deck:        deck,
		nbRounds:    opts.Hands,
		playerHands: make([]hand, 1),
		playersTurn: true,
	}

	return gs
}

func (gs *gameState) Play(ai Player) int {
	gs.ai = ai
	lowNbCards := 20
	for r := 0; r < gs.nbRounds; r++ {
		shuffled := false
		if len(gs.deck) < lowNbCards {
			shuffleDeck(gs)
			shuffled = true
		}
		playerBets(gs, shuffled)
		distribute(gs)
		playerPlays(gs)
		dealerPlays(gs)
		ai.EndRound(gs.playerHands, gs.dealerHand)
		updateBalance(gs)
		resetGame(gs)
	}
	return gs.cumWinnings
}

func playerBets(gs *gameState, shuffled bool) {
	var b string
	var err error
	for {
		b = gs.ai.Bet(shuffled)
		if gs.initBet, err = strconv.Atoi(b); err == nil {
			if gs.initBet%10 != 0 {
				fmt.Println("Bet must be a multiple of 10.")
			} else {
				break
			}
		}
	}
}

func playerPlays(gs *gameState) {
	for gs.playersTurn {
		pCards := gs.getCurPlayerHand().cards
		cpCards := make([]cardDeck.Card, len(pCards))
		copy(cpCards, pCards)
		move := gs.ai.Play(cpCards, gs.dealerHand[0])
		move(gs)
	}
}

// Type englobing the possible moves a player can play
// It is the return type of the Play method for the Player interface.
type Action func(*gameState) error

// Ask for another card
func Hit(gs *gameState) error {
	drawPlayer(gs)
	if s, _ := Score(gs.playerHands[gs.idxHand].cards); s > 21 {
		gs.idxHand++
	}
	if gs.idxHand >= len(gs.playerHands) {
		gs.playersTurn = false
	}
	return nil
}

// Stop asking for cards
func Stand(gs *gameState) error {
	gs.idxHand++
	if gs.idxHand >= len(gs.playerHands) {
		gs.playersTurn = false
	}
	return nil
}

// If the 2 cards in the player's hand have same rank, the player can split its hand
func Split(gs *gameState) error {
	pHand := gs.getCurPlayerHand()
	if len(pHand.cards) != 2 || pHand.cards[0].Rank != pHand.cards[1].Rank {
		return errors.New("Can only double with 2 cards of same rank.")
	}
	var card1, card2 []cardDeck.Card
	split1 := hand{
		cards: append(card1, pHand.cards[0]),
		bet:   pHand.bet,
	}
	split2 := hand{
		cards: append(card2, pHand.cards[1]),
		bet:   pHand.bet,
	}
	gs.playerHands[gs.idxHand] = split1
	gs.playerHands = append(gs.playerHands, split2)
	return nil
}

// If the player has 2 cards in its hand, it can double its bet and ask for exactly one more card
func Double(gs *gameState) error {
	pCards := gs.getCurPlayerHand().cards
	if len(pCards) != 2 {
		return errors.New("Can only double with 2 cards.")
	}
	gs.playerHands[gs.idxHand].bet *= 2
	Hit(gs)
	Stand(gs)
	return nil
}

func drawDealer(gs *gameState) {
	var c cardDeck.Card
	c, gs.deck = gs.deck[0], gs.deck[1:]
	gs.dealerHand = append(gs.dealerHand, c)
}
func drawPlayer(gs *gameState) {
	var c cardDeck.Card
	c, gs.deck = gs.deck[0], gs.deck[1:]
	gs.playerHands[gs.idxHand].cards = append(gs.playerHands[gs.idxHand].cards, c)
}

func distribute(gs *gameState) {
	for i := 0; i < 2; i++ {
		drawPlayer(gs)
		drawDealer(gs)
	}
	gs.playerHands[0].bet = gs.initBet
}

func dealerPlays(gs *gameState) {
	s, soft := Score(gs.dealerHand)
	for s < 17 || (s == 17 && soft) {
		drawDealer(gs)
		s, soft = Score(gs.dealerHand)
	}
}

func updateBalance(gs *gameState) {
	sd, _ := Score(gs.dealerHand)
	bjd := Blackjack(gs.dealerHand)

	for _, h := range gs.playerHands {
		sp, _ := Score(h.cards)
		bjp := Blackjack(h.cards)

		switch {
		case bjp && bjd:
			// Tie
		case bjd:
			gs.cumWinnings -= h.bet
		case bjp:
			gs.cumWinnings += int(1.5 * float64(h.bet))
		case sp > 21:
			gs.cumWinnings -= h.bet
		case sd > 21:
			gs.cumWinnings += h.bet
		case sp > sd:
			gs.cumWinnings += h.bet
		case sp < sd:
			gs.cumWinnings -= h.bet
		case sp == sd:
			// Tie
		}
	}
}

func resetGame(gs *gameState) {
	gs.dealerHand = nil
	gs.playerHands = make([]hand, 1)
	gs.playersTurn = true
	gs.idxHand = 0
}

// Return true if hand is a blackjack
func Blackjack(h []cardDeck.Card) bool {
	s, _ := Score(h)
	return len(h) == 2 && s == 21
}

// Return the score for a slice of cards and a bool to tell if score is soft or not
func Score(h []cardDeck.Card) (int, bool) {
	s := 0
	soft := false
	for _, c := range h {
		s += min(int(c.Rank), 10)
		if c.Rank == cardDeck.Ace {
			soft = true
		}
	}
	if soft && s < 12 {
		s += 10
	}
	return s, soft
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func (h hand) String() string {
	cards := h.cards
	str := make([]string, len(cards))
	for i, c := range cards {
		str[i] = c.String()
	}
	return strings.Join(str, ", ")
}

// Helper function to shuffle deck in game state
func shuffleDeck(gs *gameState) {
	gs.deck = cardDeck.NewDeck(cardDeck.WithMultipleDecks(gs.nbDecks), cardDeck.WithShuffle())
	fmt.Println("Deck has been shuffled")
}

// Helper function to return current player hand
func (gs *gameState) getCurPlayerHand() hand {
	return gs.playerHands[gs.idxHand]
}
