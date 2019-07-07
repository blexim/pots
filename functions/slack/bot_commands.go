package main

import (
	"log"
	"regexp"
	"strconv"
)

type BuyinCommand struct {
	Player string
	Value int
}

type CashoutCommand struct {
	Player string
	Value int
}

type TransferCommand struct {
	From string
	To string
	Value int
}

type SettleCommand struct {
}

type HelpCommand struct {
}


func ParseBuyin(user string, text string) *BuyinCommand {
	re := regexp.MustCompile(`in £?(\d+)(\.\d+)?`)
	m := re.FindStringSubmatch(text)

	if m != nil {
		pounds, err := strconv.Atoi(m[1])

		if err != nil {
			log.Printf("Couldn't parse value in buyin command: %v", err);
			return nil
		}

		value := pounds * 100

		if len(m) > 1 && len(m[2]) > 0 {
			pence, err := strconv.Atoi(m[2][1:])

			if err != nil {
				log.Printf("Couldn't parse value in buyin command: %v", err);
				return nil
			}

			value += pence
		}

		return &BuyinCommand{user, value}
	}

	return nil
}

func ParseCashout(user string, text string) *CashoutCommand {
	re := regexp.MustCompile(`out £?(\d+)(\.\d+)?`)
	m := re.FindStringSubmatch(text)

	if m != nil {
		pounds, err := strconv.Atoi(m[1])

		if err != nil {
			log.Printf("Couldn't parse value in cashout command: %v", err);
			return nil
		}

		value := pounds * 100

		if len(m) > 1 && len(m[2]) > 0 {
			pence, err := strconv.Atoi(m[2][1:])

			if err != nil {
				log.Printf("Couldn't parse value in cashout command: %v", err);
				return nil
			}

			value += pence
		}

		return &CashoutCommand{user, value}
	}

	return nil
}

func ParseTransfer(user string, text string) *TransferCommand {
	re := regexp.MustCompile(`paid <(@.+)> £?(\d+)(\.\d+)?`)
	m := re.FindStringSubmatch(text)

	if m != nil {
		to := m[1][1:]

		pounds, err := strconv.Atoi(m[2])

		if err != nil {
			log.Printf("Couldn't parse value in cashout command: %v", err);
			return nil
		}

		value := pounds * 100

		if len(m) > 2 && len(m[3]) > 0 {
			pence, err := strconv.Atoi(m[3][1:])

			if err != nil {
				log.Printf("Couldn't parse value in cashout command: %v", err);
				return nil
			}

			value += pence
		}

		return &TransferCommand{user, to, value}
	}

	return nil
}

func ParseSettle(user string, text string) *SettleCommand {
	re := regexp.MustCompile(`settle`)
	if re.FindString(text) != "" {
		return &SettleCommand{}
	}

	return nil
}

func ParseHelp(user string, text string) *HelpCommand {
	re := regexp.MustCompile(`help`)
	if re.FindString(text) != "" {
		return &HelpCommand{}
	}

	return nil
}
