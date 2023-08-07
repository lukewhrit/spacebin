/*
* Copyright 2020-2023 Luke Whritenour

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
	"strings"
	"testing"
)

func TestGeneratePhrase(t *testing.T) {
	phrase := GeneratePhrase(2)

	phraseArray := strings.Split(phrase, "-")

	if len(phraseArray) != 2 {
		t.Error("didn't generate phrase of correct length")
	}
}

func TestGenerateKey(t *testing.T) {
	key := GenerateKey(8)

	if len(key) != 8 {
		t.Error("didn't generate key of correct length")
	}
}
