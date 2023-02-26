package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
)

type cronCacheType struct {
	Cron   string `json:"cron"`
	ChatId int64  `json:"chat_id"`
	Msg    string `json:"msg"`
}

var cronCache map[string]cronCacheType = make(map[string]cronCacheType)

func init() {
	file, err := os.ReadFile("./.bot_cache")
	if err == nil {
		json.Unmarshal(file, &cronCache)
		for _, j := range cronCache {
			c := cron.New()
			c.AddFunc(j.Cron, func() {
				SendMsg(j.ChatId, j.Msg)
				c.Stop()
			})
			c.Start()
		}
	}
}

func addCron(month, day, hour, minute int, chatId int64, msg string) {
	id := GetRandomString(10)

	c := cron.New()
	cs := fmt.Sprintf("0 %d %d %d %d ?", minute, hour, day, month)
	c.AddFunc(cs, func() {
		SendMsg(chatId, msg)
		c.Stop()
	})

	cronCache[id] = cronCacheType{
		Cron:   cs,
		ChatId: chatId,
		Msg:    msg,
	}

	jsonStr, _ := json.Marshal(cronCache)
	WriteFile("./.bot_cache", jsonStr)
	c.Start()
}
