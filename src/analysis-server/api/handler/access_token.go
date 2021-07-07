package handler

import (
	"analysis-server/api/service"
	"common/log"
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AccessTokenHandler struct {
	loginCheckMu      sync.RWMutex
	tokenToNameMap    map[string]string
	expirationCheckMu sync.RWMutex
	tokenToTimeMap    map[string]int64
	quitCheckCh       chan bool
	authService       *service.AuthenService
	logger            *log.Logger
}

func NewAccessTokenHandler() *AccessTokenHandler {
	accTokenHandler := AccessTokenHandler{}
	accTokenHandler.tokenToNameMap = make(map[string]string)
	accTokenHandler.tokenToTimeMap = make(map[string]int64)
	accTokenHandler.quitCheckCh = make(chan bool, 1)
	return &accTokenHandler
}

func (at *AccessTokenHandler) InitAccessTokenHandler(authService *service.AuthenService, logger *log.Logger) {
	at.authService = authService
	at.logger = logger
}

func (at *AccessTokenHandler) insertToken(accessToken, userName string) {
	at.loginCheckMu.Lock()
	at.tokenToNameMap[accessToken] = userName
	at.loginCheckMu.Unlock()
	at.expirationCheckMu.Lock()
	at.tokenToTimeMap[accessToken] = time.Now().Unix() + (int64)(keepOnlineTime)
	at.expirationCheckMu.Unlock()
	return
}

func (at *AccessTokenHandler) delToken(accessToken string) {
	at.loginCheckMu.Lock()
	delete(at.tokenToNameMap, accessToken)
	at.loginCheckMu.Unlock()
	at.expirationCheckMu.Lock()
	delete(at.tokenToTimeMap, accessToken)
	at.expirationCheckMu.Unlock()
	return
}

func (at *AccessTokenHandler) delBatchToken(accessTokenSlice []string) {
	if len(accessTokenSlice) == 0 {
		return
	}
	var expiredUserName = []string{}
	at.loginCheckMu.Lock()
	for _, token := range accessTokenSlice {
		expiredUserName = append(expiredUserName, at.tokenToNameMap[token])
		delete(at.tokenToNameMap, token)
	}
	at.loginCheckMu.Unlock()
	at.expirationCheckMu.Lock()
	for _, v := range accessTokenSlice {
		delete(at.tokenToTimeMap, v)
	}
	at.expirationCheckMu.Unlock()
	//call logout service
	ctx := context.TODO()
	for _, name := range expiredUserName {
		ccErr := at.authService.Logout(ctx, name)
		if ccErr != nil {
			at.logger.WarnContext(ctx, "[delBatchToken] [AuthService.Logout,failed,errInfo: %s]", ccErr.Detail())
			return
		}
	}
	return
}

func (at *AccessTokenHandler) modifyTokenExpiredTime(accessToken string) {
	at.expirationCheckMu.Lock()
	at.tokenToTimeMap[accessToken] = time.Now().Unix() + (int64)(keepOnlineTime)
	at.expirationCheckMu.Unlock()
}

//检查用户是否登录过。
func (at *AccessTokenHandler) LoginCheck(action string, r *http.Request) (bool, error) {
	bIsPass := false
	cookie, err := r.Cookie("access_token")
	if err != nil {
		strErrMsg := err.Error()
		if action == "Login" && strings.Contains(strErrMsg, "http: lack the named cookie") {
			return true, nil
		}
		return bIsPass, err
	}
	accessToken := cookie.Value
	at.loginCheckMu.RLock()
	defer at.loginCheckMu.RUnlock()
	if _, ok := at.tokenToNameMap[accessToken]; ok {
		if action == "Login" {
			err = errors.New("The user has been to login,please logout the user.")
		} else {
			bIsPass = true
		}
	}
	return bIsPass, err
}

func (at *AccessTokenHandler) ExpirationCheck() {
	at.logger.LogDebug("expirationCheck,begin")
	for {
		select {
		case <-at.quitCheckCh:
			goto end
		case <-time.Tick(time.Second * 60):
			at.expirationCheckMu.RLock()
			var expirationToken = []string{}
			curTime := time.Now().Unix()
			for k, v := range at.tokenToTimeMap {
				if curTime > v {
					expirationToken = append(expirationToken, k)
				}
			}
			at.expirationCheckMu.RUnlock()
			at.delBatchToken(expirationToken)
		}
	}
end:
	at.logger.LogDebug("expirationCheck,end")
	return
}

func (at *AccessTokenHandler) QuitExpirationCheckService() {
	//at.quitCheckCh <- true
	close(at.quitCheckCh)
}
