package pots

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type SqlStorage struct {
  gameId int64
}

func openDb() (*sql.DB, error) {
  return sql.Open("mysql", "root:password@tcp(127.0.0.1)/pots")
}

func NewGame() (SqlStorage, error) {
  db, err := openDb()

  if err != nil {
    return SqlStorage{-1}, err
  }

  defer db.Close()

  res, err := db.Exec("INSERT INTO games VALUES(NULL, NOW(), NULL)")

  if err != nil {
    return SqlStorage{-1}, err
  }

  gameId, err := res.LastInsertId()

  if err != nil {
    return SqlStorage{-1}, err
  }

  return SqlStorage{gameId}, nil
}

func (s SqlStorage) AddCredit(name string, value int) error {
  db, err := openDb()

  if err != nil {
    return err
  }

  _, err = db.Exec("INSERT INTO ledger VALUES(NULL, ?, ?, ?, NOW())", s.gameId, name, value)
  return err
}

func (s SqlStorage) AddDebit(name string, value int) error {
  return s.AddCredit(name, -value)
}

func (s SqlStorage) GetBalances() ([]BalanceEntry, error) {
  db, err := openDb()

  if err != nil {
    return nil, err
  }

  rows, err := db.Query("SELECT player, sum(value) FROM ledger WHERE game_id=? GROUP BY PLAYER",
    s.gameId)

  if err != nil {
    return nil, err
  }

  defer rows.Close()
  ret := make([]BalanceEntry, 0)

  for rows.Next() {
    var name string
    var balance int
    rows.Scan(&name, &balance)
    ret = append(ret, BalanceEntry{name, balance})
  }

  if err = rows.Err(); err != nil {
    return nil, err
  }

  return ret, nil
}

func (s SqlStorage) EndGame() error {
  db, err := openDb()

  if err != nil {
    return err
  }

  defer db.Close()

  _, err = db.Exec("UPDATE games SET end_time=NOW() WHERE game_id=?", s.gameId)
  return err
}

