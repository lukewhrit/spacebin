package util

import (
	"math/rand"
	"time"

	"github.com/lukewhrit/phrase"
)

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GeneratePhrase(length int) string {
	return phrase.Default.Generate(length).String()
}

func GenerateKey(length int) string {
	// Default key generation
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, length)

	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(b)
}

func GenerateID(t string, l int) string {
	if t == "phrase" {
		return GeneratePhrase(l)
	}

	return GenerateKey(l)
}
