package pots

import "fmt"

func ExampleSettle1() {
  people := []BalanceEntry{
    {
      Player: "alice",
      Balance: 10,
    },
    {
      Player: "bob",
      Balance: -10,
    },
  }

  transfers := Settle(people)

  fmt.Println(transfers)
  // Output: [{bob alice 10}]
}

func ExampleSettle2() {
  people := []BalanceEntry{
    {
      Player: "alice",
      Balance: 10,
    },
    {
      Player: "bob",
      Balance: -7,
    },
    {
      Player: "charlie",
      Balance: 5,
    },
    {
      Player: "dave",
      Balance: -8,
    },
  }

  transfers := Settle(people)

  fmt.Println(transfers)
  // Output: [{dave alice 8} {bob alice 2} {bob charlie 5}]
}
