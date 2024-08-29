package handler

import "net/http"

func getToken(r *http.Request) string {
	if v := r.Header.Get("Authorization"); v != "" {
		return v
	}

	cookie, err := r.Cookie("Authorization")
	if err != nil {
		return ""
	}

	return cookie.Value
}
