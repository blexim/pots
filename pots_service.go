package pots

type Storage interface {
  AddCredit(name string, value int)
  AddDebit(name string, value int)
  GetBalances() []BalanceEntry
}

type PotsService {
  storage Storage
}

func (service PotsService) AddCredit(name string, value int) {
}
