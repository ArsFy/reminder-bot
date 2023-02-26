package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var BotToken string
var BotInfo tgbotapi.User
var Bot *tgbotapi.BotAPI

func init() {
	BotToken = os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		file, err := os.ReadFile("./token.conf")
		isToken := strings.Contains(string(file), ":")
		if err != nil || !isToken {
			fmt.Println("ConfigErr:", "You must set a Bot Token: https://github.com/ArsFy/reminder-bot/blob/main/README.md#bottoken")
			os.Exit(1)
		}
	} else if !strings.Contains(BotToken, ":") {
		fmt.Println("EnvErr:", "You must set a Bot Token: https://github.com/ArsFy/reminder-bot/blob/main/README.md#bottoken")
		os.Exit(1)
	}
}

func main() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		fmt.Println("BotErr:", err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := Bot.GetUpdatesChan(u)

	BotInfo, err := Bot.GetMe()
	if err != nil || BotInfo.UserName == "" {
		fmt.Println("BotErr:", err)
		os.Exit(1)
	}
	fmt.Println("Start Bot...")

	for update := range updates {
		if update.Message != nil {
			s, _ := json.Marshal(update.Message.Entities)
			fmt.Println(string(s))
			if update.Message.Entities[0].Type == "bot_command" {
				var commandClear []string
				command := strings.Split(update.Message.Text, " ")
				for _, j := range command {
					if j != "" {
						commandClear = append(commandClear, j)
					}
				}
				if len(commandClear) > 0 {
					switch strings.Replace(commandClear[0], "@"+BotInfo.UserName, "", 1) {
					case "/start":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join([]string{
							"透過 Telegram Bot 創建提醒事項\n",
							"/add\\_reminder `2-26(可省略)` `11:20` `提醒文本`",
							"/add\\_reminder\\_month `日 (數字)` `11:20` `提醒文本`",
							"\nGitHub: https://github\\.com/ArsFy/reminder\\-bot",
						}, "\n"))
						msg.ParseMode = "MarkdownV2"
						msg.ReplyToMessageID = update.Message.MessageID
						MsgErr(Bot.Send(msg))
					case "/add_reminder":
						if len(commandClear) >= 3 {
							var day, thisTime [2]int
							if strings.Contains(commandClear[1], "-") {
								dayList := strings.Split(commandClear[1], "-")
								if len(dayList) == 2 {
									day = [2]int{Int(dayList[0]), Int(dayList[1])}
									if strings.Contains(commandClear[2], ":") {
										timeList := strings.Split(commandClear[2], ":")
										if len(timeList) == 2 {
											thisTime = [2]int{Int(timeList[0]), Int(timeList[1])}
										} else {
											ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "時間格式錯誤: /add_reminder 2-26 11:20 提醒文本")
										}
									}
								} else {
									ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "日期格式錯誤: /add_reminder 2-26 11:20 提醒文本")
								}
							} else {
								timeList := strings.Split(commandClear[1], ":")
								if len(timeList) == 2 {
									thisTime = [2]int{Int(timeList[0]), Int(timeList[1])}
									now := time.Now()
									if thisTime[0] < now.Hour() || (thisTime[0] == now.Hour() && thisTime[1] < now.Minute()) {
										now = time.Now().Add(time.Hour * 24)
										dayList := strings.Split(now.Format("01-02"), "-")
										day = [2]int{Int(dayList[0]), Int(dayList[1])}
									}
								} else {
									ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "時間格式錯誤: /add_reminder 2-26 11:20 提醒文本")
								}
							}
							if CheckTime(day, thisTime) {
								fmt.Println(day, thisTime)
								// addCron
							} else {
								ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "時間超出範圍: /add_reminder 2023-2-26 11:20 提醒文本")
							}
						} else {
							ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "缺少參數: /add_reminder 2023-2-26 11:20 提醒文本")
						}
					}
				}
			}
		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "pass":
				// callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Pass")
				// if _, err := Bot.Request(callback); err != nil {
				// 	panic(err)
				// }
				// delmsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				// if _, err := Bot.Request(delmsg); err != nil {
				// 	panic(err)
				// }
			}
		}
	}
}
