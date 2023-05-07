package handler

import (
	"financeMgr/src/analysis-server/api/service"
	"financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
	"financeMgr/src/common/log"
	"net/http"
	"strconv"
	"strings"
	//"time"
)

type AuthenHandlers struct {
	CCHandler
	Logger         *log.Logger
	AuthService    *service.AuthenService
	ComService     *service.CompanyService
	OptInfoService *service.OperatorInfoService
}

func (ah *AuthenHandlers) ListLoginInfo(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[AuthenHandlers/ListLoginInfo/ServerHTTP] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrLogin, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if isLackBaseParams([]string{"operatorId"}, params.Filter) {
		if !GAccessTokenH.isRootRequest(r) {
			ah.Logger.ErrorContext(r.Context(), "lack base param  operatorId")
			ce := service.NewError(service.ErrLogin, service.ErrMiss, service.ErrId, service.ErrNull)
			ah.Response(r.Context(), ah.Logger, w, ce, nil)
			return
		}
	}
	//todo:通过operatorId，在operatorInfo表里获取所有属于该公司的所有操作员的登录信息。
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["operatorId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["name"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["clientIp"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["status"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrLogin, service.ErrInvalid, service.ErrField, service.ErrNull)
			ah.Response(r.Context(), ah.Logger, w, ce, nil)
			return
		}
	}
	if (params.Order != nil) && (len(params.Order) > 0) {
		switch *params.Order[0].Field {
		case "BeginedAt":
			*params.Order[0].Field = "begined_at"
		case "EndedAt":
			*params.Order[0].Field = "ended_at"
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
			ah.Response(r.Context(), ah.Logger, w, ce, nil)
			return
		}
		switch *params.Order[0].Direction {
		case cons.Order_Asc, cons.Order_Desc:
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, strconv.Itoa(*params.Order[0].Direction))
			ah.Response(r.Context(), ah.Logger, w, ce, nil)
			return
		}
	}
	if (params.DescOffset != nil) && (*params.DescOffset < 0) {
		ce := service.NewError(service.ErrLogin, service.ErrInvalid, service.ErrOffset, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ce, nil)
		return
	}
	if (params.DescLimit != nil) && (*params.DescLimit < -1) {
		ce := service.NewError(service.ErrLogin, service.ErrInvalid, service.ErrLimit, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ce, nil)
		return
	}

	optViews, count, ccErr := ah.AuthService.ListLoginInfo(r.Context(), params)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[AuthenHandlers/ListLoginInfo/ServerHTTP] [AuthService.ListLoginInfo: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), optViews}
	ah.Response(r.Context(), ah.Logger, w, nil, dataBuf)
	return
}

func (ah *AuthenHandlers) StatusCheckout(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[AuthenHandlers/StatusCheckout] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrNull, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}

	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrNull, service.ErrMiss, service.ErrId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	requestId := ah.GetTraceId(r)
	stCheckoutView, ccErr := ah.AuthService.StatusCheckout(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[AuthenHandlers/StatusCheckout] [AuthService.StatusCheckout: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	//add lease
	cookie, err := r.Cookie("access_token")
	if err != nil {
		ah.Logger.ErrorContext(r.Context(),
			"[AuthenHandlers/StatusCheckout] :get cookie access_token,failed.err: %v]", err)
		return
	}
	GAccessTokenH.modifyTokenExpiredTime(cookie.Value)

	ah.Response(r.Context(), ah.Logger, w, nil, stCheckoutView)
	return
}

func (ah *AuthenHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var params = new(model.AuthenInfoParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[AuthenHandlers/login/ServerHTTP] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrLogin, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil && *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrLogin, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	requestId := ah.GetTraceId(r)
	if params.Name == nil && *params.Name == "" {
		ccErr := service.NewError(service.ErrLogin, service.ErrMiss, service.ErrName, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.Password == nil {
		ccErr := service.NewError(service.ErrLogin, service.ErrMiss, service.ErrPasswd, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	//authentication the user and password
	optView, ccErr := ah.OptInfoService.GetOperatorInfoByName(r.Context(), *params.Name, *params.CompanyID, requestId)
	if ccErr != nil {
		if ccErr.GetCode() == cons.CodeOptInfoNotExist {
			ccErr.SetCode(cons.CodeUserNameWrong)
		}
		ah.Logger.ErrorContext(r.Context(), "[AuthenHandlers/login/ServerHTTP] [authentication failed. error: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	//暂时先这样认证
	if *params.Password != optView.Password {
		ah.Logger.ErrorContext(r.Context(), "[AuthenHandlers/login/ServerHTTP] [authentication failed. error: the password is wrong]")
		ccErr := service.NewCcError(cons.CodePasswdWrong, service.ErrLogin, service.ErrInvalid, service.ErrPasswd, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	//generate login information
	var loginInfo model.LoginInfoParams
	loginInfo.Name = params.Name
	loginInfo.OperatorID = &(optView.OperatorID)
	clientAddr := (strings.Split(r.RemoteAddr, ":"))[0]
	loginInfo.ClientIp = &clientAddr
	logInfoView, ccErr := ah.AuthService.Login(r.Context(), &loginInfo, requestId)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[AuthenHandlers/login/ServerHTTP] [AuthService.Login: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Logger.InfoContext(r.Context(), "login succeed.")
	GAccessTokenH.insertToken(logInfoView.AccessToken, logInfoView.OperatorID)
	ah.Response(r.Context(), ah.Logger, w, nil, logInfoView)
	return
}

func (ah *AuthenHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[AuthenHandlers/Logout/ServerHTTP] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrLogout, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrLogout, service.ErrMiss, service.ErrId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ccErr := ah.AuthService.Logout(r.Context(), *params.ID)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[AuthenHandlers/Logout/ServerHTTP] [AuthService.Logout: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	cookie, err := r.Cookie("access_token")
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[AuthenHandlers/Logout/ServerHTTP] [Get AccessToken,failed: %v]", err)
		ccErr := service.NewError(service.ErrLogout, service.ErrMiss, service.ErrCookie, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	GAccessTokenH.delToken(cookie.Value)
	ah.Logger.InfoContext(r.Context(), "logout succeed.")
	ah.Response(r.Context(), ah.Logger, w, nil, nil)
	return
}
