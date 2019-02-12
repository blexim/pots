package pots

import (
  "os"
	"time"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

)

type DynamoStorage struct {
}

type LedgerEntry struct {
	From string`json:"from"`
	To string`json:"to"`
	Value int`json:"value"`
	Timestamp int`json:"timestamp"`
}

type BalanceEntry struct {
	Player string`json:"player"`
	Balance int`json:"balance"`
}

sess, err := session.NewSession(&aws.Config{
    Region: aws.String("eu-west-2")},
)

// Create DynamoDB client
svc := dynamodb.New(sess)

func (s DynamoStorage) AddCredit(name string, value int) error {
  now := time.Now()
	ledgerEntry := LedgerEntry{name, "pot", value, time.Now().Unix(),}
	av, err := dynamodbattribute.MarshalMap(item)

  if err != nil {
    return err
  }

	input := &dynamodb.PutItemInput{
    Item: av,
    TableName: aws.String("pots-ledger"),
	}

  if _, err = svc.PutItem(input); err != nil {
    return err
  }

  input := &dynamodb.UpdateItemInput{
    TableName: aws.String("pots-balances"),
    Key: map[string]*dynamodb.AttributeValue{
      "player": {
        S: aws.String(name),
      },
    },
    UpdateExpression: "ADD balance :v",
  }
}

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

