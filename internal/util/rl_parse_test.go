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
	"strconv"
	"testing"
	"time"

	"github.com/lukewhrit/spacebin/internal/util"
	"github.com/stretchr/testify/require"
)

func TestParseRatelimiterTooManyParts(t *testing.T) {
	rlString := "200x5x10"
	_, _, err := util.ParseRatelimiterString(rlString)
	require.Error(t, err, util.ErrTooManyParts)
}

func TestParseRatelimiterInvalidSyntax(t *testing.T) {
	rlString := "1,000x2.0"
	_, _, err := util.ParseRatelimiterString(rlString)
	require.Error(t, err, strconv.ErrSyntax)
}

func TestParseRatelimiterOutOfRange(t *testing.T) {
	rlString := "9223372036854775808x5"
	_, _, err := util.ParseRatelimiterString(rlString)
	require.Error(t, err, strconv.ErrRange)
}

func TestParseRatelimiterString(t *testing.T) {
	rlString := "200x5"

	reqs, secs, err := util.ParseRatelimiterString(rlString)

	require.NoError(t, err, nil)
	require.Equal(t, reqs, 200)
	require.Equal(t, secs, 5*time.Second)
}
