package main

import (
	"fmt"
	"math/rand"
	"time"
)

const Jack  = 11
const Queen = 12
const King  = 13
const Ace   = 14
const Joker = 15

const Spades   = 0
const Hearts   = 1
const Diamonds = 2
const Clubs    = 3

// Value - 2 is the index into this array
var ValueNames = []string {"Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King", "Ace", "Joker"}
// Suit is index into this array
var SuitNames  = []string {"Spades", "Hearts", "Diamonds", "Clubs"}


type Card struct {
	value   int  // 2-10 standard value, 11 = Jack, 12 = Queen, 13 = King, 14 = Ace, 15 = Joker
	suit    int  // 0 = Spades, 1 = Hearts, 2 = Diamonds, 3 = Clubs
	sleeved bool // A sleeved card is considered equal to all other cards
	// NOTE: Should this be equal to or greater than all other cards?
}

func (card Card) IsRed() bool {
	return card.suit == Diamonds || card.suit == Hearts
}

func (card Card) IsDroppable() bool {
	// droppable iff != a deuce or a joker
	return card.value != 2 && card.value != Joker
}

// Standard compare function, returns -1 if a < b, 0 if a == b, 1 if a > b
// NOTE: Sleeved cards are currently considered equal to all other cards
//       This may change so that they are always considered greater
// Red Joker is considered greater than everything but another RJ or Sleeved Card
// Black Joker is considered less than everything but another BJ or Sleeved Card
func Compare(a Card, b Card) int {
	// Check sleeved or equal
	if (a.sleeved || b.sleeved) || (a.value == b.value && a.suit == b.suit) {
		return 0
	}
	if a.value == Joker && (!a.IsRed()) { // A is Black Joker
		return -1
	}
	if b.value == Joker && (!b.IsRed()) { // B is Black Joker
		return 1
	}

	if a.value > b.value {
		return 1
	} else if a.value < b.value {
		return -1
	} else {
		if a.suit < b.suit {
			return 1
		} else if a.suit > b.suit {
			return -1
		}
	}
	return 0 // we should never get here
}

// A deck is a slice of cards which contains up to 54 cards
// Will this create a slice with capacity 54 or something that always reports length of 54?
type Deck struct {
	deck      []Card
	generator *rand.Rand // random number generator will be seeded in NewDeck function
}

func InitDeck() *Deck {
	d := new(Deck)
	d.deck = make([]Card, 54, 54)
	generator := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	d.generator = generator
	d.Shuffle()
	return d
}

// refills the array with all possible cards it can hold. 
// The deck will not be shuffled in the conventional sense, but the random draw can get any card
func (deck *Deck) Shuffle() {
	// Reslice to capacity (54 cards)
	deck.deck = deck.deck[:cap(deck.deck)]

	// add jokers
	deck.deck[0] = Card{Joker, Spades, false}
	deck.deck[1] = Card{Joker, Hearts, false}
	// other values
	values := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, Jack, Queen, King, Ace}
	curr := 2
	for i := 0; i < 13; i++ { // all non-joker values
		for j := 0; j < 4; j++ { // all suits
			deck.deck[curr] = Card{values[i], j, false}
			// curr++ doesn't seem to work
			curr = curr + 1
		}
	}
	return
}


// @TODO: Finish this
// What to do on error?
func (deck *Deck) DrawCard() Card {
	if len(deck.deck) == 0 {
		return Card{Joker + 1, Clubs + 1, true}
	}
	i   := deck.generator.Intn(len(deck.deck))
	res := deck.deck[i]
	deck.deck = append(deck.deck[:i], deck.deck[i+1:]...)
	return res
}
// deleting the ith item of a slice s = append(s[:i], s[i+1:]...) 
/*
// @TODO: Finish this
// Question: if x is greater than the number of cards currently within the deck, what then?
// Option A: shuffle deck, redraw the remaining
// Option B: stop with empty deck
func (deck *Deck) DrawXCards(x int) []Card {
	
}
*/
func (card Card) ToString () string {
	if !IsValid(card) {
		return "*Invalid Card (" + string(card.value) + ", " + string(card.suit) + ")*"
	}
	if card.value == Joker {
		if card.IsRed() {
			return "Red Joker"
		}
		return "Black Joker"
	}
	// do in the form "blah" of "blah"
	
	return ValueNames[card.value - 2] + " of " + SuitNames[card.suit]
}

func IsValid (card Card) bool {
	return card.value >= 2 && card.value <= Joker && card.suit >= Spades && card.suit <= Clubs
}

//func SortCards (hand []Card) []Card {
//	
//}

func main() {
	d := InitDeck()
	fmt.Println("deck.deck:", d.deck)
	for i := 0; i < len(d.deck); i++{
		fmt.Println(d.deck[i].ToString())
	}
	
	fmt.Println()
	for i := 0; i < 10; i++ {
		c := d.DrawCard()		
		fmt.Printf("Draw %d: %s\n", i, c.ToString())
	}

}
