package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	GenerateGameID(8)
}

const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateGameID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}
