package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MsgErr(msg tgbotapi.Message, err error) tgbotapi.Message {
	if err != nil {
		fmt.Println("SendErr:", err)
	}
	return msg
}
