package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBuyin(t *testing.T) {
	tests := map[string]int{
		"£13.70": 1370,
		"£13":    1300,
		"20":     2000,
		"0.1":    10,
		"0.10":   10,
		"0.01":   1,
		"0.009":  0,
		"0.019":  1,
	}
	for input, expected := range tests {
		actual := ParseBuyin("alice", "@pokerbot in "+input)
		assert.Equal(t, expected, actual.Value, "%s parsed to unexpected value", input)
	}
}
