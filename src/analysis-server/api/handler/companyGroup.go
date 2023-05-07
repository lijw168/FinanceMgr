package handler

import (
	"net/http"
	"strconv"
	"unicode/utf8"

	"financeMgr/src/analysis-server/api/service"
	"financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
	"financeMgr/src/common/log"
)

type CompanyGroupHandlers struct {
	CCHandler
	Logger          *log.Logger
	ComGroupService *service.CompanyGroupService
}

func (ch *CompanyGroupHandlers) ListCompanyGroup(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/ListCompanyGroup] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrComGroup, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["companyGroupId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["groupName"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["groupStatus"] = utils.Attribute{Type: utils.T_Int, Val: nil}

		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrComGroup, service.ErrValue, service.ErrField, service.ErrNull)
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
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, strconv.Itoa(*params.Order[0].Direction))
			ch.Response(r.Context(), ch.Logger, w, ce, nil)
			return
		}
	}
	if (params.DescOffset != nil) && (*params.DescOffset < 0) {
		ce := service.NewError(service.ErrComGroup, service.ErrInvalid, service.ErrOffset, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ce, nil)
		return
	}
	if (params.DescLimit != nil) && (*params.DescLimit < -1) {
		ce := service.NewError(service.ErrComGroup, service.ErrInvalid, service.ErrLimit, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ce, nil)
		return
	}

	comGroupViews, count, ccErr := ch.ComGroupService.ListCompanyGroup(r.Context(), params)
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(), "[company/ListCompanyGroup/ServerHTTP] [ComGroupService.ListCompanyGroup: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), comGroupViews}
	ch.Response(r.Context(), ch.Logger, w, nil, dataBuf)
	return
}

func (ch *CompanyGroupHandlers) GetCompanyGroup(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/GetCompany] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrComGroup, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrComGroup, service.ErrMiss, service.ErrId, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	requestId := ch.GetTraceId(r)

	comView, ccErr := ch.ComGroupService.GetCompanyGroupById(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(),
			"[company/GetCompanyGroupById/ServerHTTP] [ComGroupService.GetCompanyGroupById: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ch.Response(r.Context(), ch.Logger, w, nil, comView)
	return
}

func (ch *CompanyGroupHandlers) CreateCompanyGroup(w http.ResponseWriter, r *http.Request) {
	var params = new(model.CreateCompanyGroupParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/CreateCompany] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrComGroup, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	if params.GroupName == nil || *params.GroupName == "" {
		ccErr := service.NewError(service.ErrComGroup, service.ErrMiss, service.ErrName, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	if utf8.RuneCountInString(*params.GroupName) > NameMaxLen || !utils.VerStrP(*params.GroupName) {
		ccErr := service.NewError(service.ErrComGroup, service.ErrInvalid, service.ErrName, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	if params.GroupStatus == nil || *params.GroupStatus < 0 {
		ccErr := service.NewError(service.ErrComGroup, service.ErrMiss, service.ErrStatus, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	requestId := ch.GetTraceId(r)

	comGroupView, ccErr := ch.ComGroupService.CreateCompanyGroup(r.Context(), params, requestId)
	ch.Logger.InfoContext(r.Context(), "CreateCompanyGroup in CreateCompanyGroup.")
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(),
			"[company/CreateCompanyGroup/ServerHTTP] [ComGroupService.CreateCompanyGroup: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ch.Response(r.Context(), ch.Logger, w, nil, comGroupView)
	return
}

func (ch *CompanyGroupHandlers) UpdateCompanyGroup(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ModifyCompanyGroupParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/UpdateCompanyGroup] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrComGroup, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	if params.CompanyGroupID == nil || *params.CompanyGroupID <= 0 {
		ccErr := service.NewError(service.ErrComGroup, service.ErrMiss, service.ErrId, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}

	updateFields := make(map[string]interface{})
	if params.GroupName != nil {
		if utf8.RuneCountInString(*params.GroupName) > NameMaxLen || !utils.VerStrP(*params.GroupName) {
			ccErr := service.NewError(service.ErrComGroup, service.ErrInvalid, service.ErrName, service.ErrNull)
			ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
			return
		}
		updateFields["groupName"] = *params.GroupName
	}
	if params.GroupStatus != nil || *params.GroupStatus > 0 {
		updateFields["groupStatus"] = *params.GroupStatus
	}
	if len(updateFields) == 0 {
		ccErr := service.NewError(service.ErrComGroup, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ccErr := ch.ComGroupService.UpdateCompanyGroupById(r.Context(), *params.CompanyGroupID, updateFields)
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(), "[company/UpdateCompanyGroup/ServerHTTP] [ComGroupService.UpdateCompanyGroupById: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ch.Response(r.Context(), ch.Logger, w, nil, nil)
}

func (ch *CompanyGroupHandlers) DeleteCompanyGroup(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteIDParams)
	err := ch.HttpRequestParse(r, params)
	if err != nil {
		ch.Logger.ErrorContext(r.Context(), "[company/DeleteCompanyGroup] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrComGroup, service.ErrMalformed, service.ErrNull, err.Error())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrComGroup, service.ErrMiss, service.ErrId, service.ErrNull)
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	requestId := ch.GetTraceId(r)
	ccErr := ch.ComGroupService.DeleteCompanyGroupByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ch.Logger.WarnContext(r.Context(), "[company/DeleteCompanyGroup/ServerHTTP] [ComGroupService.DeleteCompanyGroupByID: %s]", ccErr.Detail())
		ch.Response(r.Context(), ch.Logger, w, ccErr, nil)
		return
	}
	ch.Response(r.Context(), ch.Logger, w, nil, nil)
	return
}
