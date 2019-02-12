package pots

import "fmt"

func ExampleSettle1() {
  people := []BalanceEntry{
    {
      name: "alice",
      balance: 10,
    },
    {
      name: "bob",
      balance: -10,
    },
  }

  transfers := Settle(people)

  fmt.Println(transfers)
  // Output: [{bob alice 10}]
}

func ExampleSettle2() {
  people := []BalanceEntry{
    {
      name: "alice",
      balance: 10,
    },
    {
      name: "bob",
      balance: -7,
    },
    {
      name: "charlie",
      balance: 5,
    },
    {
      name: "dave",
      balance: -8,
    },
  }

  transfers := Settle(people)

  fmt.Println(transfers)
  // Output: [{dave alice 8} {bob alice 2} {bob charlie 5}]
}
