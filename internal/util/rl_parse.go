package util

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func ParseRatelimiterString(rl string) (int, time.Duration, error) {
	array := strings.Split(rl, "x")

	if len(array) != 2 {
		return 0, 0, errors.New("ratelimiter string invalid: too many parts")
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
