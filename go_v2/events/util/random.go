package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Randomly generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEventName() string {
	return RandomString(9)
}

func RandomDate() time.Time {
	min := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC).Unix()

	sec := RandomInt(min, max)
	return time.Unix(sec, 0)
}

func RandomLocation() string {
	return fmt.Sprintf("%s %s", RandomString(6), RandomString(9))
}

func RandomClass() string {
	class := []string{"Elite", "Mid", "Normal"}
	return class[rand.Intn(len(class))]
}
