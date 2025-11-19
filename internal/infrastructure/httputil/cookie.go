package httputil

import (
	"net/http"
	"time"
)

func SetRefreshTokenCookie(w http.ResponseWriter, token string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   token,
		Path:    "/api/v1/auth",
		Expires: time.Now().Add(duration),

		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
