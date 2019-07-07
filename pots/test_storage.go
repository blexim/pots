package pots

type TestStorage struct {
	ledger   []LedgerEntry
	balances map[string]int
}

func (s TestStorage) Transfer(from string, to string, value int) error {
	s.ledger = append(s.ledger, LedgerEntry{from, to, value, 0})

	if _, ok := s.balances[from]; ok {
		s.balances[from] -= value
	} else {
		s.balances[from] = -value
	}

	if _, ok := s.balances[to]; ok {
		s.balances[to] += value
	} else {
		s.balances[to] = value
	}

	return nil
}

func (s TestStorage) GetBalances() ([]BalanceEntry, error) {
	ret := []BalanceEntry{}

	for player, balance := range s.balances {
		ret = append(ret, BalanceEntry{player, balance})
	}

	return ret, nil
}

func GetTestStorage() TestStorage {
	return TestStorage{nil, make(map[string]int)}
}
