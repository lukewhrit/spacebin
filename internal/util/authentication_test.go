package util_test

import (
	"bytes"
	"testing"

	"github.com/lukewhrit/spacebin/internal/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashAndSalt(t *testing.T) {
	hash, err := util.HashAndSalt([]byte("testpassword"))
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(hash), []byte("testpassword")))
}

func TestHashAndSaltDifferentInputs(t *testing.T) {
	h1, _ := util.HashAndSalt([]byte("password1"))
	h2, _ := util.HashAndSalt([]byte("password2"))
	require.NotEqual(t, h1, h2)
}

func TestHashAndSaltTooLong(t *testing.T) {
	_, err := util.HashAndSalt(bytes.Repeat([]byte("a"), 73))
	require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
}

func TestParseTokenValid(t *testing.T) {
	tok, err := util.ParseToken("v1.public.secret.salt")
	require.NoError(t, err)
	require.Equal(t, "v1", tok.Version)
	require.Equal(t, "public", tok.Public)
	require.Equal(t, "secret", tok.Secret)
	require.Equal(t, "salt", tok.Salt)
}

func TestParseTokenThreeParts(t *testing.T) {
	tok, err := util.ParseToken("v1.public.secret")
	require.NoError(t, err)
	require.Equal(t, "v1", tok.Version)
	require.Equal(t, "public", tok.Public)
	require.Equal(t, "secret", tok.Secret)
	require.Empty(t, tok.Salt)
}

func TestParseTokenInvalid(t *testing.T) {
	_, err := util.ParseToken("invalid")
	require.Error(t, err)

	_, err = util.ParseToken("only.two")
	require.Error(t, err)
}

func TestMakeToken(t *testing.T) {
	tok := util.MakeToken(util.Token{
		Version: "v1",
		Public:  "pub",
		Secret:  "sec",
		Salt:    "salt",
	})
	require.Equal(t, "v1.pub.sec.salt", tok)
}

func TestPrngString(t *testing.T) {
	s := util.PrngString()
	require.NotEmpty(t, s)
	require.Len(t, s, 20) // 10 random bytes → 20 hex chars

	// Each call produces a different value
	s2 := util.PrngString()
	require.NotEqual(t, s, s2)
}

func TestGenerateStrings(t *testing.T) {
	a, b, c := util.GenerateStrings()
	require.NotEmpty(t, a)
	require.NotEmpty(t, b)
	require.NotEmpty(t, c)
	// All three values should be distinct
	require.NotEqual(t, a, b)
	require.NotEqual(t, b, c)
}

func TestParseAndMakeTokenRoundtrip(t *testing.T) {
	original := util.Token{
		Version: "v1",
		Public:  "pubkey",
		Secret:  "secretkey",
		Salt:    "saltvalue",
	}
	encoded := util.MakeToken(original)
	decoded, err := util.ParseToken(encoded)
	require.NoError(t, err)
	require.Equal(t, original, decoded)
}
