package main

import (
  "context"
  "log"
  "pots"
  "github.com/aws/aws-lambda-go/lambda"
)

type TransferRequest struct {
  From string`json:"from"`
  To string`json:"to"`
  Value int`json:"value"`
}

type TransferResponse struct {
}

var service pots.PotsService

func init() {
  service = pots.GetDynamoPotsService()
}

func HandleRequest(ctx context.Context, req TransferRequest) (TransferResponse, error) {
  log.Printf("Handling a transfer request: %v", req)

  var err error

  if req.From == "" {
    err = service.AddCredit(req.To, req.Value)
  } else if req.To == "" {
    err = service.AddDebit(req.From, req.Value)
  } else {
    err = service.Transfer(req.From, req.To, req.Value)
  }

  return TransferResponse{}, err
}

func main() {
  lambda.Start(HandleRequest)
}

