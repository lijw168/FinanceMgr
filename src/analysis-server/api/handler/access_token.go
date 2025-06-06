package handler

import (
	"context"
	"errors"
	"financeMgr/src/analysis-server/api/service"
	"financeMgr/src/common/log"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type onlineUserInfo struct {
	iOperatorID int
	iRole       int
}

type AccessTokenHandler struct {
	loginCheckMu      sync.RWMutex
	tokenToOptIDMap   map[string]*onlineUserInfo
	expirationCheckMu sync.RWMutex
	tokenToTimeMap    map[string]int64
	quitCheckCh       chan bool
	authService       *service.AuthenService
	optInfoService    *service.OperatorInfoService
	logger            *log.Logger
	//for statistics;test
	checkStatusCount int
}

func NewAccessTokenHandler() *AccessTokenHandler {
	accTokenHandler := AccessTokenHandler{}
	accTokenHandler.tokenToOptIDMap = make(map[string]*onlineUserInfo)
	accTokenHandler.tokenToTimeMap = make(map[string]int64)
	accTokenHandler.quitCheckCh = make(chan bool, 1)
	accTokenHandler.checkStatusCount = 0
	return &accTokenHandler
}

func (at *AccessTokenHandler) InitAccessTokenHandler(authService *service.AuthenService,
	optInfoService *service.OperatorInfoService, logger *log.Logger) {
	at.authService = authService
	at.optInfoService = optInfoService
	at.logger = logger
}

func (at *AccessTokenHandler) insertToken(accessToken string, iOptID int) {
	at.loginCheckMu.Lock()
	userInfo := onlineUserInfo{iOperatorID: iOptID, iRole: 0}
	at.tokenToOptIDMap[accessToken] = &userInfo
	at.loginCheckMu.Unlock()
	at.expirationCheckMu.Lock()
	at.tokenToTimeMap[accessToken] = time.Now().Unix() + (int64)(keepOnlineTime)
	at.expirationCheckMu.Unlock()
}

func (at *AccessTokenHandler) delToken(accessToken string) {
	at.loginCheckMu.Lock()
	delete(at.tokenToOptIDMap, accessToken)
	at.loginCheckMu.Unlock()
	at.expirationCheckMu.Lock()
	delete(at.tokenToTimeMap, accessToken)
	at.expirationCheckMu.Unlock()
}

func (at *AccessTokenHandler) delBatchToken(accessTokenSlice []string) {
	if len(accessTokenSlice) == 0 {
		return
	}
	var expiredOptID = []int{}
	at.loginCheckMu.Lock()
	for _, token := range accessTokenSlice {
		if usrInfo, ok := at.tokenToOptIDMap[token]; ok {
			expiredOptID = append(expiredOptID, usrInfo.iOperatorID)
		} else {
			errMsg := fmt.Sprintf("can't find the %s token in tokenToOptIDMap", token)
			at.logger.Error(errMsg)
			panic(errMsg)
		}
		delete(at.tokenToOptIDMap, token)
	}
	at.loginCheckMu.Unlock()
	at.expirationCheckMu.Lock()
	for _, v := range accessTokenSlice {
		delete(at.tokenToTimeMap, v)
	}
	at.expirationCheckMu.Unlock()
	//call logout service
	ctx := context.TODO()
	for _, optId := range expiredOptID {
		ccErr := at.authService.Logout(ctx, optId)
		if ccErr != nil {
			at.logger.WarnContext(ctx, "[delBatchToken] [AuthService.Logout,failed,errInfo: %s]", ccErr.Detail())
			return
		}
		at.logger.DebugContext(ctx, "[delBatchToken] [the operator %d has been logout]", optId)
	}
}

func (at *AccessTokenHandler) modifyTokenExpiredTime(accessToken string) {
	at.expirationCheckMu.Lock()
	at.tokenToTimeMap[accessToken] = time.Now().Unix() + (int64)(keepOnlineTime)
	at.expirationCheckMu.Unlock()
	at.checkStatusCount = at.checkStatusCount + 1
}

// 用户在线检查
func (at *AccessTokenHandler) LoginCheck(action string,
	r *http.Request) (bIsPass bool, accessToken string, err error) {
	bIsPass = false
	var cookie *http.Cookie
	cookie, err = r.Cookie("access_token")
	if err != nil {
		strErrMsg := err.Error()
		//"http: named cookie not present",该错误信息，不能修改，该信息是http包里，返回的。
		if action == "Login" && strings.Contains(strErrMsg, "http: named cookie not present") {
			return true, "", nil
		}
		return
	}
	accessToken = cookie.Value
	at.loginCheckMu.RLock()
	defer at.loginCheckMu.RUnlock()
	if _, ok := at.tokenToOptIDMap[accessToken]; ok {
		if action == "Login" {
			err = errors.New("the user has been to login,please logout the user first")
		} else {
			bIsPass = true
		}
	}
	return
}

func (at *AccessTokenHandler) ExpirationCheck() {
	at.logger.LogDebug("expirationCheck,begin")
	tick := time.NewTicker(time.Second * 60)
	defer tick.Stop()
	for {
		select {
		case <-at.quitCheckCh:
			goto end
		case <-tick.C:
			at.expirationCheckMu.RLock()
			var expirationToken = []string{}
			curTime := time.Now().Unix()
			for k, v := range at.tokenToTimeMap {
				if curTime > v {
					expirationToken = append(expirationToken, k)
					at.logger.Debug("[has been expire.] [token:time[%s:%v];curTime:%v;times:%d][in ExpirationCheck]", k, v, curTime, at.checkStatusCount)
				}
			}
			at.expirationCheckMu.RUnlock()
			at.delBatchToken(expirationToken)
		}
	}
end:
	at.logger.LogDebug("expirationCheck,end")
	//return
}

func (at *AccessTokenHandler) QuitExpirationCheckService() {
	//at.quitCheckCh <- true
	close(at.quitCheckCh)
}

func (at *AccessTokenHandler) isRootToken(accessToken string) bool {
	bIsRoot := false
	at.loginCheckMu.RLock()
	usrInfo, ok := at.tokenToOptIDMap[accessToken]
	if ok {
		if usrInfo.iOperatorID == 101 {
			bIsRoot = true
		}
	}
	at.loginCheckMu.RUnlock()
	return bIsRoot
}

func (at *AccessTokenHandler) isAdminToken(accessToken string) bool {
	bIsAdmin := false
	//获取在线用户的角色
	iRole := at.getOperatorRole(accessToken)
	if (iRole & 0xF0) > 0 {
		bIsAdmin = true
	}
	return bIsAdmin
}

func (at *AccessTokenHandler) isRootRequest(r *http.Request) bool {
	bIsRoot := false
	cookie, err := r.Cookie("access_token")
	if err != nil {
		at.logger.LogDebug("r.Cookie,failed,errMsg:", err.Error())
		return bIsRoot
	}
	return at.isRootToken(cookie.Value)
}

func (at *AccessTokenHandler) isAdminRequest(r *http.Request) bool {
	bIsAdmin := false
	cookie, err := r.Cookie("access_token")
	if err != nil {
		at.logger.LogError("r.Cookie,failed,errMsg:", err.Error())
		return bIsAdmin
	}
	return at.isAdminToken(cookie.Value)
}

func (at *AccessTokenHandler) getOperatorRole(accessToken string) int {
	at.loginCheckMu.Lock()
	defer at.loginCheckMu.Unlock()
	usrInfo, ok := at.tokenToOptIDMap[accessToken]
	if ok {
		if usrInfo.iRole == 0 {
			ctx := context.TODO()
			infoView, err := at.optInfoService.GetOperatorInfoByID(ctx, usrInfo.iOperatorID,
				"getOperatorRole")
			if err != nil {
				at.logger.LogError("GetOperatorInfoByID,failed,iOperatorID:", usrInfo.iOperatorID,
					"errMsg:", err.Error())
				return 0
			}
			usrInfo.iRole = infoView.Role
		}
	} else {
		return 0
	}
	return usrInfo.iRole
}
