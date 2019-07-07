package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"

	"github.com/blexim/pots/pots"
)

var potsService pots.PotsService
var slackClient *slack.Client

func init() {
	potsService = pots.GetDynamoPotsService()
	slackClient = slack.New(os.Getenv("SLACK_BOT_TOKEN"))
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sv, err := slack.NewSecretsVerifier(http.Header(req.MultiValueHeaders), os.Getenv("SLACK_SIGNING_SECRET"))

	if err != nil {
		log.Printf("Couldn't make verifier: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Auth failed",
		}, nil
	}

	io.WriteString(&sv, req.Body)

	if err = sv.Ensure(); err != nil {
		log.Printf("Bad signature: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Auth failed",
		}, nil
	}

	msg, err := slackevents.ParseEvent(json.RawMessage(req.Body), slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Printf("Couldn't parse json from Slack: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad request",
		}, nil
	}

	if msg.InnerEvent.Type == "app_mention" {
		mention := msg.InnerEvent.Data.(*slackevents.AppMentionEvent)
		if err = handleMention(mention.User, mention.Channel, mention.Text); err != nil {
			log.Printf("Error handling mention: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Internal error",
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "ok",
	}, nil
}

func handleMention(user string, channel string, text string) error {
	log.Printf("Handling mention from [%v] in channel [%v]: [%v]", user, channel, text)

	if buyin := ParseBuyin(user, text); buyin != nil {
		return handleBuyin(channel, buyin)
	}

	if cashout := ParseCashout(user, text); cashout != nil {
		return handleCashout(channel, cashout)
	}

	if transfer := ParseTransfer(user, text); transfer != nil {
		return handleTransfer(channel, transfer)
	}

	if settle := ParseSettle(user, text); settle != nil {
		return handleSettle(channel)
	}

	if help := ParseHelp(user, text); help != nil {
		return handleHelp(channel)
	}

	return nil
}

func handleBuyin(channel string, buyin *BuyinCommand) error {
	if err := potsService.AddCredit(buyin.Player, buyin.Value); err != nil {
		return err
	}

	reply := fmt.Sprintf("Confirmed <@%s> borrowed £%.02f from the bank", buyin.Player, float32(buyin.Value)/100)
	_, _, err := slackClient.PostMessage(channel, slack.MsgOptionText(reply, false))
	return err
}

func handleCashout(channel string, cashout *CashoutCommand) error {
	if err := potsService.AddDebit(cashout.Player, cashout.Value); err != nil {
		return err
	}

	reply := fmt.Sprintf("Confirmed <@%s> deposited £%.02f into the bank", cashout.Player, float32(cashout.Value)/100)
	_, _, err := slackClient.PostMessage(channel, slack.MsgOptionText(reply, false))
	return err
}

func handleTransfer(channel string, transfer *TransferCommand) error {
	if err := potsService.Transfer(transfer.From, transfer.To, transfer.Value); err != nil {
		return err
	}

	reply := fmt.Sprintf("Confirmed <@%s> paid £%.02f to <@%s>", transfer.From, float32(transfer.Value)/100,
		transfer.To)
	_, _, err := slackClient.PostMessage(channel, slack.MsgOptionText(reply, false))
	return err
}

func handleSettle(channel string) error {
	ledger, err := potsService.Settle()

	if err != nil {
		return err
	}

	if len(ledger) == 0 {
		_, _, err := slackClient.PostMessage(channel, slack.MsgOptionText("Everybody is even!", false))
		return err
	}

	for _, debt := range ledger {
		msg := fmt.Sprintf("<@%s> should pay £%.02f to <@%s>", debt.From, float32(debt.Value)/100, debt.To)
		_, _, err := slackClient.PostMessage(channel, slack.MsgOptionText(msg, false))

		if err != nil {
			_, _, _ = slackClient.PostMessage(channel, slack.MsgOptionText("Error :(", false))
			return err
		}
	}

	return nil
}

func handleHelp(channel string) error {
	_, _, err := slackClient.PostMessage(channel, slack.MsgOptionText(
		`I understand the following commands:
@pokerbot in £10  (use this if you took £10 of chips without paying cash)
@pokerbot out £20 (use this if you put back £20 of chips without taking any cash)
@pokerbot pay @someone £5  (use this if you transferred £5 cash to someone)
@pokerbot settle  (use this to find out who owes what)
@pokerbot help`, false))
	return err
}

func main() {
	lambda.Start(HandleRequest)
}
