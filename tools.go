package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func Int(num string) int {
	i, _ := strconv.Atoi(num)
	return i
}

func CheckTime(d, t [2]int) bool {
	now := time.Now()
	return (d[0] > int(now.Month()) || (d[0] == int(now.Month()) && d[1] >= now.Day())) && (d[0] < 13 && d[1] < 32 && t[0] < 13 && t[1] < 61)
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
