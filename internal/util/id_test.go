package util_test

import (
	"strings"
	"testing"

	"github.com/orca-group/spirit/internal/util"
)

func TestGeneratePhrase(t *testing.T) {
	phrase := util.GeneratePhrase(2)

	phraseArray := strings.Split(phrase, "-")

	if len(phraseArray) != 2 {
		t.Error("didn't generate phrase of correct length")
	}
}

func TestGenerateKey(t *testing.T) {
	key := util.GenerateKey(8)

	if len(key) != 8 {
		t.Error("didn't generate key of correct length")
	}
}
