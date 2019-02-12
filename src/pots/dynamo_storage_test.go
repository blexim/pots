package pots

import (
  "fmt"
  "testing"
)

func TestDynamoCredit(t *testing.T) {
  s := GetDynamo()

  if err := s.Transfer("@alice", "@bob", 10); err != nil {
    t.Errorf("Error: %v", err)
  }
}

func ExampleDynamoBalances() {
  s := GetDynamo()

  fmt.Println(s.GetBalances())

  // Output: {}
}
