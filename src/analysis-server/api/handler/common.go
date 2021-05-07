package handler

import (
	//"errors"
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
