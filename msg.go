package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func ReplyMsg(chatId int64, replyId int, msg string) {
	sendmsg := tgbotapi.NewMessage(chatId, msg)
	sendmsg.ReplyToMessageID = replyId
	MsgErr(Bot.Send(sendmsg))
}

func SendMsg(chatId int64, msg string) {
	sendmsg := tgbotapi.NewMessage(chatId, msg)
	MsgErr(Bot.Send(sendmsg))
}
