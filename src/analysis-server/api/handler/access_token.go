package handler

import (
	"errors"
	"net/http"
	"strings"
	"sync"
)

type AccessTokenHandler struct {
	mu             sync.RWMutex
	tokenToNameMap map[string]string
}

func NewAccessTokenHandler() *AccessTokenHandler {
	accTokenHandler := AccessTokenHandler{}
	accTokenHandler.tokenToNameMap = make(map[string]string)
	return &accTokenHandler
}

func (at *AccessTokenHandler) insertElem(accessToken, userName string) {
	at.mu.Lock()
	defer at.mu.Unlock()
	at.tokenToNameMap[accessToken] = userName
	return
}

func (at *AccessTokenHandler) delElem(r *http.Request) error {
	at.mu.Lock()
	defer at.mu.Unlock()
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return err
	}
	accessToken := cookie.Value
	delete(at.tokenToNameMap, accessToken)
	return nil
}

//检查用户是否登录过。
func (at *AccessTokenHandler) Checkout(action string, r *http.Request) (bool, error) {
	bIsPass := false
	cookie, err := r.Cookie("access_token")
	if err != nil {
		strErrMsg := err.Error()
		if action == "Login" && strings.Contains(strErrMsg, "http: named cookie not present") {
			return true, nil
		}
		return bIsPass, err
	}
	accessToken := cookie.Value
	at.mu.RLock()
	defer at.mu.RUnlock()
	if _, ok := at.tokenToNameMap[accessToken]; ok {
		if action == "Login" {
			err = errors.New("the user has been to login,please logout the user.")
		} else {
			bIsPass = true
		}
	}
	return bIsPass, err
}
