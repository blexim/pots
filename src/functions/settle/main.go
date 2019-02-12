package main

import (
  "context"
  "log"
  "pots"
  "github.com/aws/aws-lambda-go/lambda"
)

type SettleRequest struct {
}

type SettleResponse struct {
  transfers []pots.LedgerEntry  `json:"transfers"`
}

var service pots.PotsService

func init() {
  service = pots.GetDynamoPotsService()
}

func HandleRequest(ctx context.Context) (SettleResponse, error) {
  log.Print("Handling a settle request")

  if entries, err := service.Settle(); err != nil {
    return SettleResponse{entries}, nil
  } else {
    return SettleResponse{}, err
  }
}

func main() {
  lambda.Start(HandleRequest)
}

