package slug

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// nonAlphanumericRegex tìm các ký tự không phải chữ cái, số, hoặc dấu gạch ngang
var nonAlphanumericRegex = regexp.MustCompile(`[^\p{L}0-9]+`)

func removeDiacritics(s string) string {
	// Logic chuẩn hóa Unicode và loại bỏ dấu
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)

	s, _, _ = transform.String(t, s)
	return s
}

// Generate tạo ra một chuỗi slug an toàn cho URL từ một chuỗi đầu vào.
func Generate(title string) string {
	s := removeDiacritics(title)
	s = strings.ToLower(s)
	s = nonAlphanumericRegex.ReplaceAllString(s, "-")
	s = strings.ReplaceAll(s, "--", "-")
	s = strings.Trim(s, "-")

	return s
}
