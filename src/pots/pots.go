package pots

import (
  "sort"
)

type LedgerEntry struct {
  From string`json:"from"`
  To string`json:"to"`
  Value int`json:"value"`
  Timestamp int64`json:"timestamp"`
}

type BalanceEntry struct {
  Player string`json:"player"`
  Balance int`json:"balance"`
}

type byAbsBalance []BalanceEntry

func absi(i int) int { if i < 0 { return -i } else { return i } }
func mini(i, j int) int { if i < j { return i } else { return j } }

func (a byAbsBalance) Len() int { return len(a) }
func (a byAbsBalance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byAbsBalance) Less(i, j int) bool { return absi(a[i].Balance) > absi(a[j].Balance) }

func Settle(balances []BalanceEntry) []LedgerEntry {
  var creditors, debtors []BalanceEntry

  for _, entry := range balances {
    if entry.Balance < 0 {
      creditors = append(creditors, entry)
    } else {
      debtors = append(debtors, entry)
    }
  }

  sort.Sort(byAbsBalance(creditors))
  sort.Sort(byAbsBalance(debtors))

  var transfers []LedgerEntry

  for ci, di := 0, 0; ci < len(creditors) && di < len(debtors); {
    c := &creditors[ci]
    d := &debtors[di]
    value := mini(-c.Balance, d.Balance)
    transfer := LedgerEntry{ d.Player, c.Player, value, 0, }
    transfers = append(transfers, transfer)
    c.Balance += value
    d.Balance -= value

    if c.Balance == 0 {
      ci++
    }

    if d.Balance == 0 {
      di++
    }
  }

  return transfers
}

