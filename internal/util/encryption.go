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

package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Bounds on the PBKDF2 iterations of an envelope. Too few makes passwords easy to brute force, too many
// would freeze the browser of anyone opening the document.
const (
	MinEnvelopeIterations = 100_000
	MaxEnvelopeIterations = 2_000_000
)

// Envelope is the content of an encrypted document. Browsers derive an AES-256-GCM key from the password
// with PBKDF2-SHA256, so the server never sees the password or the plaintext. Binary fields are standard
// base64.
type Envelope struct {
	Version    int    `json:"v"`
	Iterations int    `json:"iter"`
	Salt       string `json:"salt"` // 16 bytes
	IV         string `json:"iv"`   // 12 bytes
	Ciphertext string `json:"ct"`   // Includes the 16 byte GCM tag
}

// ValidateEnvelope checks that content is a well-formed Envelope. It can't check that the ciphertext
// decrypts, as only the author knows the password.
func ValidateEnvelope(content string) error {
	var envelope Envelope

	decoder := json.NewDecoder(bytes.NewReader([]byte(content)))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&envelope); err != nil {
		return fmt.Errorf("encrypted content is not a valid envelope: %w", err)
	}

	if _, err := decoder.Token(); err != io.EOF {
		return errors.New("encrypted content has data after the envelope")
	}

	if envelope.Version != 1 {
		return fmt.Errorf("unsupported envelope version %d", envelope.Version)
	}

	if envelope.Iterations < MinEnvelopeIterations || envelope.Iterations > MaxEnvelopeIterations {
		return fmt.Errorf("envelope iterations must be between %d and %d", MinEnvelopeIterations, MaxEnvelopeIterations)
	}

	fields := []struct {
		name    string
		value   string
		minSize int
		maxSize int
	}{
		{"salt", envelope.Salt, 16, 16},
		{"iv", envelope.IV, 12, 12},
		{"ct", envelope.Ciphertext, 16, -1},
	}

	for _, field := range fields {
		decoded, err := base64.StdEncoding.DecodeString(field.value)

		if err != nil {
			return fmt.Errorf("envelope %s is not valid base64", field.name)
		}

		if len(decoded) < field.minSize || (field.maxSize != -1 && len(decoded) > field.maxSize) {
			return fmt.Errorf("envelope %s has the wrong length", field.name)
		}
	}

	return nil
}
