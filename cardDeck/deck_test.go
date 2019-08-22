package cardDeck

import (
	"fmt"
	"testing"
)

func ExampleCard_String() {
	c := Card{Six, Heart}
	fmt.Println(c.String())

	// Output: Six of Hearts
}

func TestDeck_New(t *testing.T) {
	d := NewDeck()
	if len(d) != 52 {
		t.Errorf("Expected 52 cards, got:%d", len(d))
	}
}
