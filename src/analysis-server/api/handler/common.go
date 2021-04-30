package handler

import (
	"errors"
	"net/http"
)

func Auth(secretIds []string) bool {
	auth := func(secretId string) bool {
		return true
	}
	for _, secretId := range secretIds {
		if auth(secretId) {
			return true
		}
	}
	return false
}

func isAdmin(req *http.Request) bool {
	return Auth(req.Header["Secret-Id"])
}

//检查用户是否登录过。
func Checkout(action string, r *http.Request) (bool, error) {
	bIsPass := false
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return bIsPass, err
	}
	accessToken := cookie.String()
	if _, ok := tokenToNameMap[accessToken]; ok {
		if action == "Login" {
			bIsPass = false
			err = errors.New("the user has been to login,please logout the user.")
		} else {
			bIsPass = true
		}
	} else {
		if action == "Login" {
			bIsPass = true
		} else {
			bIsPass = false
		}
	}
	return bIsPass, err
}
