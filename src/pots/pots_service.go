package pots

type Storage interface {
  Transfer(from string, to string, value int) error
  GetBalances() ([]BalanceEntry, error)
}

type PotsService struct {
  storage Storage
}

const POT_ACCOUNT string = "pot"

func (service PotsService) AddCredit(name string, value int) error {
  return service.storage.Transfer(POT_ACCOUNT, name, value)
}

func (service PotsService) AddDebit(name string, value int) error {
  return service.storage.Transfer(name, POT_ACCOUNT, value)
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
