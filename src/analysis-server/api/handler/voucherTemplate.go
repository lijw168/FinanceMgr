package handler

import (
	"net/http"
	//"unicode/utf8"

	"analysis-server/api/service"
	"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
)

type VoucherTemplateHandlers struct {
	CCHandler
	Logger             *log.Logger
	VoucherTempService *service.VoucherTemplateService
}

func (vt *VoucherTemplateHandlers) ListVoucherTemplate(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListParams)
	err := vt.HttpRequestParse(r, params)
	if err != nil {
		vt.Logger.ErrorContext(r.Context(), "[voucherTemplate/ListVoucherTemplate] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucherTemplate, service.ErrMalformed, service.ErrNull, err.Error())
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["voucherTemplateId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["refVoucherId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["voucherYear"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["illustration"] = utils.Attribute{Type: utils.T_String, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrVoucherTemplate, service.ErrInvalid, service.ErrField, service.ErrNull)
			vt.Response(r.Context(), vt.Logger, w, ce, nil)
			return
		}
	}
	if (params.Order != nil) && (len(params.Order) > 0) {
		switch *params.Order[0].Field {
		case "voucherTemplateId":
			*params.Order[0].Field = "voucherTemplateId"
		case "createdAt":
			*params.Order[0].Field = "createdAt"
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
			vt.Response(r.Context(), vt.Logger, w, ce, nil)
			return
		}
		switch *params.Order[0].Direction {
		case cons.Order_Asc, cons.Order_Desc:
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
			vt.Response(r.Context(), vt.Logger, w, ce, nil)
			return
		}
	}
	if (params.DescOffset != nil) && (*params.DescOffset < 0) {
		ce := service.NewError(service.ErrVoucherTemplate, service.ErrInvalid, service.ErrOffset, service.ErrNull)
		vt.Response(r.Context(), vt.Logger, w, ce, nil)
		return
	}
	if (params.DescLimit != nil) && (*params.DescLimit < -1) {
		ce := service.NewError(service.ErrVoucherTemplate, service.ErrInvalid, service.ErrLimit, service.ErrNull)
		vt.Response(r.Context(), vt.Logger, w, ce, nil)
		return
	}

	tmpViews, count, ccErr := vt.VoucherTempService.ListVoucherTemplate(r.Context(), params)
	if ccErr != nil {
		vt.Logger.ErrorContext(r.Context(), "[voucherTemplate/ListVoucherTemplate/ServerHTTP] [VoucherTempService.ListVoucherTemplate: %s]", ccErr.Detail())
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), tmpViews}
	vt.Response(r.Context(), vt.Logger, w, nil, dataBuf)
	return
}

func (vt *VoucherTemplateHandlers) GetVoucherTemplate(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := vt.HttpRequestParse(r, params)
	if err != nil {
		vt.Logger.ErrorContext(r.Context(), "[voucherTemplate/GetVoucherTemplate] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucherTemplate, service.ErrMalformed, service.ErrNull, err.Error())
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	//the id is voucherTemplateID
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrVoucherTemplate, service.ErrMiss, service.ErrVoucherTemplateID, service.ErrNull)
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}

	requestId := vt.GetTraceId(r)
	tmpView, ccErr := vt.VoucherTempService.GetVoucherTemplate(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		vt.Logger.WarnContext(r.Context(), "[voucherTemplate/GetVoucherTemplate/ServerHTTP] [VoucherTempService.GetVoucherTemplate: %s]", ccErr.Detail())
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	vt.Response(r.Context(), vt.Logger, w, nil, tmpView)
	return
}

func (vt *VoucherTemplateHandlers) CreateVoucherTemplate(w http.ResponseWriter, r *http.Request) {
	var params = new(model.VoucherTemplateParams)
	err := vt.HttpRequestParse(r, params)
	if err != nil {
		vt.Logger.ErrorContext(r.Context(), "[voucherTemplate/CreateVoucherTemplate] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucherTemplate, service.ErrMalformed, service.ErrNull, err.Error())
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	if params.RefVoucherID == nil || *params.RefVoucherID <= 0 {
		ccErr := service.NewError(service.ErrVoucherTemplate, service.ErrMiss, service.ErrId, service.ErrNull)
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	if params.VoucherYear == nil || *params.VoucherYear <= 0 {
		ccErr := service.NewError(service.ErrVoucherTemplate, service.ErrMiss, service.ErrVouYear, service.ErrNull)
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	requestId := vt.GetTraceId(r)

	voucherTemplateID, ccErr := vt.VoucherTempService.CreateVoucherTemplate(r.Context(), params, requestId)
	if ccErr != nil {
		vt.Logger.WarnContext(r.Context(), "[voucherTemplate/CreateVoucherTemplate/ServerHTTP] [VoucherTempService.CreateVoucherTemplate: %s]", ccErr.Detail())
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	vt.Response(r.Context(), vt.Logger, w, nil, voucherTemplateID)
	return
}

func (vt *VoucherTemplateHandlers) DeleteVoucherTemplate(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteSubjectParams)
	err := vt.HttpRequestParse(r, params)
	if err != nil {
		vt.Logger.ErrorContext(r.Context(), "[vouchertemplate/DeleteOperator] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucherTemplate, service.ErrMalformed, service.ErrNull, err.Error())
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	//the id is voucherTemplateID
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrVoucherTemplate, service.ErrMiss, service.ErrVoucherTemplateID, service.ErrNull)
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	requestId := vt.GetTraceId(r)
	ccErr := vt.VoucherTempService.DeleteVoucherTemplate(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		vt.Logger.WarnContext(r.Context(), "[vouchertemplate/DeleteVoucherTemplate/ServerHTTP] [VoucherTempService.DeleteVoucherTemplate: %s]", ccErr.Detail())
		vt.Response(r.Context(), vt.Logger, w, ccErr, nil)
		return
	}
	vt.Response(r.Context(), vt.Logger, w, nil, nil)
	return
}
