package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Int(num string) int {
	i, _ := strconv.Atoi(num)
	return i
}

func CheckTime(t TimeType) bool {
	now := time.Now()
	return (t.Month > int(now.Month()) || (t.Month == int(now.Month()) && t.Day > now.Day()) || (t.Month == int(now.Month()) && t.Day == now.Day() && !(t.Hour < now.Hour() || (t.Hour == now.Hour() && t.Minute < now.Minute())))) && (t.Month < 13 && t.Day < 32 && t.Hour < 25 && t.Minute < 61)
}

func WriteFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func AddZero(t int) string {
	if t < 10 {
		return fmt.Sprintf("0%d", t)
	} else {
		return fmt.Sprint(t)
	}
}

func AddZeroStr(t string) string {
	if len(t) == 1 && t != "*" {
		return "0" + t
	} else if t != "*" {
		return t
	} else {
		return "每月"
	}
}

func GetUser(user *tgbotapi.User) string {
	username := "@" + user.UserName
	if username == "@" {
		username = fmt.Sprintf(`[%s %s](tg://user?id=%d)`, user.LastName, user.FirstName, user.ID)
	}
	return username
}
