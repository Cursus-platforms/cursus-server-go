package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
)

// GenerateRandomOTP tạo ra một chuỗi OTP ngẫu nhiên gồm n chữ số.
func GenerateRandomOTP(length int) string {
	// Dùng crypto/rand để đảm bảo tính an toàn (secure random)
	if length <= 0 || length > 10 {
		return ""
	}

	b := make([]byte, (length+1)/2) // Tạo byte slice đủ lớn
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	// Chuyển đổi byte ngẫu nhiên thành số và cắt theo độ dài mong muốn
	// Ví dụ: OTP 6 số.
	otpValue := strconv.Itoa(int(b[0]))

	for i := 1; i < len(b); i++ {
		otpValue += strconv.Itoa(int(b[i]))
	}

	if len(otpValue) > length {
		return otpValue[:length]
	}
	// Nếu không đủ, có thể thêm số 0 vào đầu (nếu cần thiết cho display)
	return otpValue
}
