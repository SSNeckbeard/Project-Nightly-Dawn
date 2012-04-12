// dice.go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Creates a closure with a seeded generator for a specific die type
func DieType (NumSides int) func() (sum int, rolls []int) {
	generator := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	return func () (sum int, rolls []int) {
		sum = 0
		roll := generator.Intn(NumSides) + 1
		for {
			sum += roll
			rolls = append (rolls, roll)
			if roll != NumSides {
				break
			}
			roll = generator.Intn(NumSides) + 1
		}
		return
	}	
}

// 

func main() {
	roll := 0
	list := []int{}
	d4   := DieType( 4)
	d6   := DieType( 6)
	d8   := DieType( 8)
	d10  := DieType(10)
	d12  := DieType(12)
	
	for i := 0; i < 16; i++ {
		roll, list = d4()
		fmt.Println("Rolling 1d4: ", roll, list)
	}
	fmt.Println()
	for i := 0; i < 16; i++ {
		roll, list = d6()
		fmt.Println("Rolling 1d6: ", roll, list)
	}
	fmt.Println()
	for i := 0; i < 16; i++ {
		roll, list = d8()
		fmt.Println("Rolling 1d8: ", roll, list)
	}
	fmt.Println()
	for i := 0; i < 16; i++ {
		roll, list = d10()
		fmt.Println("Rolling 1d10: ", roll, list)
	}
	fmt.Println()
	for i := 0; i < 16; i++ {
		roll, list = d12()
		fmt.Println("Rolling 1d12: ", roll, list)
	}
}
