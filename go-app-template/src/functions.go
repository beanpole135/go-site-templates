package main

import (
	"math/rand"
	"os"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomNumber(max int) int {
	if max < 1 {
		return 0
	}
	return (rand.Intn(max) + 1)
}

func randomString(length int) string {
	if length < 1 {
		return ""
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func randomID() string {
	return randomString(4) + time.Now().Format("123456")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
