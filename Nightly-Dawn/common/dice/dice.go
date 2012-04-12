// dice.go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Represents a single die roll (with explosions)
// NOTE: Should the struct be modified in order to mention if it has an explosion?
// This can be figured out if list has more than one element.
type DieRoll struct {
	Total uint16
	List  []uint16
}

// Creates a closure with a seeded generator for a specific die type
func DieType (NumSides int) func() (res *DieRoll) {
	generator := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	return func () (res *DieRoll) {
		res       = new(DieRoll)
		res.Total = 0
		roll     := uint16(generator.Intn(NumSides) + 1)
		for {
			res.Total += roll
			res.List   = append (res.List, roll)
			if roll != uint16(NumSides) {
				break
			}
			roll = uint16(generator.Intn(NumSides) + 1)
		}
		return
	}	
}

// Rolls a die a set number of times. Returns as an array of DieRoll structs
// Report without the explosion unless explosion array has two or more elements
func RollDice (die func() *DieRoll, NumDice uint8) []DieRoll {
	rolls := make([]DieRoll, NumDice)
	for i := uint8(0); i < NumDice; i++ {
		rolls[i] = *die()
	}
	return rolls
}

func main() {
	roll := new(DieRoll)
	d4   := DieType( 4)
	d6   := DieType( 6)
	d8   := DieType( 8)
	d10  := DieType(10)
	d12  := DieType(12)
	
	for i := 0; i < 16; i++ {
		roll = d4()
		fmt.Println("Rolling 1d4: ", roll.Total, roll.List)
	}
	fmt.Println()
	for i := 0; i < 16; i++ {
		roll = d6()
		fmt.Println("Rolling 1d6: ", roll.Total, roll.List)
	}
	fmt.Println()
	for i := 0; i < 16; i++ {
		roll = d8()
		fmt.Println("Rolling 1d8: ", roll.Total, roll.List)
	}
	fmt.Println()
	for i := 0; i < 16; i++ {
		roll = d10()
		fmt.Println("Rolling 1d10: ", roll.Total, roll.List)
	}
	fmt.Println()
	for i := 0; i < 16; i++ {
		roll = d12()
		fmt.Println("Rolling 1d12: ", roll.Total, roll.List)
	} 
	fmt.Println()
	fmt.Println("4d12", RollDice(d12, 4))
}
