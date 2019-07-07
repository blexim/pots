package pots

import (
	"fmt"
)

func ExampleSettle() {
	storage := GetTestStorage()
	service := PotsService{storage}

	service.AddCredit("alice", 10)
	service.AddCredit("alice", 5)
	service.AddDebit("bob", 15)

	transfers, _ := service.Settle()

	fmt.Println(transfers)

	// Output: [{alice bob 15 0}]
}
