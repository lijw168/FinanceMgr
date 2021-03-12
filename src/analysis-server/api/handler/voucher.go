package handler

import (
	"net/http"
	"unicode/utf8"

	"common/log"
	"analysis-server/api/service"
	"analysis-server/api/utils"
	"analysis-server/model"
)

// var (
// 	DescriptionMaxlen = 256
// 	NameMaxLen        = 32
// )

type VoucherHandlers struct {
	CCHandler
	Logger *log.Logger
	Vis    service.VoucherInfoService
	Vrs    service.VoucherRecordService
	Vs     service.VoucherService
}

func (vh *VoucherHandlers) ListVoucherInfo(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucherInfo/ListVoucherInfo] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["voucherId"] = utils.Attribute{utils.T_Int_Arr, nil}
		filterMap["companyId"] = utils.Attribute{utils.T_Int, nil}
		filterMap["voucherMonth"] = utils.Attribute{utils.T_Int_Arr, nil}
		filterMap["numOfMonth"] = utils.Attribute{utils.T_Int, nil}
		filterMap["voucherDate"] = utils.Attribute{utils.T_Int, nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrVoucher, service.ErrDesc, service.ErrField, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
	}
	// if (params.Order != nil) && (len(params.Order) > 0) {
	// 	switch *params.Order[0].Field {
	// 	case "create_time":
	// 		*params.Order[0].Field = "created_at"
	// 	case "update_time":
	// 		*params.Order[0].Field = "updated_at"
	// 	case "delete_time":
	// 		*params.Order[0].Field = "deleted_at"
	// 	default:
	// 		ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
	// 		vh.Response(r.Context(), vh.Logger, w, ce, nil)
	// 		return
	// 	}
	// 	switch *params.Order[0].Direction {
	// 	case cons.Order_Asc, cons.Order_Desc:
	// 	default:
	// 		ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
	// 		vh.Response(r.Context(), vh.Logger, w, ce, nil)
	// 		return
	// 	}
	// }
	if (params.DescOffset != nil) && (*params.DescOffset < 0) {
		ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrOffset, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ce, nil)
		return
	}
	if (params.DescLimit != nil) && (*params.DescLimit < -1) {
		ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrLimit, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ce, nil)
		return
	}
	vouInfoViews, count, ccErr := vh.Vis.ListVoucherInfo(r.Context(), params)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucherInfo/ListVoucherInfo/ServeHTTP] [VoucherInfoService.ListVoucherInfo: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), vouInfoViews}
	vh.Response(r.Context(), vh.Logger, w, nil, dataBuf)
	return
}

func (vh *VoucherHandlers) ListVoucherRecords(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucherInfo/ListVoucherRecords] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["voucherId"] = utils.Attribute{utils.T_Int, nil}
		filterMap["recordId"] = utils.Attribute{utils.T_Int_Arr, nil}
		filterMap["subjectName"] = utils.Attribute{utils.T_String, nil}
		filterMap["summary"] = utils.Attribute{utils.T_String, nil}
		filterMap["subId1"] = utils.Attribute{utils.T_Int, nil}
		filterMap["subId2"] = utils.Attribute{utils.T_Int, nil}
		filterMap["subId3"] = utils.Attribute{utils.T_Int, nil}
		filterMap["subId4"] = utils.Attribute{utils.T_Int, nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrVoucher, service.ErrDesc, service.ErrField, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
	}
	// if (params.Order != nil) && (len(params.Order) > 0) {
	// 	switch *params.Order[0].Field {
	// 	case "create_time":
	// 		*params.Order[0].Field = "created_at"
	// 	case "update_time":
	// 		*params.Order[0].Field = "updated_at"
	// 	case "delete_time":
	// 		*params.Order[0].Field = "deleted_at"
	// 	case "delete_time":
	// 		*params.Order[0].Field = "deleted_at"
	// 	default:
	// 		ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
	// 		vh.Response(r.Context(), vh.Logger, w, ce, nil)
	// 		return
	// 	}
	// 	switch *params.Order[0].Direction {
	// 	case cons.Order_Asc, cons.Order_Desc:
	// 	default:
	// 		ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
	// 		vh.Response(r.Context(), vh.Logger, w, ce, nil)
	// 		return
	// 	}
	// }
	if (params.DescOffset != nil) && (*params.DescOffset < 0) {
		ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrOffset, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ce, nil)
		return
	}
	if (params.DescLimit != nil) && (*params.DescLimit < -1) {
		ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrLimit, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ce, nil)
		return
	}

	vouRecordViews, count, ccErr := vh.Vrs.ListVoucherRecords(r.Context(), params)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucherInfo/ListVoucherRecords/ServeHTTP] [VoucherRecordService.ListVoucherRecords: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), vouRecordViews}
	vh.Response(r.Context(), vh.Logger, w, nil, dataBuf)
	return
}

func (vh *VoucherHandlers) GetVoucherInfo(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucherHandlers/GetVoucherInfo] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucherInfo, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}

	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrVoucherInfo, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	requestId := vh.GetTraceId(r)
	voucherView, ccErr := vh.Vis.GetVoucherInfoByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucherHandlers/GetVoucherInfo/ServeHTTP] [Vis.GetVoucherInfoByID: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, voucherView)
	return
}

