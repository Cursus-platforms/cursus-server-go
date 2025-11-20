package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// Khai báo bộ ký tự an toàn cho hậu tố Slug
const slugCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

// RandomString tạo ra một chuỗi ngẫu nhiên có độ dài xác định
func RandomString(length int) string {
	var result strings.Builder
	charsetLen := big.NewInt(int64(len(slugCharset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			fmt.Printf("Warning: Failed to generate crypto random number: %v. Using default '0'.\n", err)
			result.WriteByte('0')
			continue
		}
		result.WriteByte(slugCharset[randomIndex.Int64()])
	}

	return result.String()
}
