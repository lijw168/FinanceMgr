package business

import (
	//"financeMgr/src/analysis-server/model"
	"financeMgr/src/analysis-server/sdk/options"
	sUtil "financeMgr/src/analysis-server/sdk/util"
	sdkUtil "financeMgr/src/analysis-server/sdk/util"
	"financeMgr/src/client/util"
	"sync"
)

type Authen struct {
	OperatorID     int
	strUserName    string
	strAccessToken string
	userStatus     int
	mu             sync.Mutex
}

func (auth *Authen) setAuthenInfo(strUserName, strAccessToken string, iOperatorID, iUserStatus int) {
	auth.OperatorID = iOperatorID
	auth.strUserName = strUserName
	auth.strAccessToken = strAccessToken
	auth.userStatus = iUserStatus
	logger.LogInfo("authenInfo: {operatorId:", iOperatorID, ";user name:", strUserName,
		";access_token:", strAccessToken, ";userStatus:", iUserStatus, "}")
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

func (auth *Authen) UserLogin(param []byte) (errCode int, errMsg string) {
	logger.LogInfo("userLogin begin")
	errCode = util.ErrNull
	if view, err := cSdk.Login_json(param); err != nil {
		errCode = util.ErrUserLoginFailed
		logger.Error("the Login failed,err:%v", err.Error())
		if resErr, ok := err.(*sdkUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrUserLoginFailed
			errMsg = "the Login failed,internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("Login succeed;view:%v", view)
		auth.setAuthenInfo(view.Name, view.AccessToken, view.OperatorID, view.Status)
		cSdk.SetAccessToken(view.AccessToken)
	}
	logger.LogInfo("userLogin end")
	return errCode, errMsg
}

func (auth *Authen) Logout() (errCode int, errMsg string) {
	logger.LogInfo("logout,begin")
	errCode = util.ErrNull
	var opts options.BaseOptions
	opts.ID = auth.OperatorID
	if err := cSdk.Logout(&opts); err != nil {
		if resErr, ok := err.(*sUtil.RespErr); ok {
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrUserLogoutFailed
			errMsg = "the Logout failed, internal error"
		}
		logger.LogError(errMsg)
	} else {
		logger.Debug("logout succeed")
	}
	auth.setAuthenInfo("", "", 0, util.Offline)
	cSdk.SetAccessToken("")
	logger.LogInfo("logout,end")
	return errCode, errMsg
}

func (auth *Authen) OnlineCheck() (errCode int, errMsg string) {
	errCode = util.ErrNull
	var opts options.BaseOptions
	opts.ID = auth.OperatorID
	if view, err := cSdk.StatusCheckout(&opts); err != nil {
		if resErr, ok := err.(*sUtil.RespErr); ok {
			//errInfo:please login first
			if resErr.Code == -2 {
				auth.setAuthenInfo("", "", 0, util.Offline)
				cSdk.SetAccessToken("")
				logger.Debug("the user status has been to convert to the %d", util.Offline)
			}
			errCode = resErr.Code
			errMsg = resErr.Err.Error()
		} else {
			errCode = util.ErrOnlineCheckout
			errMsg = "the OnlineCheck failed,internal error"
		}
		logger.LogError(errMsg)

	} else {
		if auth.GetUserStatus() != view.Status {
			auth.setAuthenInfo("", "", 0, view.Status)
			cSdk.SetAccessToken("")
			logger.Debug("OnlineCheck succeed;the user status is %d", view.Status)
		}
	}
	return errCode, errMsg
}

func (auth *Authen) ListLoginInfo(param []byte) (resData []byte, errCode int, errMsg string) {
	return listCmdJson(resource_type_login_info, param, cSdk.ListLoginInfo_json)
}
