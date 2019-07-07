package pots

import "fmt"

func ExampleSettle1() {
	people := []BalanceEntry{
		{
			Player:  "alice",
			Balance: 10,
		},
		{
			Player:  "bob",
			Balance: -10,
		},
	}

	transfers := Settle(people)

	fmt.Println(transfers)
	// Output: [{alice bob 10 0}]
}

func ExampleSettle2() {
	people := []BalanceEntry{
		{
			Player:  "alice",
			Balance: 10,
		},
		{
			Player:  "bob",
			Balance: -7,
		},
		{
			Player:  "charlie",
			Balance: 5,
		},
		{
			Player:  "dave",
			Balance: -8,
		},
	}

	transfers := Settle(people)

	fmt.Println(transfers)
	// Output: [{alice dave 8 0} {alice bob 2 0} {charlie bob 5 0}]
}
