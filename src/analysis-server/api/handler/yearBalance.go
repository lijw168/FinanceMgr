package handler

import (
	"analysis-server/api/service"
	//"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"net/http"
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
	requestId := yh.GetTraceId(r)

	yearBal, ccErr := yh.YearBalService.GetYearBalance(r.Context(), *(params.Year), *(params.SubjectID), requestId)
	if ccErr != nil {
		yh.Logger.WarnContext(r.Context(), "[yearBalance/GetYearBalance/ServerHTTP] [YearBalService.GetYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, yearBal)
	return
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
	return
}

func (yh *YearBalHandlers) BatchCreateYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.OptYearBalsParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/BatchCreateYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	// for _, param := range params {
	// 	if param.Year == nil || *(param.Year) <= 0 {
	// 		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
	// 		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
	// 		return
	// 	}
	// 	if param.SubjectID == nil || *(param.SubjectID) <= 0 {
	// 		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
	// 		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
	// 		return
	// 	}
	// 	if param.Balance == nil {
	// 		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrBalance, service.ErrNull)
	// 		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
	// 		return
	// 	}
	// }
	requestId := yh.GetTraceId(r)
	ccErr := yh.YearBalService.BatchCreateYearBalance(r.Context(), params.OptYearBals, requestId)
	yh.Logger.InfoContext(r.Context(), "YearBalService.BatchCreateYearBalance in BatchCreateYearBalance.")
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/BatchCreateYearBalance/ServerHTTP] [YearBalService.BatchCreateYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
	return
}

func (yh *YearBalHandlers) UpdateYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.OptYearBalanceParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.SubjectID == nil || *params.SubjectID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Year == nil || *(params.Year) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	updateFields := make(map[string]interface{})
	if params.Balance != nil {
		updateFields["balance"] = *params.Balance
	}
	if len(updateFields) == 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	ccErr := yh.YearBalService.UpdateYearBalance(r.Context(), *(params.Year), *(params.SubjectID), updateFields)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers/ServerHTTP] [YearBalService.UpdateYearBalanceById: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
	return
}

func (yh *YearBalHandlers) BatchUpdateYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.OptYearBalsParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	// for _, param := range params {
	// 	if param.SubjectID == nil || *param.SubjectID <= 0 {
	// 		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
	// 		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
	// 		return
	// 	}
	// 	if param.Year == nil || *param.Year <= 0 {
	// 		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
	// 		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
	// 		return
	// 	}
	// 	updateFields := make(map[string]interface{})
	// 	if param.Balance != nil {
	// 		updateFields["balance"] = *param.Balance
	// 	} else {
	// 		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
	// 		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
	// 		return
	// 	}
	// }
	requestId := yh.GetTraceId(r)
	ccErr := yh.YearBalService.BatchUpdateYearBalance(r.Context(), params.OptYearBals, requestId)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers/ServerHTTP] [YearBalService.BatchUpdateYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
	return
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
	if params.SubjectID == nil || *params.SubjectID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Year == nil || *(params.Year) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYear, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	requestId := yh.GetTraceId(r)
	ccErr := yh.YearBalService.DeleteYearBalance(r.Context(), *(params.Year), *(params.SubjectID), requestId)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/DeleteYearBalance/ServerHTTP] [YearBalService.DeleteYearBalance: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
	return
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
	if isLackBaseParams([]string{"subjectId", "companyId"}, params.Filter) {
		yh.Logger.ErrorContext(r.Context(), "lack base param  operatorId")
		ce := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
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
		case cons.Order_Asc, cons.Order_Desc:
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
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
	return
}
