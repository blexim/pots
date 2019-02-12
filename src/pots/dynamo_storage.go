package pots

import (
	"strconv"
	"time"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoStorage struct {
  svc *dynamodb.DynamoDB
  ledgerTable string
  balanceTable string
}

func GetDynamo() DynamoStorage {
  sess := session.Must(session.NewSession(&aws.Config{
      Region: aws.String("eu-west-2"),
    }))

  return DynamoStorage{dynamodb.New(sess), "pots-ledger", "pots-balance"}
}

func GetTestDynamo() DynamoStorage {
  sess := session.Must(session.NewSession(&aws.Config{
      Region: aws.String("eu-west-2"),
    }))

  return DynamoStorage{dynamodb.New(sess), "pots-test-ledger", "pots-test-balance"}
}

func (s DynamoStorage) Transfer(from string, to string, value int) error {
  if err := s.addLedger(from, to, value); err != nil {
    return err
  }

  if err := s.addBalance(from, -value); err != nil {
    return err
  }

  if err := s.addBalance(to, value); err != nil {
    return err
  }

  return nil
}

func (s DynamoStorage) addLedger(from string, to string, value int) error {
	ledgerEntry := LedgerEntry{from, to, value, time.Now().Unix(),}
	av, err := dynamodbattribute.MarshalMap(ledgerEntry)

  if err != nil {
    return err
  }

	input := &dynamodb.PutItemInput{
    Item: av,
    TableName: aws.String(s.ledgerTable),
	}

  _, err = s.svc.PutItem(input)
  return err
}

func (s DynamoStorage) addBalance(player string, value int) error {
  update := &dynamodb.UpdateItemInput{
    TableName: aws.String(s.balanceTable),
    Key: map[string]*dynamodb.AttributeValue{
      "player": {
        S: aws.String(player),
      },
    },
    UpdateExpression: aws.String("ADD balance :val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {
				N: aws.String(strconv.Itoa(value)),
			},
    },
  }

  _, err := s.svc.UpdateItem(update)
  return err
}

func (s DynamoStorage) GetBalances() ([]BalanceEntry, error) {
  scanParams := &dynamodb.ScanInput{
    TableName: aws.String(s.balanceTable),
    ProjectionExpression: aws.String("player, balance"),
  }

  scanOutput, err := s.svc.Scan(scanParams)

  if err != nil {
    return nil, err
  }

  ret := []BalanceEntry{}
  err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &ret)

  if err != nil {
    return nil, err
  }

  return ret, nil
}

