package handler

import (
	"fmt"
	"math/rand"
	"time"
)

var seeded bool = false

// CreateTrackingId generates a tracing ID
func CreateTrackingId(salt string) string {
	if !seeded {
		seeded = true
		rand.Seed(time.Now().UnixNano())
	}
	return fmt.Sprintf("%c%c-%d%s-%d%s", getRandomLetterCode(),
		getRandomLetterCode(), len(salt)/2, getRandomNumber(3),
		len(salt)/2, getRandomNumber(7))
}

func getRandomLetterCode() uint32 {
	return 65 + uint32(rand.Intn(25))
}

func getRandomNumber(digits int) string {
	str := ""
	for i := 0; i < digits; i++ {
		str = fmt.Sprintf("%s%d", str, rand.Intn(10))
	}
	return str
}
