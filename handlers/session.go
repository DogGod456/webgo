package handlers

import "github.com/gorilla/sessions"

func Store(SecretKey string) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(SecretKey))
}
