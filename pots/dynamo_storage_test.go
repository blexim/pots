package pots

import (
	"fmt"
	"testing"
)

func TestDynamoCredit(t *testing.T) {
	s := GetTestDynamo()

	if err := s.Transfer("@alice", "@bob", 10); err != nil {
		t.Errorf("Error: %v", err)
	}
}

func ExampleDynamoBalances() {
	s := GetTestDynamo()

	balances, err := s.GetBalances()

	if err != nil {
		fmt.Println(balances)
	}
}
