package validator

import (
	"regexp"
	"strings"
)

const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`

var (
	emailCompiledRegex = regexp.MustCompile(emailRegex)
)

func IsRequired(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

func IsValidPrice(price float64) bool {
	return price > 0
}

func IsValidEmail(email string) bool {
	if !IsRequired(email) {
		return false
	}
	return emailCompiledRegex.MatchString(email)
}

func IsValidStatus(status string, allowedStatuses []string) bool {
	if !IsRequired(status) {
		return false
	}
	for _, allowed := range allowedStatuses {
		if status == allowed {
			return true
		}
	}
	return false
}
