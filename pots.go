package pots

import (
  "sort"
)

type BalanceEntry struct {
  name string
  balance int
}

type Transfer struct {
  from, to string
  value int
}

type byAbsBalance []BalanceEntry

func absi(i int) int { if i < 0 { return -i } else { return i } }
func mini(i, j int) int { if i < j { return i } else { return j } }

func (a byAbsBalance) Len() int { return len(a) }
func (a byAbsBalance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byAbsBalance) Less(i, j int) bool { return absi(a[i].balance) > absi(a[j].balance) }

func Settle(balances []BalanceEntry) []Transfer {
  var creditors, debtors []BalanceEntry

  for _, entry := range balances {
    if entry.balance > 0 {
      creditors = append(creditors, entry)
    } else {
      debtors = append(debtors, entry)
    }
  }

  sort.Sort(byAbsBalance(creditors))
  sort.Sort(byAbsBalance(debtors))

  var transfers []Transfer

  for ci, di := 0, 0; ci < len(creditors) && di < len(debtors); {
    c := &creditors[ci]
    d := &debtors[di]
    value := mini(c.balance, -d.balance)
    transfer := Transfer{ d.name, c.name, value, }
    transfers = append(transfers, transfer)
    c.balance -= value
    d.balance += value

    if c.balance == 0 {
      ci++
    }

    if d.balance == 0 {
      di++
    }
  }

  return transfers
}

