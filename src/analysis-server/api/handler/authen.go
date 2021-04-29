package handler

import (
	"analysis-server/api/service"
	"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"net/http"
	"strings"
	"time"
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
		ah.Logger.ErrorContext(r.Context(), "[loginInfo/ListLoginInfo] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrLogin, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["name"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["clientIp"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["role"] = utils.Attribute{Type: utils.T_Int, Val: nil}
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
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
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
		ah.Logger.WarnContext(r.Context(), "[loginInfo/ListLoginInfo/ServerHTTP] [AuthService.ListLoginInfo: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), optViews}
	ah.Response(r.Context(), ah.Logger, w, nil, dataBuf)
	return
}

// func (ah *AuthenHandlers) GetLoginInfo(w http.ResponseWriter, r *http.Request) {
// 	var params = new(model.DescribeNameParams)
// 	err := ah.HttpRequestParse(r, params)
// 	if err != nil {
// 		ah.Logger.ErrorContext(r.Context(), "[loginInfo/GetLoginInfo] [HttpRequestParse: %v]", err)
// 		ccErr := service.NewError(service.ErrLogin, service.ErrMalformed, service.ErrNull, err.Error())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}

// 	if params.Name == nil || *params.Name == "" {
// 		ccErr := service.NewError(service.ErrLogin, service.ErrMiss, service.ErrName, service.ErrNull)
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	requestId := ah.GetTraceId(r)
// 	optView, ccErr := ah.AuthService.GetLoginInfoByName(r.Context(), *params.Name, requestId)
// 	if ccErr != nil {
// 		ah.Logger.WarnContext(r.Context(), "[loginInfo/GetLoginInfo/ServerHTTP] [AuthService.GetLoginInfoByName: %s]", ccErr.Detail())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	ah.Response(r.Context(), ah.Logger, w, nil, optView)
// 	return
// }

func (ah *AuthenHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var params = new(model.AuthenInfoParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[Login] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrLogin, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil && *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrLogin, service.ErrMiss, service.ErrId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	//check out ,whether the company exists.
	requestId := ah.GetTraceId(r)
	comView, ccErr := ah.ComService.GetCompanyById(r.Context(), *params.CompanyID, requestId)
	if comView == nil || ccErr != nil {
		ah.Logger.ErrorContext(r.Context(), "the company is not exist. companyId:%d", *params.CompanyID)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
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
	optView, ccErr := ah.OptInfoService.GetOperatorInfoByName(r.Context(), *params.Name, requestId)
	if ccErr != nil {
		ah.Logger.ErrorContext(r.Context(), "[login/ServerHTTP] [authentication failed. error: %s]", ccErr.Detail())
		if ccErr.GetCode() == cons.CodeOptInfoNotExist {
			ccErr.SetCode(cons.CodeUserNameWrong)
		}
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if *params.Password != optView.Password {
		ah.Logger.ErrorContext(r.Context(), "[login/ServerHTTP] [authentication failed. error: the password is wrong]")
		ccErr := service.NewCcError(cons.CodePasswdWrong, service.ErrLogin, service.ErrInvalid, service.ErrPasswd, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	//generate login information
	var loginInfo model.LoginInfoParams
	loginInfo.Name = params.Name
	loginInfo.Role = &optView.Role
	clientAddr := (strings.Split(r.RemoteAddr, ":"))[0]
	loginInfo.ClientIp = &clientAddr
	optInfoView, ccErr := ah.AuthService.Login(r.Context(), &loginInfo, requestId)
	ah.Logger.InfoContext(r.Context(), "CreateLoginInfo in login.")
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[login/CreateLoginInfo/ServerHTTP] [AuthService.CreateLoginInfo: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, optInfoView)
	return
}

func (ah *AuthenHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeNameParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[volumes/Logout] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrLogin, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.Name == nil || *params.Name == "" {
		ccErr := service.NewError(service.ErrLogin, service.ErrMiss, service.ErrName, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	updateFields := make(map[string]interface{})
	updateFields["EndedAt"] = time.Now()

	ccErr := ah.AuthService.Logout(r.Context(), *params.Name, updateFields)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[opreator/Logout/ServerHTTP] [AuthService.UpdateLoginInfo: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, nil)
	return
}
