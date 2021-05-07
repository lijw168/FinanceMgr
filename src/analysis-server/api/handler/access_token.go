package handler

import (
	"errors"
	"net/http"
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
	accessToken := cookie.String()
	delete(at.tokenToNameMap, accessToken)
	return nil
}

//检查用户是否登录过。
func (at *AccessTokenHandler) Checkout(action string, r *http.Request) (bool, error) {
	bIsPass := false
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return bIsPass, err
	}
	accessToken := cookie.String()
	at.mu.RLock()
	defer at.mu.RUnlock()
	if _, ok := at.tokenToNameMap[accessToken]; ok {
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
