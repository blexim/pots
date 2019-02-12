package main

import (
  "context"
  "log"
  "pots"
  "github.com/aws/aws-lambda-go/lambda"
)

type TransferRequest struct {
  from string `json:"from"`
  to string `json:"to"`
  value int `json:"value"`
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

  if req.from == "" {
    err = service.AddCredit(req.to, req.value)
  } else if req.to == "" {
    err = service.AddDebit(req.from, req.value)
  } else {
    err = service.Transfer(req.from, req.to, req.value)
  }

  return TransferResponse{}, err
}

func main() {
  lambda.Start(HandleRequest)
}

