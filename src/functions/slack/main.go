package main

import (
  "context"
  "encoding/json"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/nlopes/slack"
  "github.com/nlopes/slack/slackevents"
  "io"
  "log"
  "net/http"
  "os"
)

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  log.Printf("Handling a Slack request: %v", req)
  log.Printf("Body: %v", req.Body)

  sv, err := slack.NewSecretsVerifier(http.Header(req.MultiValueHeaders), os.Getenv("SLACK_SIGNING_SECRET"))

  if err != nil {
    log.Printf("Couldn't make verifier: %v", err)
    return events.APIGatewayProxyResponse{
      StatusCode: http.StatusUnauthorized,
      Body: "Auth failed",
    }, nil
  }

  io.WriteString(&sv, req.Body)

  if err = sv.Ensure(); err != nil {
    log.Printf("Bad signature: %v", err)
    return events.APIGatewayProxyResponse{
      StatusCode: http.StatusUnauthorized,
      Body: "Auth failed",
    }, nil
  }

  log.Printf("Signature verification succeeded!")

  msg, err := slackevents.ParseEvent(json.RawMessage(req.Body), slackevents.OptionNoVerifyToken())
  if err != nil {
    log.Printf("Couldn't parse json from Slack: %v", err)
    return events.APIGatewayProxyResponse{
      StatusCode: http.StatusBadRequest,
      Body: "Bad request",
    }, nil
  }

  if msg.InnerEvent.Type == "app_mention" {
    log.Printf("Parsed mention: %v", msg.InnerEvent.Data.(slackevents.AppMentionEvent))
  } else {
    log.Printf("Message type was [%v]", msg.InnerEvent.Type)
  }

  return events.APIGatewayProxyResponse{
    StatusCode: http.StatusOK,
    Body: "ok",
  }, nil
}

func main() {
  lambda.Start(HandleRequest)
}
