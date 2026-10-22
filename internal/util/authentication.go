package util

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd []byte) (string, error) {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	// It fails if pwd is longer than 72 bytes.
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), err
}

func PrngString() string {
	b := make([]byte, 10)
	rand.Read(b) // Never returns an error, see crypto/rand docs

	return fmt.Sprintf("%x", b)
}

func GenerateStrings() (a, b, c string) {
	return PrngString(), PrngString(), PrngString()
}

func ParseToken(token string) (Token, error) {
	var tok Token
	toks := strings.Split(token, ".")

	if len(toks) < 3 {
		return Token{}, errors.New("invalid token")
	}

	tok.Version = toks[0]
	tok.Public = toks[1]
	tok.Secret = toks[2]

	if len(toks) == 4 {
		tok.Salt = toks[3]
	}

	return tok, nil
}

func MakeToken(token Token) string {
	return fmt.Sprintf("%s.%s.%s.%s", token.Version, token.Public, token.Secret, token.Salt)
}
