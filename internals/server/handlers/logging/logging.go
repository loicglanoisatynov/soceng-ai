package logging

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	db_cookies "soceng-ai/database/tables/cookies"
	db_users "soceng-ai/database/tables/users"
)

func randString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func IssueCookie(identifier string) string {
	cookie := randString(32)
	id := db_users.Get_user_id_by_username_or_email(identifier)
	err := db_cookies.Register_cookie(id, cookie)
	if err != nil {
		fmt.Println("Error registering cookie in database: ", err)
		return ""
	}
	fmt.Println("Cookie registered in database")
	return cookie
}

func IsCookieValid(cookie string) bool {
	id := db_cookies.Get_user_id_by_cookie(cookie)
	if id == -1 {
		return false
	}
	return true
}
