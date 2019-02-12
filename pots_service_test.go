package pots

import (
  "fmt"
)

func ExampleSettle() {
  sqlStorage, _ := NewGame()
  service := PotsService{sqlStorage}

  service.AddCredit("alice", 10)
  service.AddCredit("alice", 5)
  service.AddDebit("bob", 15)

  transfers, _ := service.Settle()

  fmt.Println(transfers)

  // Output: [{bob alice 15}]
}

