package handler

import (
	"net/http"
	"unicode/utf8"

	cons "common/constant"
	"common/log"
	"analysis-server/api/service"
	"analysis-server/api/utils"
	"analysis-server/model"
)

var (
	DescriptionMaxlen = 256
	NameMaxLen        = 32
)

type OperatorInfoHandlers struct {
	CCHandler
	Logger         *log.Logger
	OptInfoService service.OperatorInfoService
	ComService     service.CompanyService
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
		filterMap["companyId"] = utils.Attribute{utils.T_String, nil}
		filterMap["job"] = utils.Attribute{utils.T_String, nil}
		filterMap["name"] = utils.Attribute{utils.T_String, nil}
		filterMap["department"] = utils.Attribute{utils.T_String, nil}
		filterMap["status"] = utils.Attribute{utils.T_Int, nil}
		filterMap["role"] = utils.Attribute{utils.T_Int, nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrOperator, service.ErrDesc, service.ErrField, service.ErrNull)
			oh.Response(r.Context(), oh.Logger, w, ce, nil)
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
	// 		oh.Response(r.Context(), oh.Logger, w, ce, nil)
	// 		return
	// 	}
	// 	switch *params.Order[0].Direction {
	// 	case cons.Order_Asc, cons.Order_Desc:
	// 	default:
	// 		ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
	// 		oh.Response(r.Context(), oh.Logger, w, ce, nil)
	// 		return
	// 	}
	// }
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

	optViews, count, ccErr := oh.OptInfoService.DescribeOperators(r.Context(), params)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[operatorInfo/ListOperatorInfo/ServeHTTP] [OptInfoService.DescribeOperators: %s]", ccErr.Detail())
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
	requestId := ah.GetTraceId(r)
	optView, ccErr := oh.OptInfoService.GetOperatorInfoByName(r.Context(), params, requestId)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[operatorInfo/GetOperatorInfo/ServeHTTP] [OptInfoService.GetOperatorInfoByName: %s]", ccErr.Detail())
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

	if params.Size == nil || !(*params.Size >= 10 && *params.Size <= 16000 && *params.Size%10 == 0) {
		ccErr = service.NewCcError(cons.CodeVolInvalSize, service.ErrOperator, service.ErrInvalid, service.ErrSize, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.VolumeTypeName == nil || *params.VolumeTypeName == "" {
		ccErr = service.NewError(service.ErrOperator, service.ErrMiss, service.ErrType, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}

	requestId := oh.GetTraceId(r)

	volumeView, ccErr := oh.OptInfoService.CreateOptInfo(r.Context(), params, requestId, "")
	oh.Logger.InfoContext(r.Context(), "CreateOptInfo in CreateOperator.")
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[operatorInfo/CreateOperator/ServeHTTP] [OptInfoService.CreateOptInfo: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, volumeView)
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
		updateFields["Password"] = *params.CompanyIDPassword
	}
	if params.CompanyID != nil {
		//要判断companyID是否存在
		requestId := oh.GetTraceId(r)
		comView, ccErr := oh.ComService.GetCompanyById(r.Context(), *params.CompanyID, requestId)
		if comView == nil || ccErr != nil {
			oh.Logger.WarnContext(r.Context(), "[opreator/UpdateOperator/ServeHTTP] [ComService.GetCompanyById: %s]", ccErr.Detail())
			oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
			return
		}
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
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ccErr := oh.OptInfoService.UpdateOperator(r.Context(), *params.Name, updateFields)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[opreator/UpdateOperator/ServeHTTP] [OptInfoService.UpdateOperator: %s]", ccErr.Detail())
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
	ccErr = oh.OptInfoService.DeleteOperatorInfoByName(r.Context(), params, requestId)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[opreator/DeleteOperator/ServeHTTP] [OptInfoService.DeleteOperatorInfoByName: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, nil)
	return
}
