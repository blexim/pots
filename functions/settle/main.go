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
  Transfers []pots.LedgerEntry    `json:"transfers"`
}

var service pots.PotsService

func init() {
  service = pots.GetDynamoPotsService()
}

func HandleRequest(ctx context.Context) (SettleResponse, error) {
  log.Print("Handling a settle request")

  if entries, err := service.Settle(); err != nil {
    log.Printf("Got an error: %v", err)
    return SettleResponse{}, err
  } else {
    log.Printf("Found a solution: %v", entries)
    return SettleResponse{entries}, nil
  }
}

func main() {
  lambda.Start(HandleRequest)
}
