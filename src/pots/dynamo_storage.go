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
}

func GetDb() *dynamodb.DynamoDB {
  sess := session.Must(session.NewSession(&aws.Config{
      Region: aws.String("eu-west-2")},
  ))

  // Create DynamoDB client
  return dynamodb.New(sess)
}

func (s DynamoStorage) AddCredit(name string, value int) error {
  svc := GetDb()
	ledgerEntry := LedgerEntry{name, "pot", value, time.Now().Unix(),}
	av, err := dynamodbattribute.MarshalMap(ledgerEntry)

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

  update := &dynamodb.UpdateItemInput{
    TableName: aws.String("pots-balance"),
    Key: map[string]*dynamodb.AttributeValue{
      "player": {
        S: aws.String(name),
      },
    },
    UpdateExpression: aws.String("ADD balance :val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {
				N: aws.String(strconv.Itoa(value)),
			},
    },
  }

  _, err = svc.UpdateItem(update)

  return err
}