func (vh *VoucherHandlers) GetVoucher(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucherHandlers/GetVoucher] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}

	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	requestId := vh.GetTraceId(r)
	voucherView, ccErr := vh.Vs.GetVoucherByVoucherID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucherHandlers/GetVoucher/ServeHTTP] [Vs.GetVoucherByVoucherID: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, voucherView)
	return
}

//CreateVoucher ... 创建voucher时，只要voucher record翻页，就保存，所以此函数插入的voucher record不会太多。
func (vh *VoucherHandlers) CreateVoucher(w http.ResponseWriter, r *http.Request) {
	var params = new(model.VoucherParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucherHandlers/CreateVoucher] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.InfoParams.VoucherMonth == nil || *(params.InfoParams.VoucherMonth) <= 0 {
		ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouMon, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.InfoParams.VoucherMonth == nil || *(params.InfoParams.VoucherMonth) <= 0 {
		ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouMon, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	for _, recParam := range params.RecordsParams {
		if recParam.SubjectName == nil || *recParam.SubjectName == "" {
			ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecSub, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.CreditMoney == nil || *recParam.CreditMoney <= 0.001 {
			ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecCredit, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.DebitMoney == nil || *recParam.DebitMoney <= 0.001 {
			ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecDebit, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
	}
	requestId := vh.GetTraceId(r)

	IdSlice, ccErr := vh.Vs.CreateVoucher(r.Context(), params, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/CreateVoucher/ServeHTTP] [Vs.CreateVoucher: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{len(IdSlice), IdSlice}
	vh.Response(r.Context(), vh.Logger, w, nil, dataBuf)
	return
}

func (vh *VoucherHandlers) DeleteVoucher(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteIDParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucher/DeleteVoucher] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	requestId := vh.GetTraceId(r)
	ccErr = vh.Vs.DeleteVoucher(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/DeleteVoucher/ServeHTTP] [Vs.DeleteVoucher: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, nil)
	return
}

func (vh *VoucherHandlers) CreateVoucherRecords(w http.ResponseWriter, r *http.Request) {
	var recordsParams []model.VoucherRecordParams
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucher/CreateVoucherRecords] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if len(recordsParams) == 0 {
		ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVoucherRecord, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	
	for _, recParam := recordsParams {
		if recParam.VoucherID == nil || *recParam.VoucherID == 0 {
			ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrIds, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.SubjectName == nil || *recParam.SubjectName == "" {
			ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecSub, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.CreditMoney == nil || *recParam.CreditMoney <= 0.001 {
			ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecCredit, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.DebitMoney == nil || *recParam.DebitMoney <= 0.001 {
			ccErr = service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecDebit, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
	}
	requestId := vh.GetTraceId(r)

	IdSlice, ccErr := vh.Vrs.CreateVoucherRecords(r.Context(), recordsParams, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/CreateVoucherRecords/ServeHTTP] [Vrs.CreateVoucherRecords: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{len(IdSlice), IdSlice}
	vh.Response(r.Context(), vh.Logger, w, nil, dataBuf)
	return
}

func (vh *VoucherHandlers) UpdateVoucherRecord(w http.ResponseWriter, r *http.Request) {
	var params = new(model.VoucherRecordParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucher/UpdateVoucherRecord] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.VoucherID == nil || *params.VoucherID <= 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	updateFields := make(map[string]interface{})
	if params.SubjectName != nil {
		if *params.SubjectName == "" {
			ccErr := service.NewError(service.ErrVoucher, service.ErrNotAllowed, service.ErrEmpty, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		updateFields["SubjectName"] = *params.SubjectName
	}
	if params.Summary != nil {		
		updateFields["Summary"] = *params.Summary
	}
	if params.BillCount != nil {
		updateFields["BillCount"] = *params.BillCount
	}
	if params.CreditMoney != nil {
		updateFields["CreditMoney"] = *params.CreditMoney
	}
	if params.DebitMoney != nil {
		updateFields["DebitMoney"] = *params.DebitMoney
	}
	if params.SubID1 != nil {
		updateFields["SubID1"] = *params.SubID1
	}
	if params.SubID2 != nil {
		updateFields["SubID2"] = *params.SubID2
	}
	if params.SubID3 != nil {
		updateFields["SubID3"] = *params.SubID3
	}
	if params.SubID4 != nil {
		updateFields["SubID4"] = *params.SubID4
	}
	if len(updateFields) == 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	ccErr := vh.Vrs.UpdateVoucherRecord(r.Context(), *params.VoucherID, updateFields)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/UpdateVoucherRecord/ServeHTTP] [Vrs.UpdateVoucherRecord: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, nil)
	return
}

func (vh *VoucherHandlers) DeleteVoucherRecord(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteIDParams)
	err := oh.HttpRequestParse(r, params)
	if err != nil {
		oh.Logger.ErrorContext(r.Context(), "[voucher/DeleteVoucherRecord] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	requestId := oh.GetTraceId(r)
	ccErr = vh.Vrs.DeleteVoucherRecordByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[voucher/DeleteVoucherRecord/ServeHTTP] [Vrs.DeleteVoucherRecordByID: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, nil)
	return
}
