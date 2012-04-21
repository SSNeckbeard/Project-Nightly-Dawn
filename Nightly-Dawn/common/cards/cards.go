package main

import (
	"fmt"
	"math/rand"
	"time"
	"sort"
)

////////////////////////////////////////////////////////////////////////////////
// Constants
const Jack  = 11
const Queen = 12
const King  = 13
const Ace   = 14
const Joker = 15

const Spades   = 0
const Hearts   = 1
const Diamonds = 2
const Clubs    = 3

var ValueNames = []string {"Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King", "Ace", "Joker"}
	// Value - 2 is the index into this array
var SuitNames  = []string {"Spades", "Hearts", "Diamonds", "Clubs"}
	// Suit is index into this array

////////////////////////////////////////////////////////////////////////////////
// Data Structures
type Card struct {
	value   int  // 2-10 standard value, 11 = Jack, 12 = Queen, 13 = King, 14 = Ace, 15 = Joker
	suit    int  // 0 = Spades, 1 = Hearts, 2 = Diamonds, 3 = Clubs
	sleeved bool // A sleeved card is considered equal to all other cards
	// NOTE: Should this be equal to or greater than all other cards?
}

type Cards []Card
   // QUESTION: Should this be of type []*Card?

// A deck is a slice of cards which contains up to 54 cards and a random number
// generator to decide which one is drawn when queried.
// You should only use this struct after getting one from the InitDeck function.
type Deck struct {
	deck      Cards
	generator *rand.Rand 
           // random number generator will be seeded in InitDeck function
}

type ErrMtDeck int
   // 0: Drawing from empty deck
   // 1: Trying to draw more cards than deck contains, reshuffling not enabled
 
func (e ErrMtDeck) Error() string {
	if e == 0 {
		return "ERROR: Cannot draw cards from an empty deck." }
	if e == 1 {
		return "ERROR: Deck does not contain enough cards to draw." }
	return "Error: Deck is empty."
}

func InitDeck() *Deck {
	d := new(Deck)
	d.deck      = make([]Card, 54, 54)
	d.generator = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	d.Shuffle()
	return d
}

////////////////////////////////////////////////////////////////////////////////
// Shuffling and Drawing Cards

// Refills the slice with all possible cards it can hold. 
// The deck will not be shuffled in the conventional sense, but the random draw 
// can get any card.
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

// Gets a random card from the possibilities within the deck, removes it from
// available possibilities.
// Returns ErrMtDeck if there are no available possibilities for this deck.
func (deck *Deck) DrawCard() (Card, error) {
	if len(deck.deck) == 0 {
		return Card{Joker + 1, Clubs + 1, true}, ErrMtDeck(0) }
	i   := deck.generator.Intn(len(deck.deck))
	res := deck.deck[i]
	deck.deck = append(deck.deck[:i], deck.deck[i+1:]...)
	return res, nil
}

// Gets x random cards from the possibilities within the deck, removes them from
// available possibilities.
// Question: if x is greater than the number of cards currently within the deck, what then?
// Option A: shuffle deck, redraw the remaining  * Defaulting to this
// Option B: stop with empty deck
// Option C: don't draw anything
func (deck *Deck) DrawXCards(x int, redraw bool) (Cards, error) {
	res   := make(Cards, x, x)
	var   err error
	
	if x > len(deck.deck) && !redraw { return res, ErrMtDeck(1) }
	for i := 0; i < x; i++ {
		if len(deck.deck) == 0 { deck.Shuffle() }
		res[i], err = deck.DrawCard()
		if err != nil { return res, err }
	}
	sort.Sort(CardSorter{res})
	return res, nil
}

////////////////////////////////////////////////////////////////////////////////
// Utils

func (card Card) ToString () string {
	if !IsValid(card) {
		return "*Invalid Card (" + string(card.value) + ", " + string(card.suit) + ")*"
	}
	res := ""
	if card.sleeved { res = "*Sleeved* " }
	if card.value == Joker {
		if card.IsRed() { return res + "Red Joker" }
		return res + "Black Joker"
	}
	return res + ValueNames[card.value - 2] + " of " + SuitNames[card.suit]
}

// Only use cards which pass this check.
func IsValid (card Card) bool {
	return card.value >= 2 && card.value <= Joker && card.suit >= Spades && card.suit <= Clubs
}

func (card Card) IsRed() bool {
	return card.suit == Diamonds || card.suit == Hearts
}

func (card Card) IsDroppable() bool {
	// droppable iff != a deuce or a joker
	return card.value != 2 && card.value != Joker
}

////////////////////////////////////////////////////////////////////////////////
// Sorting
func (s Cards) Len() int      { return len(s) }
func (s Cards) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// CardSorter implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded Cards value.
type CardSorter struct{ Cards }

func (s CardSorter) Less(i, j int) bool { 
	return Compare(s.Cards[i], s.Cards[j]) == -1 }

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


////////////////////////////////////////////////////////////////////////////////
// Tests
// @TODO: Make tests

////////////////////////////////////////////////////////////////////////////////
// Demo main
func main() {
	d := InitDeck()
	fmt.Println("deck.deck:", d.deck)
	sort.Sort(CardSorter{d.deck})
	fmt.Println("Sorted deck.deck:", d.deck)
	//for i := 0; i < len(d.deck); i++{
	//	fmt.Println(d.deck[i].ToString())
	//}
	
	//fmt.Println()
	for i := 0; i < 60; i++ {
		c, e := d.DrawCard()
		if e == nil {		
			fmt.Printf("Draw %d: %s\n", i + 1, c.ToString())
		} else {
			fmt.Println(e.Error())
			fmt.Printf("\tReshuffling\n")
			d.Shuffle()
			i--
		}
		
	}
	// to sort: sort.Sort(CardSorter{<slice>})
	
	for i:= 0; i < 2; i++ {
		fmt.Println()
		hand, e := d.DrawXCards(5, true)
		if e != nil {
			fmt.Println(e)
			break
		}
		for j:= 0; j < len(hand); j++ {
			fmt.Printf("\tCard %d: %s\n", j + 1, hand[j].ToString())
		}	
	}
	
	
	
}
