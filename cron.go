package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
)

type TimeType struct {
	Month  int
	Day    int
	Hour   int
	Minute int
}

type cronCacheType struct {
	Type   int8   `json:"type"`
	UserId int64  `json:"user_id"`
	Cron   string `json:"cron"`
	ChatId int64  `json:"chat_id"`
	Msg    string `json:"msg"`
}

var cronCache map[string]cronCacheType = make(map[string]cronCacheType)
var cronTask map[string]*cron.Cron = make(map[string]*cron.Cron)

func init() {
	file, err := os.ReadFile("./.bot_cache.json")
	if err == nil {
		json.Unmarshal(file, &cronCache)
		for i, j := range cronCache {
			c := cron.New()
			c.AddFunc(j.Cron, func() {
				SendMsg(j.ChatId, j.Msg, true)
				if j.Type == 0 {
					c.Stop()
				}
			})
			c.Start()
			cronTask[i] = c
		}
	}
}

func UpCache() {
	jsonStr, _ := json.Marshal(cronCache)
	WriteFile("./.bot_cache.json", jsonStr)
}

func addCron(t TimeType, chatId int64, msg string, uid int64, ty int8) string {
	id := GetRandomString(10)

	c := cron.New()
	var cs string
	if ty == 0 {
		cs = fmt.Sprintf("%d %d %d %d ?", t.Minute, t.Hour, t.Day, t.Month)
	} else {
		cs = fmt.Sprintf("%d %d %d * ?", t.Minute, t.Hour, t.Day)
	}
	c.AddFunc(cs, func() {
		SendMsg(chatId, msg, true)
		if ty == 0 {
			delete(cronCache, id)
			UpCache()
			c.Stop()
		}
	})

	cronCache[id] = cronCacheType{
		Type:   ty,
		Cron:   cs,
		ChatId: chatId,
		Msg:    msg,
		UserId: uid,
	}

	cronTask[id] = c

	UpCache()
	c.Start()

	return id
}
