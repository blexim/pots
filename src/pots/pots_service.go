package pots

type Storage interface {
  AddCredit(name string, value int) error
  AddDebit(name string, value int) error
  GetBalances() ([]BalanceEntry, error)
  EndGame() error
}

type PotsService struct {
  storage Storage
}

func (service PotsService) AddCredit(name string, value int) error {
  return service.storage.AddCredit(name, value)
}

func (service PotsService) AddDebit(name string, value int) error {
  return service.storage.AddDebit(name, value)
}

func (service PotsService) Settle() ([]LedgerEntry, error) {
  balances, err := service.storage.GetBalances()

  if err != nil {
    return nil, err
  }

  if err = service.storage.EndGame(); err != nil {
    return nil, err
  }

  return Settle(balances), nil
}
