package handler

import (
	"financeMgr/src/analysis-server/api/service"
	"financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	"financeMgr/src/common/log"
	"net/http"
	"strconv"
)

type YearBalHandlers struct {
	CCHandler
	Logger         *log.Logger
	YearBalService *service.YearBalanceService
}

func (yh *YearBalHandlers) GetYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.BasicYearBalanceParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/GetYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if ccErr := yh.checkoutYearBalanceBaseParams(params); ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/GetYearBalance] [checkoutYearBalanceBaseParams: %v]", err)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	requestId := yh.GetTraceId(r)

	yearBal, ccErr := yh.YearBalService.GetYearBalance(r.Context(), params, requestId)
	if ccErr != nil {
		yh.Logger.WarnContext(r.Context(), "[yearBalance/GetYearBalance/ServerHTTP] [YearBalService.GetYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, yearBal)
}

func (yh *YearBalHandlers) GetAccSubYearBalValue(w http.ResponseWriter, r *http.Request) {
	var params = new(model.BasicYearBalanceParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/GetAccSubYearBalValue] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if ccErr := yh.checkoutYearBalanceBaseParams(params); ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/GetAccSubYearBalValue] [checkoutYearBalanceBaseParams: %v]", err)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	requestId := yh.GetTraceId(r)

	dBalanceValue, ccErr := yh.YearBalService.GetAccSubYearBalValue(r.Context(), params, requestId)
	if ccErr != nil {
		yh.Logger.WarnContext(r.Context(), "[yearBalance/GetAccSubYearBalValue/ServerHTTP] [YearBalService.GetAccSubYearBalValue: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, dBalanceValue)
}

func (yh *YearBalHandlers) CreateYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.OptYearBalanceParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/CreateYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Year == nil || *(params.Year) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.SubjectID == nil || *(params.SubjectID) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Balance == nil {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrBalance, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	requestId := yh.GetTraceId(r)
	ccErr := yh.YearBalService.CreateYearBalance(r.Context(), params, requestId)
	yh.Logger.InfoContext(r.Context(), "YearBalService.CreateYearBalance in CreateYearBalance.")
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/CreateYearBalance/ServerHTTP] [YearBalService.CreateYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
}

func (yh *YearBalHandlers) BatchCreateYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.BatchCreateYearBalsParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/BatchCreateYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	//data check
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Year == nil || *(params.Year) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if len(params.OptSubAndBals) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrIds, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	ccErr := yh.YearBalService.BatchCreateYearBalance(r.Context(), params)
	yh.Logger.InfoContext(r.Context(), "YearBalService.BatchCreateYearBalance in BatchCreateYearBalance.")
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/BatchCreateYearBalance/ServerHTTP] [YearBalService.BatchCreateYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
}

// 该函数仅仅批量更新balance这一个字段。
func (yh *YearBalHandlers) BatchUpdateBals(w http.ResponseWriter, r *http.Request) {
	var params = new(model.BatchUpdateBalsParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Year == nil || *(params.Year) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}

	ccErr := yh.YearBalService.BatchUpdateBals(r.Context(), params)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers/ServerHTTP] [YearBalService.UpdateYearBalanceById: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
}

// 可以更新status和balance这两个字段。
func (yh *YearBalHandlers) UpdateYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.OptYearBalanceParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/UpdateYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	filter := map[string]interface{}{}
	updateField := map[string]interface{}{}
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	filter["companyId"] = *params.CompanyID

	if params.Year == nil || *(params.Year) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	filter["year"] = *params.Year

	if params.SubjectID == nil || *(params.SubjectID) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	filter["subjectId"] = *params.SubjectID

	if params.Balance != nil && *params.Balance > 0 {
		updateField["balance"] = *params.Balance
	}
	if params.Status != nil {
		switch *params.Status {
		case utils.NoAnnualClosing, utils.AnnualClosing:
			updateField["status"] = *params.Status
		}
	}
	if len(updateField) == 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrUpdateParam, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	ccErr := yh.YearBalService.UpdateYearBalance(r.Context(), filter, updateField)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/UpdateYearBalance/ServerHTTP] [YearBalService.UpdateYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
}

func (yh *YearBalHandlers) BatchDeleteYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.BatchDelYearBalsParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if len(params.SubjectIDs) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrIds, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Year == nil || *(params.Year) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	filterFields := make(map[string]interface{})
	filterFields["companyId"] = *params.CompanyID
	filterFields["year"] = *params.Year
	filterFields["subject_id"] = params.SubjectIDs
	ccErr := yh.YearBalService.BatchDeleteYearBalance(r.Context(), filterFields)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers/ServerHTTP] [YearBalService.BatchDeleteYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
}

