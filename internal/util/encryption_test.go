/*
 * Copyright 2020-2026 Luke Whritenour

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/lukewhrit/spacebin/internal/util"
	"github.com/stretchr/testify/require"
)

// browserEnvelope was made by web/static/crypto.js, encrypting browserPlaintext with browserPassword
const (
	browserEnvelope  = `{"v":1,"iter":600000,"salt":"f7taG8g1oRXn6EiOzsYviA==","iv":"UUdU8C/6+iG2lV9q","ct":"eo4MFscrTwp7NLK8h+aQPLrFDHdCL54aw0Yx+xjSaeIswl6hVIJj+R4W2w9GI0y19O48aIrjBq9Ib9o="}`
	browserPassword  = "correct horse"
	browserPlaintext = "package main\n\nfunc main() {} // héllo 🚀"
)

func TestValidateEnvelopeFromBrowser(t *testing.T) {
	require.NoError(t, util.ValidateEnvelope(browserEnvelope))
}

// TestBrowserEnvelopeFormat decrypts the browser's envelope independently, to pin down its format.
func TestBrowserEnvelopeFormat(t *testing.T) {
	var envelope util.Envelope
	require.NoError(t, json.Unmarshal([]byte(browserEnvelope), &envelope))

	salt, _ := base64.StdEncoding.DecodeString(envelope.Salt)
	iv, _ := base64.StdEncoding.DecodeString(envelope.IV)
	ciphertext, _ := base64.StdEncoding.DecodeString(envelope.Ciphertext)

	key, err := pbkdf2.Key(sha256.New, browserPassword, salt, envelope.Iterations, 32)
	require.NoError(t, err)

	block, err := aes.NewCipher(key)
	require.NoError(t, err)
	gcm, err := cipher.NewGCM(block)
	require.NoError(t, err)

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	require.NoError(t, err)
	require.Equal(t, browserPlaintext, string(plaintext))
}

// envelopeWith returns the browser's envelope with field replaced by value.
func envelopeWith(t *testing.T, field string, value any) string {
	t.Helper()

	fields := map[string]any{}
	require.NoError(t, json.Unmarshal([]byte(browserEnvelope), &fields))
	fields[field] = value

	encoded, err := json.Marshal(fields)
	require.NoError(t, err)

	return string(encoded)
}

func TestValidateEnvelopeInvalid(t *testing.T) {
	tests := map[string]struct {
		content string
		message string
	}{
		"plain text":          {"hello world", "not a valid envelope"},
		"unknown field":       {envelopeWith(t, "password", "oops"), "not a valid envelope"},
		"wrong type":          {envelopeWith(t, "iter", "600000"), "not a valid envelope"},
		"trailing data":       {browserEnvelope + "{}", "data after the envelope"},
		"version":             {envelopeWith(t, "v", 2), "unsupported envelope version 2"},
		"too few iterations":  {envelopeWith(t, "iter", util.MinEnvelopeIterations-1), "iterations must be between"},
		"too many iterations": {envelopeWith(t, "iter", util.MaxEnvelopeIterations+1), "iterations must be between"},
		"salt not base64":     {envelopeWith(t, "salt", "!!!"), "salt is not valid base64"},
		"short salt":          {envelopeWith(t, "salt", base64.StdEncoding.EncodeToString(make([]byte, 15))), "salt has the wrong length"},
		"long iv":             {envelopeWith(t, "iv", base64.StdEncoding.EncodeToString(make([]byte, 13))), "iv has the wrong length"},
		"ciphertext no tag":   {envelopeWith(t, "ct", base64.StdEncoding.EncodeToString(make([]byte, 15))), "ct has the wrong length"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.ErrorContains(t, util.ValidateEnvelope(test.content), test.message)
		})
	}
}

func TestValidateEnvelopeTrailingWhitespace(t *testing.T) {
	require.NoError(t, util.ValidateEnvelope(browserEnvelope+"\n"))
}
