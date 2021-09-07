package handler

import (
	"analysis-server/api/service"
	"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"net/http"
	"unicode/utf8"
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

//获取某个公司的某个月份的最大的凭证号，暂不考虑年度
func (vh *VoucherHandlers) GetMaxNumOfMonth(w http.ResponseWriter, r *http.Request) {
	var params = new(model.QueryMaxNumOfMonthParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucherHandlers/QueryMaxNumOfMonth] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucherInfo, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrVoucherInfo, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.VoucherMonth == nil || *(params.VoucherMonth) <= 0 {
		ccErr := service.NewError(service.ErrVoucherInfo, service.ErrMiss, service.ErrVouMon, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	requestId := vh.GetTraceId(r)
	count, ccErr := vh.Vis.GetMaxNumOfMonthByContion(r.Context(), params, requestId)
	if ccErr != nil {
		FunctionName := "voucherHandlers/GetMaxNumOfMonth/ServerHTTP"
		vh.Logger.WarnContext(r.Context(), "[requestId:%s][%s] [Vis.GetMaxNumOfMonthByContion: %s]",
			requestId, FunctionName, ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, count)
	return
}

func (vh *VoucherHandlers) GetLatestVoucherInfo(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucherHandlers/GetLatestVoucherInfo] [HttpRequestParse: %v]", err)
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
	vouInfoViews, count, ccErr := vh.Vis.GetLatestVoucherInfoByCompanyID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		FunctionName := "voucherHandlers/GetVoucherInfo/ServerHTTP"
		vh.Logger.WarnContext(r.Context(), "[requestId:%s][%s] [Vis.GetVoucherInfoByID: %s]", requestId, FunctionName, ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), vouInfoViews}
	vh.Response(r.Context(), vh.Logger, w, nil, dataBuf)
	return
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
	if isLackBaseParams([]string{"voucherId", "companyId"}, params.Filter) {
		vh.Logger.ErrorContext(r.Context(), "lack base param  Id")
		ce := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ce, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		//先暂时修改为一个值，如果以后确实需要，再进行添加。
		//filterMap["voucherId"] = utils.Attribute{Type: utils.T_Int_Arr, Val: nil}
		filterMap["voucherId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["companyId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		//filterMap["voucherMonth"] = utils.Attribute{Type: utils.T_Int_Arr, Val: nil}
		filterMap["voucherMonth"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["numOfMonth"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["voucherDate"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["voucherFiller"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["voucherAuditor"] = utils.Attribute{Type: utils.T_String, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrField, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
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

func (vh *VoucherHandlers) ListVoucherInfoByMulCondition(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListVoucherInfoParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucherInfo/ListVoucherInfoByMulCondition] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if isLackBaseParams([]string{"voucherId", "companyId"}, params.BasicFilter) {
		vh.Logger.ErrorContext(r.Context(), "lack base param Id")
		ce := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ce, nil)
		return
	}
	if params.BasicFilter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["voucherId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["companyId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["voucherMonth"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		//因为numOfMonth和 voucherDate可能是一个值，也可能是多个值
		//filterMap["numOfMonth"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		//filterMap["voucherDate"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["voucherFiller"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["voucherAuditor"] = utils.Attribute{Type: utils.T_String, Val: nil}
		if !utils.ValiFilter(filterMap, params.BasicFilter) {
			ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrField, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
	}
	if params.AuxiFilter != nil {
		filterMap := map[string]utils.Attribute{}
		//因为credit和 debit可能是一个值，也可能是多个值,所以此处就不用做检查了。
		filterMap["subjectName"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["summary"] = utils.Attribute{Type: utils.T_String, Val: nil}
		if !utils.ValiFilter(filterMap, params.AuxiFilter) {
			ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrField, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
	}
	if (params.Order != nil) && (len(params.Order) > 0) {
		switch *params.Order[0].Field {
		case "voucherId":
			*params.Order[0].Field = "voucherId"
		case "createdAt":
			*params.Order[0].Field = "createdAt"
		case "updatedAt":
			*params.Order[0].Field = "updatedAt"
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
	vouInfoViews, count, ccErr := vh.Vs.ListVoucherInfoByMulCondition(r.Context(), params)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucherHandlers/ListVoucherInfoByMulCondition/ServerHTTP] [Error: %s]", ccErr.Detail())
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
	if isLackBaseParams([]string{"voucherId", "recordId"}, params.Filter) {
		vh.Logger.ErrorContext(r.Context(), "lack base param  operatorId")
		ce := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ce, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["voucherId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		//先暂时修改为一个值，如果以后确实需要，再进行添加。
		//filterMap["recordId"] = utils.Attribute{Type: utils.T_Int_Arr, Val: nil}
		filterMap["recordId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["status"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["subjectName"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["summary"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["subId1"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["subId2"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["subId3"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["subId4"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrField, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
	}
	if (params.Order != nil) && (len(params.Order) > 0) {
		switch *params.Order[0].Field {
		case "createdAt":
			*params.Order[0].Field = "createdAt"
		case "updatedAt":
			*params.Order[0].Field = "updatedAt"
		case "recordId":
			*params.Order[0].Field = "recordId"
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
			vh.Response(r.Context(), vh.Logger, w, ce, nil)
			return
		}
		if params.Order[0].Direction != nil {
			switch *params.Order[0].Direction {
			case cons.Order_Asc, cons.Order_Desc:
			default:
				ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
				vh.Response(r.Context(), vh.Logger, w, ce, nil)
				return
			}
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
	var params = new(model.CreateVoucherParams)
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
	//由于制证日期和制证月份是一一对应的，所以可以不判断制证月份了。因为如果不传该参数，就默认是当前日期和月份。
	// if params.InfoParams.VoucherMonth == nil || *(params.InfoParams.VoucherMonth) <= 0 {
	// 	ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouMon, service.ErrNull)
	// 	vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
	// 	return
	// }
	if params.InfoParams.VoucherFiller == nil || *(params.InfoParams.VoucherFiller) == "" {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouFiller, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if utf8.RuneCountInString(*(params.InfoParams.VoucherFiller)) > NameMaxLen ||
		!utils.VerStrP(*(params.InfoParams.VoucherFiller)) {
		ccErr := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrVouFiller, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}

	for _, recParam := range params.RecordsParams {
		if recParam.Summary == nil || *recParam.Summary == "" {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouSummary, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.SubjectName == nil || *recParam.SubjectName == "" {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecSub, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.CreditMoney == nil && recParam.DebitMoney == nil {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVoucherData, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if *recParam.CreditMoney <= 0.001 && *recParam.DebitMoney <= 0.001 {
			ccErr := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrParam, service.ErrNull)
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

//update voucher,include:voucherInfo,voucherRecord
func (vh *VoucherHandlers) UpdateVoucher(w http.ResponseWriter, r *http.Request) {
	var params = new(model.UpdateVoucherParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucherHandlers/UpdateVoucher] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.ModifyInfoParams != nil {
		if params.ModifyInfoParams.VoucherID == nil {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
	}
	requestId := vh.GetTraceId(r)
	IdSlice, ccErr := vh.Vs.UpdateVoucher(r.Context(), params, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/UpdateVoucher/ServerHTTP] [Vrs.UpdateVoucher: %s]",
			ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if len(IdSlice) != 0 {
		dataBuf := &DescData{int64(len(IdSlice)), IdSlice}
		vh.Response(r.Context(), vh.Logger, w, nil, dataBuf)
	} else {
		vh.Response(r.Context(), vh.Logger, w, nil, nil)
	}
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
	err := vh.HttpRequestParse(r, &recordsParams)
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
		if recParam.Summary == nil || *recParam.Summary == "" {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouSummary, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.SubjectName == nil || *recParam.SubjectName == "" {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouRecSub, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if recParam.CreditMoney == nil && recParam.DebitMoney == nil {
			ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVoucherData, service.ErrNull)
			vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
			return
		}
		if *recParam.CreditMoney <= 0.001 && *recParam.DebitMoney <= 0.001 {
			ccErr := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrParam, service.ErrNull)
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

func (vh *VoucherHandlers) UpdateVoucherRecordByID(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ModifyVoucherRecordParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucher/UpdateVoucherRecordByID] [HttpRequestParse: %v]", err)
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
		updateFields["subjectName"] = *params.SubjectName
	}
	if params.Summary != nil {
		updateFields["summary"] = *params.Summary
	}
	if params.BillCount != nil {
		updateFields["billCount"] = *params.BillCount
	}
	if params.CreditMoney != nil {
		updateFields["creditMoney"] = *params.CreditMoney
	}
	if params.DebitMoney != nil {
		updateFields["debitMoney"] = *params.DebitMoney
	}
	if params.SubID1 != nil {
		updateFields["subId1"] = *params.SubID1
	}
	if params.SubID2 != nil {
		updateFields["subId2"] = *params.SubID2
	}
	if params.SubID3 != nil {
		updateFields["subId3"] = *params.SubID3
	}
	if params.SubID4 != nil {
		updateFields["subId4"] = *params.SubID4
	}
	if len(updateFields) == 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	ccErr := vh.Vrs.UpdateVoucherRecordByID(r.Context(), *params.VouRecordID, updateFields)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/UpdateVoucherRecordByID/ServerHTTP] [Vrs.UpdateVoucherRecordByID: %s]", ccErr.Detail())
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

func (vh *VoucherHandlers) DeleteVoucherRecords(w http.ResponseWriter, r *http.Request) {
	var params = new(model.IDsParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucher/DeleteVoucherRecords] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucherRecord, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if params.IDs == nil || len(params.IDs) == 0 {
		ccErr := service.NewError(service.ErrVoucherRecord, service.ErrMiss, service.ErrIds, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	requestId := vh.GetTraceId(r)
	ccErr := vh.Vrs.DeleteVoucherRecords(r.Context(), params, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/DeleteVoucherRecords/ServerHTTP] [Vrs.DeleteVoucherRecords: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, nil)
	return
}

func (vh *VoucherHandlers) VoucherAudit(w http.ResponseWriter, r *http.Request) {
	var params = new(model.VoucherAuditParams)
	err := vh.HttpRequestParse(r, params)
	if err != nil {
		vh.Logger.ErrorContext(r.Context(), "[voucher/DeleteVoucherRecord] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrVoucher, service.ErrMalformed, service.ErrNull, err.Error())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	updateFields := make(map[string]interface{})
	if params.VoucherID == nil || *params.VoucherID <= 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrId, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	updateFields["voucherId"] = *params.VoucherID
	if params.VoucherAuditor == nil /*|| *(params.VoucherAuditor) == "" */ {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrVouAuditor, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	if utf8.RuneCountInString(*(params.VoucherAuditor)) > NameMaxLen || !utils.VerStrP(*(params.VoucherAuditor)) {
		ccErr := service.NewError(service.ErrVoucher, service.ErrInvalid, service.ErrVouAuditor, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	updateFields["voucherAuditor"] = *params.VoucherAuditor
	if params.Status == nil || *params.Status <= 0 {
		ccErr := service.NewError(service.ErrVoucher, service.ErrMiss, service.ErrStatus, service.ErrNull)
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	updateFields["status"] = *params.Status
	requestId := vh.GetTraceId(r)
	ccErr := vh.Vs.VoucherAudit(r.Context(), *params.VoucherID, updateFields, requestId)
	if ccErr != nil {
		vh.Logger.WarnContext(r.Context(), "[voucher/DeleteVoucherRecord/ServerHTTP] [Vrs.DeleteVoucherRecordByID: %s]", ccErr.Detail())
		vh.Response(r.Context(), vh.Logger, w, ccErr, nil)
		return
	}
	vh.Response(r.Context(), vh.Logger, w, nil, nil)
	return
}
