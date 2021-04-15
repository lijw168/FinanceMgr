package handler

import (
	"analysis-server/api/service"
	"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"net/http"
)

// var (
// 	DescriptionMaxlen = 256
// 	NameMaxLen        = 32
// )

type VoucherHandlers struct {
	CCHandler
	Logger *log.Logger
	Vis    *service.VoucherInfoService
	Vrs    *service.VoucherRecordService
	Vs     *service.VoucherService
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
		filterMap["voucher_id"] = utils.Attribute{Type: utils.T_Int_Arr, Val: nil}
		filterMap["company_id"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["voucher_month"] = utils.Attribute{Type: utils.T_Int_Arr, Val: nil}
		filterMap["num_of_month"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["voucher_date"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrField, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
	}
	if (params.Order != nil) && (len(params.Order) > 0) {
		switch *params.Order[0].Field {
		case "create_time":
			*params.Order[0].Field = "created_at"
		case "update_time":
			*params.Order[0].Field = "updated_at"
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
		switch *params.Order[0].Direction {
		case cons.Order_Asc, cons.Order_Desc:
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
	}
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
		vh.Logger.WarnContext(r.Context(), "[voucherInfo/ListVoucherInfo/ServerHTTP] [VoucherInfoService.ListVoucherInfo: %s]", ccErr.Detail())
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
		filterMap["voucher_id"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["record_id"] = utils.Attribute{Type: utils.T_Int_Arr, Val: nil}
		filterMap["subject_name"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["summary"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["sub_id1"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["sub_id2"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["sub_id3"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["sub_id4"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrField, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
	}
	if (params.Order != nil) && (len(params.Order) > 0) {
		switch *params.Order[0].Field {
		case "create_time":
			*params.Order[0].Field = "created_at"
		case "update_time":
			*params.Order[0].Field = "updated_at"
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
		switch *params.Order[0].Direction {
		case cons.Order_Asc, cons.Order_Desc:
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
	}
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
		vh.Logger.WarnContext(r.Context(), "[voucherInfo/ListVoucherRecords/ServerHTTP] [VoucherRecordService.ListVoucherRecords: %s]", ccErr.Detail())
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
		vh.Logger.WarnContext(r.Context(), "[voucherHandlers/GetVoucherInfo/ServerHTTP] [Vis.GetVoucherInfoByID: %s]", ccErr.Detail())
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
		vh.Logger.WarnContext(r.Context(), "[voucherHandlers/GetVoucher/ServerHTTP] [Vs.GetVoucherByVoucherID: %s]", ccErr.Detail())
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
	if params.InfoParams.CompanyID == nil || *(params.InfoParams.CompanyID) <= 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.InfoParams.VoucherMonth == nil || *(params.InfoParams.VoucherMonth) <= 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouMon, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	for _, recParam := range params.RecordsParams {
		if recParam.SubjectName == nil || *recParam.SubjectName == "" {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecSub, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.CreditMoney == nil || *recParam.CreditMoney <= 0.001 {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecCredit, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.DebitMoney == nil || *recParam.DebitMoney <= 0.001 {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecDebit, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
	}
	requestId := vh.GetTraceId(r)

	IdSlice, ccErr := vh.Vs.CreateVoucher(r.Context(), params, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/CreateVoucher/ServerHTTP] [Vs.CreateVoucher: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{int64(len(IdSlice)), IdSlice}
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
	ccErr := vh.Vs.DeleteVoucher(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/DeleteVoucher/ServerHTTP] [Vs.DeleteVoucher: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, nil)
	return
}

func (vh *VoucherHandlers) CreateVoucherRecords(w http.ResponseWriter, r *http.Request) {
	var recordsParams []*model.CreateVoucherRecordParams
	err := vh.HttpRequestParse(r, recordsParams)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucher/CreateVoucherRecords] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if len(recordsParams) == 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVoucherRecord, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}

	for _, recParam := range recordsParams {
		if recParam.VoucherID == nil || *recParam.VoucherID == 0 {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrIds, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.SubjectName == nil || *recParam.SubjectName == "" {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecSub, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.CreditMoney == nil || *recParam.CreditMoney <= 0.001 {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecCredit, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.DebitMoney == nil || *recParam.DebitMoney <= 0.001 {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecDebit, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
	}
	requestId := vh.GetTraceId(r)

	IdSlice, ccErr := vh.Vrs.CreateVoucherRecords(r.Context(), recordsParams, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/CreateVoucherRecords/ServerHTTP] [Vrs.CreateVoucherRecords: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{int64(len(IdSlice)), IdSlice}
	vh.Response(r.Context(), vh.Logger, w, nil, dataBuf)
	return
}

func (vh *VoucherHandlers) UpdateVoucherRecord(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ModifyVoucherRecordParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucher/UpdateVoucherRecord] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.VouRecordID == nil || *params.VouRecordID <= 0 {
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
	ccErr := vh.Vrs.UpdateVoucherRecord(r.Context(), *params.VouRecordID, updateFields)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/UpdateVoucherRecord/ServerHTTP] [Vrs.UpdateVoucherRecord: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, nil)
	return
}

func (vh *VoucherHandlers) DeleteVoucherRecord(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteIDParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucher/DeleteVoucherRecord] [HttpRequestParse: %v]", err)
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
	ccErr := vh.Vrs.DeleteVoucherRecordByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/DeleteVoucherRecord/ServerHTTP] [Vrs.DeleteVoucherRecordByID: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, nil)
	return
}
