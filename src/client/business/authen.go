package business

import (
	"analysis-server/sdk/options"
	"client/util"
	"encoding/json"
	"sync"
)

type Authen struct {
	strUserName    string
	strPasswd      string
	strAccessToken string
	userStatus     int
	mu             sync.Mutex
}

func (auth *Authen) SetUserInfo(strUserName, strPasswd string) {
	auth.strUserName = strUserName
	auth.strPasswd = strPasswd
	logger.LogInfo("user name:", strUserName, ";passwd:", strPasswd)
	auth.userStatus = util.UserOffline
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

func (auth *Authen) UserLogin(param []byte) (errCode int) {
	logger.LogInfo("userLogin begin")
	errCode = util.ErrNull
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
		auth.setUserStatus(util.UserOnline)
	}
	logger.LogInfo("userLogin end")
	return errCode
}

func (auth *Authen) Logout() (errCode int) {
	logger.LogInfo("logout,begin")
	errCode = util.ErrNull
	var opts options.NameOptions
	opts.Name = auth.strUserName
	if err := cSdk.Logout(&opts); err != nil {
		errCode = util.ErrUserLogoutFailed
		logger.Error("the Logout failed,err:%v", err.Error())
	} else {
		logger.Debug("logout succeed")
		auth.setUserStatus(util.UserOffline)
	}
	logger.LogInfo("logout,end")
	return errCode
}
func (auth *Authen) UserOnlineCheck() int {
	//fmt.Println("userOnlineCheck,begin")
	logger.LogInfo("userOnlineCheck,begin")
	resCode := util.ErrNull
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// runNumber := r.Int()
	// base64UserName := util.Base64Ex([]byte(auth.strUserName))
	// strUrl := "https://" + auth.strHost + ":6688/webui/login/sslvpn_usercheck?username="
	// strUrl = strUrl + base64UserName + "&run=" + strconv.Itoa(runNumber) + "&haha110=haha120"
	// request, err := auth.generateRequest(strUrl, "GET", nil)
	// if request == nil {
	// 	//fmt.Println("generateRequest,failed")
	// 	logger.LogError("generateRequest,failed")
	// 	resCode = ErrGenHttpReqFailed
	// 	return resCode
	// }
	// urlEncodeName := url.QueryEscape(base64UserName)
	// auth.addCookie(request, "access_token", auth.strAccessToken)
	// //auth.addCookie(request, "username", urlEncodeName)
	// auth.addCookie(request, "user", urlEncodeName)
	// reqHeader := "reqHeader dump:"
	// if dump, err := httputil.DumpRequest(request, false); err == nil {
	// 	reqHeader += string(dump)
	// }
	// //fmt.Println(reqHeader)
	// logger.LogDebug(reqHeader)

	// // tr := &http.Transport{
	// // 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// // }
	// var resp *http.Response
	// //client := &http.Client{Transport: tr}
	// client := &http.Client{Transport: auth.tr}
	// resp, err = client.Do(request)
	// if err == nil && resp.StatusCode == http.StatusOK {
	// 	//fmt.Println("client.Do, succed")
	// 	logger.LogInfo("client.Do, succed")
	// } else {
	// 	//fmt.Println("Get error:", err)
	// 	logger.LogError("userOnlineCheck request,failed, error:", err.Error())
	// 	resCode = ErrHttpReqFailed
	// 	return resCode
	// }
	// if dump, err2 := httputil.DumpResponse(resp, true); err2 == nil {
	// 	putResult := fmt.Sprintf("status_code:%d,response body:%v", resp.StatusCode, string(dump))
	// 	//fmt.Println(putResult)
	// 	logger.LogInfo(putResult)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.LogError("read body ,failed, error:", err.Error())
	// 	resCode = ErrUnmarshFailed
	// } else {
	// 	logger.LogDebug("json:", string(body))
	// }
	return resCode
}
