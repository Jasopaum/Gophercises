package blackjack

import (
	"fmt"
	"gophercises/cardDeck"
	"strconv"
	"strings"
)

type Hand []cardDeck.Card

type gameState struct {
	dealerHand  Hand
	playerHand  Hand
	nbDecks     int
	deck        []cardDeck.Card
	nbRounds    int
	curBet      int
	cumWinnings int
	playersTurn bool
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
		playersTurn: true,
	}

	return gs
}

func (gs *gameState) Play(ai Player) int {
	lowNbCards := 20
	for r := 0; r < gs.nbRounds; r++ {
		shuffled := false
		if len(gs.deck) < lowNbCards {
			gs.deck = cardDeck.NewDeck(cardDeck.WithMultipleDecks(gs.nbDecks), cardDeck.WithShuffle())
			fmt.Println("Deck has been shuffled")
			shuffled = true
		}
		var b string
		var err error
		for {
			b = ai.Bet(shuffled)
			if gs.curBet, err = strconv.Atoi(b); err == nil {
				break
			}
		}
		distribute(gs)
		for gs.playersTurn {
			pHand := make([]cardDeck.Card, len(gs.playerHand))
			copy(pHand, gs.playerHand)
			move := ai.Play(pHand, gs.dealerHand[0])
			move(gs)
		}
		dealerPlays(gs)
		showResult(gs)
		resetGame(gs)
	}
	return gs.cumWinnings
}

type action func(*gameState)

func Hit(gs *gameState) {
	gs.playerHand, gs.deck = draw(gs.playerHand, gs.deck)
	if s, _ := Score(gs.playerHand); s > 21 {
		gs.playersTurn = false
	}
}
func Stand(gs *gameState) {
	gs.playersTurn = false
}
func Double(gs *gameState) {
	if len(gs.playerHand) > 2 {
		fmt.Println("Can only double with 2 cards.")
		return
	}
	gs.curBet *= 2
	Hit(gs)
	gs.playersTurn = false
}

func draw(h Hand, d []cardDeck.Card) (Hand, []cardDeck.Card) {
	c, d := d[0], d[1:]
	h = append(h, c)
	return h, d
}

func distribute(gs *gameState) {
	for i := 0; i < 2; i++ {
		gs.playerHand, gs.deck = draw(gs.playerHand, gs.deck)
		gs.dealerHand, gs.deck = draw(gs.dealerHand, gs.deck)
	}
}

func dealerPlays(gs *gameState) {
	s, soft := Score(gs.dealerHand)
	for s < 17 || (s == 17 && soft) {
		gs.dealerHand, gs.deck = draw(gs.dealerHand, gs.deck)
		s, soft = Score(gs.dealerHand)
	}
}

func showResult(gs *gameState) {
	sd, _ := Score(gs.dealerHand)
	sp, _ := Score(gs.playerHand)
	bjd := Blackjack(gs.dealerHand)
	bjp := Blackjack(gs.playerHand)

	fmt.Println("Dealer:")
	fmt.Println("\t", gs.dealerHand)
	fmt.Println("\tscore:", sd)

	fmt.Println("Player:")
	fmt.Println("\t", gs.playerHand)
	fmt.Println("\tscore:", sp)

	if bjp && bjd {
		fmt.Println("Tie")
	} else if bjd {
		fmt.Println("Dealer won")
		gs.cumWinnings -= gs.curBet
	} else if bjp {
		fmt.Println("You won")
		gs.cumWinnings += int(1.5 * float64(gs.curBet))
	} else if sp > 21 {
		fmt.Println("You busted")
		gs.cumWinnings -= gs.curBet
	} else if sd > 21 {
		fmt.Println("Dealer busted")
		gs.cumWinnings += gs.curBet
	} else if sp > sd {
		fmt.Println("You won")
		gs.cumWinnings += gs.curBet
	} else if sp < sd {
		fmt.Println("Dealer won")
		gs.cumWinnings -= gs.curBet
	} else if sp == sd {
		fmt.Println("Tie")
	}
}

func resetGame(gs *gameState) {
	gs.dealerHand = nil
	gs.playerHand = nil
	gs.playersTurn = true
}

func Blackjack(h Hand) bool {
	s, _ := Score(h)
	return len(h) == 2 && s == 21
}

func Score(h Hand) (int, bool) {
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

func (h Hand) String() string {
	str := make([]string, len(h))
	for i, c := range h {
		str[i] = c.String()
	}
	return strings.Join(str, ", ")
}

func (h Hand) StringDealer() string {
	return fmt.Sprintf("%s, ***Hidden***", h[0])
}
