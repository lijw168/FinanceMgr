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

var (
	DescriptionMaxlen = 256
	NameMaxLen        = 32
)

type OperatorInfoHandlers struct {
	CCHandler
	Logger         *log.Logger
	OptInfoService *service.OperatorInfoService
	ComService     *service.CompanyService
}

func (oh *OperatorInfoHandlers) ListOperatorInfo(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListOperatorsParams)
	err := oh.HttpRequestParse(r, params)
	if err != nil {
		oh.Logger.ErrorContext(r.Context(), "[operatorInfo/ListOperatorInfo] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrOperator, service.ErrMalformed, service.ErrNull, err.Error())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["company_id"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["job"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["name"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["department"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["status"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["role"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrOperator, service.ErrInvalid, service.ErrField, service.ErrNull)
			oh.Response(r.Context(), oh.Logger, w, ce, nil)
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
			oh.Response(r.Context(), oh.Logger, w, ce, nil)
			return
		}
		switch *params.Order[0].Direction {
		case cons.Order_Asc, cons.Order_Desc:
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
			oh.Response(r.Context(), oh.Logger, w, ce, nil)
			return
		}
	}
	if (params.DescOffset != nil) && (*params.DescOffset < 0) {
		ce := service.NewError(service.ErrOperator, service.ErrInvalid, service.ErrOffset, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ce, nil)
		return
	}
	if (params.DescLimit != nil) && (*params.DescLimit < -1) {
		ce := service.NewError(service.ErrOperator, service.ErrInvalid, service.ErrLimit, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ce, nil)
		return
	}

	optViews, count, ccErr := oh.OptInfoService.ListOperators(r.Context(), params)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[operatorInfo/ListOperatorInfo/ServerHTTP] [OptInfoService.DescribeOperators: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), optViews}
	oh.Response(r.Context(), oh.Logger, w, nil, dataBuf)
	return
}

func (oh *OperatorInfoHandlers) GetOperatorInfo(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeNameParams)
	err := oh.HttpRequestParse(r, params)
	if err != nil {
		oh.Logger.ErrorContext(r.Context(), "[operatorInfo/GetOperatorInfo] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrOperator, service.ErrMalformed, service.ErrNull, err.Error())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}

	if params.Name == nil || *params.Name == "" {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrName, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	requestId := oh.GetTraceId(r)
	optView, ccErr := oh.OptInfoService.GetOperatorInfoByName(r.Context(), *params.Name, requestId)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[operatorInfo/GetOperatorInfo/ServerHTTP] [OptInfoService.GetOperatorInfoByName: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, optView)
	return
}

func (oh *OperatorInfoHandlers) CreateOperator(w http.ResponseWriter, r *http.Request) {
	var params = new(model.OperatorInfoParams)
	err := oh.HttpRequestParse(r, params)
	if err != nil {
		oh.Logger.ErrorContext(r.Context(), "[operatorInfo/CreateOperator] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrOperator, service.ErrMalformed, service.ErrNull, err.Error())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil || *params.CompanyID <= 0 {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrId, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.Name == nil || *params.Name == "" {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrName, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.Name == nil || utf8.RuneCountInString(*params.Name) > NameMaxLen || !utils.VerStrP(*params.Name) {
		ccErr := service.NewError(service.ErrOperator, service.ErrInvalid, service.ErrName, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if *params.Password == "" {
		ccErr := service.NewError(service.ErrOperator, service.ErrNotAllowed, service.ErrEmpty, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if utf8.RuneCountInString(*params.Password) > NameMaxLen || !utils.VerStrP(*params.Password) {
		ccErr := service.NewError(service.ErrOperator, service.ErrInvalid, service.ErrPasswd, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	requestId := oh.GetTraceId(r)

	optInfoView, ccErr := oh.OptInfoService.CreateOptInfo(r.Context(), params, requestId)
	oh.Logger.InfoContext(r.Context(), "CreateOptInfo in CreateOperator.")
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[operatorInfo/CreateOperator/ServerHTTP] [OptInfoService.CreateOptInfo: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, optInfoView)
	return
}

func (oh *OperatorInfoHandlers) UpdateOperator(w http.ResponseWriter, r *http.Request) {
	var params = new(model.OperatorInfoParams)
	err := oh.HttpRequestParse(r, params)
	if err != nil {
		oh.Logger.ErrorContext(r.Context(), "[volumes/UpdateOperator] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrOperator, service.ErrMalformed, service.ErrNull, err.Error())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.Name == nil || *params.Name == "" {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrName, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	updateFields := make(map[string]interface{})
	if params.Password != nil {
		if *params.Password == "" {
			ccErr := service.NewError(service.ErrOperator, service.ErrNotAllowed, service.ErrEmpty, service.ErrNull)
			oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
			return
		}
		if utf8.RuneCountInString(*params.Password) > NameMaxLen || !utils.VerStrP(*params.Password) {
			ccErr := service.NewError(service.ErrOperator, service.ErrInvalid, service.ErrPasswd, service.ErrNull)
			oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
			return
		}
		updateFields["Password"] = *params.Password
	}
	if params.CompanyID != nil {
		//要判断companyID是否存在, 这步操作只有在用户登录时，才有用。
		// requestId := oh.GetTraceId(r)
		// comView, ccErr := oh.ComService.GetCompanyById(r.Context(), *params.CompanyID, requestId)
		// if comView == nil || ccErr != nil {
		// 	oh.Logger.WarnContext(r.Context(), "[opreator/UpdateOperator/ServerHTTP] [ComService.GetCompanyById: %s]", ccErr.Detail())
		// 	oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		// 	return
		// }
		updateFields["CompanyId"] = *params.CompanyID
	}
	if params.Job != nil {
		updateFields["Job"] = *params.Job
	}
	if params.Department != nil {
		updateFields["Department"] = *params.Department
	}
	if params.Status != nil {
		updateFields["Status"] = *params.Status
	}
	if params.Role != nil {
		updateFields["Role"] = *params.Role
	}
	if len(updateFields) == 0 {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	ccErr := oh.OptInfoService.UpdateOperator(r.Context(), *params.Name, updateFields)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[opreator/UpdateOperator/ServerHTTP] [OptInfoService.UpdateOperator: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, nil)
	return
}

func (oh *OperatorInfoHandlers) DeleteOperator(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteOperatorParams)
	err := oh.HttpRequestParse(r, params)
	if err != nil {
		oh.Logger.ErrorContext(r.Context(), "[opreator/DeleteOperator] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrOperator, service.ErrMalformed, service.ErrNull, err.Error())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.Name == nil || *params.Name == "" {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrName, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	requestId := oh.GetTraceId(r)
	ccErr := oh.OptInfoService.DeleteOperatorInfoByName(r.Context(), *params.Name, requestId)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[opreator/DeleteOperator/ServerHTTP] [OptInfoService.DeleteOperatorInfoByName: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, nil)
	return
}
