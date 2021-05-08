package business

import (
	"analysis-server/sdk/options"
	"client/util"
	"encoding/json"
	"sync"
)

type Authen struct {
	strUserName string
	//strPasswd      string
	strAccessToken string
	userStatus     int
	mu             sync.Mutex
}

func (auth *Authen) setAuthenInfo(strUserName, strAccessToken string, userStatus int) {
	auth.strUserName = strUserName
	auth.strAccessToken = strAccessToken
	auth.userStatus = userStatus
	logger.LogInfo("authenInfo: {user name:", strUserName, ";access_token:", strAccessToken,
		";userStatus:", userStatus, "}")
}

func (auth *Authen) GetUserStatus() int {
	auth.mu.Lock()
	defer auth.mu.Unlock()
	return auth.userStatus
}

func (auth *Authen) setUserStatus(status int) {
	auth.mu.Lock()
	defer auth.mu.Unlock()
	auth.userStatus = status
}

func (auth *Authen) UserLogin(param []byte) int {
	logger.LogInfo("userLogin begin")
	errCode := util.ErrNull
	var opts options.AuthenInfoOptions
	if err := json.Unmarshal(param, &opts); err != nil {
		logger.Error("the Unmarshal failed,err:%v", err.Error())
		errCode = util.ErrUnmarshalFailed
		return errCode
	}
	if view, err := cSdk.Login(&opts); err != nil {
		errCode = util.ErrUserLoginFailed
		logger.Error("the Login failed,err:%v", err.Error())
	} else {
		logger.Debug("Login succeed;view:%v", view)
		auth.setAuthenInfo(opts.Name, view.AccessToken, util.Online)
	}
	logger.LogInfo("userLogin end")
	return errCode
}

func (auth *Authen) Logout() int {
	logger.LogInfo("logout,begin")
	errCode := util.ErrNull
	var opts options.LogoutOptions
	opts.Name = auth.strUserName
	opts.AccessToken = auth.strAccessToken
	if err := cSdk.Logout(&opts); err != nil {
		errCode = util.ErrUserLogoutFailed
		logger.Error("the Logout failed,err:%v", err.Error())
	} else {
		logger.Debug("logout succeed")
		auth.setAuthenInfo("", "", util.Offline)
	}
	logger.LogInfo("logout,end")
	return errCode
}
func (auth *Authen) OnlineCheck() int {
	errCode := util.ErrNull
	//userStatus := util.InvalidStatus
	var opts options.NameOptions
	opts.Name = auth.strUserName
	if view, err := cSdk.StatusCheckout(&opts); err != nil {
		errCode = util.ErrUserLogoutFailed
		logger.Error("OnlineCheck failed,err:%v", err.Error())
	} else {
		//userStatus = view.Status
		if auth.GetUserStatus() != view.Status {
			auth.setAuthenInfo("", "", view.Status)
			logger.Debug("OnlineCheck succeed;but the user status has been to change,the new user status is:%v", view)
		}
	}
	return errCode
}
