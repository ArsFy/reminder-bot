package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var BotToken string
var timeZoneSet int = 100
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

	timeZone := os.Getenv("TIMEZONE")
	if timeZone != "" {
		timeZoneNum, err := strconv.Atoi(timeZone)
		if err == nil && timeZoneNum > -13 && timeZoneNum < 13 {
			time.FixedZone("UTC", timeZoneNum*3600)
			timeZoneSet = timeZoneNum
		}
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
			if len(update.Message.Entities) > 0 && update.Message.Entities[0].Type == "bot_command" {
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
						msgArr := []string{
							"透過 Telegram Bot 創建提醒事項\n",
							"列出所有提醒: /my\\_reminders",
							"新增提醒: /add\\_reminder `2-26(可省略)` `11:20` `提醒文本`",
							"新增每月提醒: /add\\_reminder\\_month `日(數字)` `11:20` `提醒文本`",
							"\nGitHub: https://github\\.com/ArsFy/reminder\\-bot",
						}
						if timeZoneSet != 100 {
							msgArr = append(msgArr, []string{"", fmt.Sprintf("這個 Bot 現在運行在 UTC%s%d", Operator3(timeZoneSet <= 0, "-", "\\+"), timeZoneSet)}...)
						}
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(msgArr, "\n"))
						msg.ParseMode = "MarkdownV2"
						msg.ReplyToMessageID = update.Message.MessageID
						MsgErr(Bot.Send(msg))
					case "/my_reminders":
						var rList [][]tgbotapi.InlineKeyboardButton = [][]tgbotapi.InlineKeyboardButton{}
						for i, j := range cronCache {
							if j.UserId == update.Message.From.ID {
								tL := strings.Split(j.Cron, " ")
								rList = append(rList, tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("刪除 %s-%s %s:%s %s", AddZeroStr(tL[3]), AddZeroStr(tL[2]), AddZeroStr(tL[1]), AddZeroStr(tL[0]), strings.Split(j.Msg, "提醒內容: ")[1]), fmt.Sprintf("remove_%s_%d", i, update.Message.From.ID)),
								))
							}
						}
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("↓ 你有 %d 個提醒事件", len(rList)))
						msg.ReplyToMessageID = update.Message.MessageID
						if len(rList) > 0 {
							msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rList...)
						}
						Bot.Send(msg)
					case "/add_reminder":
						if len(commandClear) > 2 {
							testText := strings.Join(commandClear[1:], ",")
							var tstru TimeType
							var content string
							if Dtime.MatchString(testText) {
								tstri := Dtime.FindStringSubmatch(testText)
								tstru = TimeType{
									Month:  Int(tstri[1]),
									Day:    Int(tstri[2]),
									Hour:   Int(tstri[3]),
									Minute: Int(tstri[4]),
								}
								content = tstri[5]
							} else if Htime.MatchString(testText) {
								now := time.Now()
								tstri := Htime.FindStringSubmatch(testText)
								if Int(tstri[1]) < now.Hour() || (Int(tstri[1]) == now.Hour() && Int(tstri[2]) < now.Minute()) {
									now = time.Now().Add(time.Hour * 24)
								}
								tstru = TimeType{
									Month:  int(now.Month()),
									Day:    now.Day(),
									Hour:   Int(tstri[1]),
									Minute: Int(tstri[2]),
								}
								content = tstri[3]
							} else {
								ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "時間格式錯誤: /add_reminder 2023-2-26 11:20 提醒文本")
								return
							}
							if CheckTime(tstru) {
								username := GetUser(update.Message.From)
								id := addCron(tstru, update.Message.Chat.ID, fmt.Sprintf(
									"%s 事件提醒\n設定時間: %s\n觸發時間: %d\\-%s\\-%s %s:%s\n提醒內容: %s",
									username,
									time.Now().Format("2006\\-01\\-02 15:04"),
									time.Now().Year(),
									AddZero(tstru.Month), AddZero(tstru.Day), AddZero(tstru.Hour), AddZero(tstru.Minute),
									content,
								), update.Message.From.ID, 0)
								sendmsg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
									"設定新提醒\n觸發時間: %d-%s-%s %s:%s\n提醒內容: %s",
									time.Now().Year(),
									AddZero(tstru.Month), AddZero(tstru.Day), AddZero(tstru.Hour), AddZero(tstru.Minute),
									content,
								))
								sendmsg.ReplyToMessageID = update.Message.MessageID
								sendmsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("× 取消", fmt.Sprintf("remove_%s_%d", id, update.Message.From.ID)),
									),
								)
								MsgErr(Bot.Send(sendmsg))
							} else {
								ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "錯誤: 時間超出範圍/已過期")
							}
						} else {
							ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "缺少參數: /add_reminder 2023-2-26 11:20 提醒文本")
						}
					case "/add_reminder_month":
						if len(commandClear) == 4 {
							testText := strings.Join(commandClear[1:], ",")
							if Mtime.MatchString(testText) {
								tstri := Mtime.FindStringSubmatch(testText)
								tstru := TimeType{
									Month:  12,
									Day:    Int(tstri[1]),
									Hour:   Int(tstri[2]),
									Minute: Int(tstri[3]),
								}

								if CheckTime(tstru) {
									username := GetUser(update.Message.From)
									id := addCron(tstru, update.Message.Chat.ID, fmt.Sprintf(
										"%s 事件提醒\n設定時間: %s\n觸發時間: 每月 %s 日 %s:%s\n提醒內容: %s",
										username,
										time.Now().Format("2006\\-01\\-02 15:04"),
										AddZero(tstru.Day), AddZero(tstru.Hour), AddZero(tstru.Minute),
										tstri[4],
									), update.Message.From.ID, 1)
									sendmsg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
										"設定新每月提醒\n觸發時間: 每月 %s 日 %s:%s\n提醒內容: %s",
										AddZero(tstru.Day), AddZero(tstru.Hour), AddZero(tstru.Minute),
										tstri[4],
									))
									sendmsg.ReplyToMessageID = update.Message.MessageID
									sendmsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
										tgbotapi.NewInlineKeyboardRow(
											tgbotapi.NewInlineKeyboardButtonData("× 取消", fmt.Sprintf("remove_%s_%d", id, update.Message.From.ID)),
										),
									)
									MsgErr(Bot.Send(sendmsg))
								} else {
									ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "錯誤: 時間超出範圍")
								}
							} else {
								ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "時間格式錯誤: /add_reminder_month 26 11:20 提醒文本")
							}
						} else {
							ReplyMsg(update.Message.Chat.ID, update.Message.MessageID, "缺少參數: /add_reminder_month 26 11:20 提醒文本")
						}
					}
				}
			}
		} else if update.CallbackQuery != nil {
			commList := strings.Split(update.CallbackQuery.Data, "_")
			switch commList[0] {
			case "remove":
				if len(commList) == 3 && commList[2] == fmt.Sprint(update.CallbackQuery.From.ID) {
					if cronTask[commList[1]] != nil {
						cronTask[commList[1]].Stop()
					}
					delete(cronTask, commList[1])
					delete(cronCache, commList[1])
					UpCache()
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "已取消")
					if _, err := Bot.Request(callback); err != nil {
						panic(err)
					}
					if !strings.Contains(update.CallbackQuery.Message.Text, "↓") {
						editmsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "提醒已取消")
						if _, err := Bot.Request(editmsg); err != nil {
							panic(err)
						}
					} else {
						var rList [][]tgbotapi.InlineKeyboardButton = [][]tgbotapi.InlineKeyboardButton{}
						for i, j := range cronCache {
							if j.UserId == update.Message.From.ID {
								tL := strings.Split(j.Cron, " ")
								rList = append(rList, tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("刪除 %s-%s %s:%s %s", AddZeroStr(tL[3]), AddZeroStr(tL[2]), AddZeroStr(tL[1]), AddZeroStr(tL[0]), strings.Split(j.Msg, "提醒內容: ")[1]), fmt.Sprintf("remove_%s_%d", i, update.Message.From.ID)),
								))
							}
						}
						editmsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("↓ 你有 %d 個提醒事件", len(rList)))
						if _, err := Bot.Request(editmsg); err != nil {
							panic(err)
						}
						if len(rList) > 0 {
							editmsg2 := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, tgbotapi.NewInlineKeyboardMarkup(rList...))
							if _, err := Bot.Request(editmsg2); err != nil {
								panic(err)
							}
						}
					}
					go func(cq *tgbotapi.CallbackQuery) {
						time.Sleep(time.Minute)
						delmsg := tgbotapi.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
						if _, err := Bot.Request(delmsg); err != nil {
							panic(err)
						}
					}(update.CallbackQuery)
				} else {
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "只允許設定自己的提醒事件")
					if _, err := Bot.Request(callback); err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
