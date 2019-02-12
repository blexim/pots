package pots

import (
  "fmt"
  "testing"
)

func TestNewGame(t *testing.T) {
  sqlStorage, err := NewGame()

  if err != nil || sqlStorage.gameId < 0 {
    t.Errorf("Error: %v", err)
  }
}

func ExampleGetBalances() {
  s, _ := NewGame()

  s.AddCredit("alice", 10)
  s.AddCredit("alice", 5)
  s.AddDebit("bob", 15)

  balances, _ := s.GetBalances()

  fmt.Println(balances)

  // Output: [{alice 15} {bob -15}]
}

