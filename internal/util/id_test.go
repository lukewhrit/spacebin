/*
* Copyright 2020-2024 Luke Whritenour

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
	"strings"
	"testing"

	"github.com/orca-group/spirit/internal/util"
	"github.com/stretchr/testify/require"
)

func TestGeneratePhrase(t *testing.T) {
	phrase := util.GeneratePhrase(2)
	phraseArray := strings.Split(phrase, "-")

	require.Len(t, phraseArray, 2)
}

func TestGenerateKey(t *testing.T) {
	key := util.GenerateKey(8)
	require.Len(t, key, 8)

}

func TestGenerateID(t *testing.T) {
	phrase := util.GenerateID("phrase", 2)
	phraseArray := strings.Split(phrase, "-")

	require.Len(t, phraseArray, 2)

	key := util.GenerateID("key", 8)
	require.Len(t, key, 8)
}
