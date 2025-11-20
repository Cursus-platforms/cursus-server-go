package utils

import (
	"crypto/rand"
	"io"
	"strconv"
)

func GenerateRandomOTP(length int) string {
	if length <= 0 || length > 10 {
		return ""
	}

	b := make([]byte, (length+1)/2)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	otpValue := strconv.Itoa(int(b[0]))

	for i := 1; i < len(b); i++ {
		otpValue += strconv.Itoa(int(b[i]))
	}

	if len(otpValue) > length {
		return otpValue[:length]
	}
	return otpValue
}
