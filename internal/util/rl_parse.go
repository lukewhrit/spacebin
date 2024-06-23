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

package util

import (
	"strconv"
	"strings"
	"time"
)

func ParseRatelimiterString(rl string) (int, time.Duration, error) {
	array := strings.Split(rl, "x")

	if len(array) != 2 {
		return 0, 0, ErrTooManyParts
	}

	intArray := make([]int, 0)

	for i := range array {
		newInt, err := strconv.Atoi(array[i])

		if err != nil {
			return 0, 0, err
		}

		intArray = append(intArray, newInt)
	}

	return intArray[0], time.Duration(intArray[1]) * time.Second, nil
}
