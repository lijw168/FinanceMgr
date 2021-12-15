package handler

import (
	"net/http"
	"unicode/utf8"

	"analysis-server/api/service"
	"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
)

type CompanyHandlers struct {
	CCHandler
	Logger     *log.Logger
	ComService *service.CompanyService
}

func (ch *CompanyHandlers) ListCompany(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListCompanyParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/Listcompany] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrCompany, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["companyId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["companyGroupId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["companyName"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["abbreName"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["corporator"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["phone"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["e_mail"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["companyAddr"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["backup"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["startAccountPeriod"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["latestAccountYear"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrCompany, service.ErrValue, service.ErrField, service.ErrNull)
			ch.Response(r.Context(), ch.Logger, w, ce, nil)
			return
		}
	}
	if (params.Order != nil) && (len(params.Order) > 0) {
		switch *params.Order[0].Field {
		case "createdAt":
			*params.Order[0].Field = "createdAt"
		case "updatedAt":
			*params.Order[0].Field = "updatedAt"
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
			ch.Response(r.Context(), ch.Logger, w, ce, nil)
			return
		}
		switch *params.Order[0].Direction {
		case cons.Order_Asc, cons.Order_Desc:
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
			ch.Response(r.Context(), ch.Logger, w, ce, nil)
			return
		}
	}
	if (params.DescOffset != nil) && (*params.DescOffset < 0) {
		ce := service.NewError(service.ErrCompany, service.ErrInvalid, service.ErrOffset, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ce, nil)
		return
	}
	if (params.DescLimit != nil) && (*params.DescLimit < -1) {
		ce := service.NewError(service.ErrCompany, service.ErrInvalid, service.ErrLimit, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ce, nil)
		return
	}

	comViews, count, ccErr := ch.ComService.ListCompany(r.Context(), params)
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(), "[company/ListCompany/ServerHTTP] [ComService.ListCompany: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), comViews}
	ch.Response(r.Context(), ch.Logger, w, nil, dataBuf)
	return
}

func (ch *CompanyHandlers) GetCompany(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/GetCompany] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrCompany, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrCompany, service.ErrMiss, service.ErrId, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	requestId := ch.GetTraceId(r)

	comView, ccErr := ch.ComService.GetCompanyById(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(), "[company/GetCompanyById/ServerHTTP] [ComService.GetCompanyById: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ch.Response(r.Context(), ch.Logger, w, nil, comView)
	return
}

func (ch *CompanyHandlers) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var params = new(model.CreateCompanyParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/CreateCompany] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrCompany, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyName == nil || *params.CompanyName == "" {
		ccErr := service.NewError(service.ErrCompany, service.ErrMiss, service.ErrName, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	if utf8.RuneCountInString(*params.CompanyName) > NameMaxLen || !utils.VerStrP(*params.CompanyName) {
		ccErr := service.NewError(service.ErrCompany, service.ErrInvalid, service.ErrName, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	if params.StartAccountPeriod == nil || *params.StartAccountPeriod <= 0 {
		ccErr := service.NewError(service.ErrCompany, service.ErrMiss, service.ErrAccountPeriod, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	requestId := ch.GetTraceId(r)

	comView, ccErr := ch.ComService.CreateCompany(r.Context(), params, requestId)
	ch.Logger.InfoContext(r.Context(), "CreateCompany in CreateCompany.")
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(), "[company/CreateCompany/ServerHTTP] [ComService.CreateCompany: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ch.Response(r.Context(), ch.Logger, w, nil, comView)
	return
}

func (ch *CompanyHandlers) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ModifyCompanyParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/UpdateCompany] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrCompany, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrCompany, service.ErrMiss, service.ErrId, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	updateFields := make(map[string]interface{})
	if params.CompanyName != nil {
		if utf8.RuneCountInString(*params.CompanyName) > NameMaxLen || !utils.VerStrP(*params.CompanyName) {
			ccErr := service.NewError(service.ErrCompany, service.ErrInvalid, service.ErrName, service.ErrNull)
			ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
			return
		}
		updateFields["CompanyName"] = *params.CompanyName
	}
	if params.AbbrevName != nil {
		if utf8.RuneCountInString(*params.AbbrevName) > NameMaxLen || !utils.VerStrP(*params.AbbrevName) {
			ccErr := service.NewError(service.ErrCompany, service.ErrInvalid, service.ErrName, service.ErrNull)
			ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
			return
		}
		updateFields["AbbrevName"] = *params.AbbrevName
	}
	if params.Corporator != nil {
		updateFields["Corporator"] = *params.Corporator
	}
	if params.Phone != nil {
		updateFields["Phone"] = *params.Phone
	}
	if params.Email != nil {
		updateFields["E_mail"] = *params.Email
	}
	if params.CompanyAddr != nil {
		updateFields["CompanyAddr"] = *params.CompanyAddr
	}
	if params.Backup != nil {
		updateFields["Backup"] = *params.Backup
	}
	if params.LatestAccountYear != nil {
		updateFields["LatestAccountYear"] = *params.LatestAccountYear
	}
	if len(updateFields) == 0 {
		ccErr := service.NewError(service.ErrCompany, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ccErr := ch.ComService.UpdateCompanyById(r.Context(), *params.CompanyID, updateFields)
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(), "[company/UpdateCompany/ServerHTTP] [ComService.UpdateCompanyById: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ch.Response(r.Context(), ch.Logger, w, nil, nil)
}

func (ch *CompanyHandlers) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteIDParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/DeleteCompany] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrCompany, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrCompany, service.ErrMiss, service.ErrId, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	requestId := ch.GetTraceId(r)
	ccErr := ch.ComService.DeleteCompanyByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(), "[company/DeleteCompany/ServerHTTP] [ComService.DeleteCompanyByID: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ch.Response(r.Context(), ch.Logger, w, nil, nil)
	return
}

func (ch *CompanyHandlers) AssociatedCompanyGroup(w http.ResponseWriter, r *http.Request) {
	var params = new(model.AssociatedCompanyGroupParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/AssociatedCompanyGroup] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrCompany, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrCompany, service.ErrMiss, service.ErrId, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	if params.CompanyGroupID == nil || *params.CompanyGroupID <= 0 {
		ccErr := service.NewError(service.ErrCompany, service.ErrMiss, service.ErrId, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	if params.IsAttach == nil {
		ccErr := service.NewError(service.ErrCompany, service.ErrMiss, service.ErrAttachParam, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	requestId := ch.GetTraceId(r)
	ccErr := ch.ComService.AssociatedCompanyGroup(r.Context(), params, requestId)
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(),
			"[company/AssociatedCompanyGroup/ServerHTTP] [ComService.AssociatedCompanyGroup: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ch.Response(r.Context(), ch.Logger, w, nil, nil)
}