func (yh *YearBalHandlers) DeleteYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.BasicYearBalanceParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/DeleteYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if ccErr := yh.checkoutYearBalanceBaseParams(params); ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/DeleteYearBalance] [checkoutYearBalanceBaseParams: %v]", err)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}

	requestId := yh.GetTraceId(r)
	ccErr := yh.YearBalService.DeleteYearBalance(r.Context(), params, requestId)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/DeleteYearBalance/ServerHTTP] [YearBalService.DeleteYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
}

func (yh *YearBalHandlers) ListYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[accSub/ListYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if isLackBaseParams([]string{"subjectId", "companyId", "year"}, params.Filter) {
		yh.Logger.ErrorContext(r.Context(), "lack base param")
		ce := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrBaseParam, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ce, nil)
		return
	}
	// if params.Filter != nil {
	// 	filterMap := map[string]utils.Attribute{}
	// 	//对通过c++库的jsoncpp 编码成的整型数组进行验证，有问题。所有此处不再验证该参数。
	//问题是：//从客户端发过来的整型数据，go解析json时，会解析成float64
	// 	filterMap["subjectId"] = utils.Attribute{Type: utils.T_Int_Arr, Val: nil}
	// 	filterMap["year"] = utils.Attribute{Type: utils.T_Int, Val: nil}
	// 	if !utils.ValiFilter(filterMap, params.Filter) {
	// 		ce := service.NewError(service.ErrYearBalance, service.ErrValue, service.ErrField, service.ErrNull)
	// 		yh.Response(r.Context(), yh.Logger, w, ce, nil)
	// 		return
	// 	}
	// }
	if (params.Order != nil) && (len(params.Order) > 0) {
		switch *params.Order[0].Field {
		case "subjectId":
			*params.Order[0].Field = "subjectId"
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
			yh.Response(r.Context(), yh.Logger, w, ce, nil)
			return
		}
		switch *params.Order[0].Direction {
		case utils.OrderAsc, utils.OrderDesc:
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, strconv.Itoa(*params.Order[0].Direction))
			yh.Response(r.Context(), yh.Logger, w, ce, nil)
			return
		}
	}
	yearBalViews, count, ccErr := yh.YearBalService.ListYearBalance(r.Context(), params)
	if ccErr != nil {
		yh.Logger.WarnContext(r.Context(), "[accSub/ListYearBalance/ServerHTTP] [YearBalService.ListYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), yearBalViews}
	yh.Response(r.Context(), yh.Logger, w, nil, dataBuf)
}

// 年度结算
func (yh *YearBalHandlers) AnnualClosing(w http.ResponseWriter, r *http.Request) {
	var params = new(model.BatchCreateYearBalsParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/AnnualClosing] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	//data check
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Year == nil || *(params.Year) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if len(params.OptSubAndBals) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrBaseParam, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	ccErr := yh.YearBalService.AnnualClosing(r.Context(), params)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/AnnualClosing/ServerHTTP] [YearBalService.AnnualClosing: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
}

// 取消年度结算
func (yh *YearBalHandlers) CancelAnnualClosing(w http.ResponseWriter, r *http.Request) {
	var params = new(model.BatchDelYearBalsParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if len(params.SubjectIDs) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrIds, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Year == nil || *(params.Year) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	ccErr := yh.YearBalService.CancelAnnualClosing(r.Context(), params)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers/ServerHTTP] [YearBalService.BatchDeleteYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
}

func (yh *YearBalHandlers) GetAnnualClosingStatus(w http.ResponseWriter, r *http.Request) {
	var params = new(model.BasicYearBalanceParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/GetAnnualClosingStatus] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Year == nil || *params.Year <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}

	iStatus, ccErr := yh.YearBalService.GetAnnualClosingStatus(r.Context(), *params.CompanyID, *params.Year)
	if ccErr != nil {
		yh.Logger.WarnContext(r.Context(), "[yearBalance/GetAnnualClosingStatus/ServerHTTP] [YearBalService.GetAnnualClosingStatus: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, iStatus)
}

func (yh *YearBalHandlers) checkoutYearBalanceBaseParams(params *model.BasicYearBalanceParams) service.CcError {
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		return service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrCompanyId, service.ErrNull)
	}
	if params.SubjectID == nil || *params.SubjectID <= 0 {
		return service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
	}
	if params.Year == nil || *(params.Year) <= 0 {
		return service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
	}
	return nil
}
