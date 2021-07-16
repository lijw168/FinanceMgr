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
	if isLackBaseParams([]string{"operatorId", "companyId"}, params.Filter) {
		oh.Logger.ErrorContext(r.Context(), "lack base param  operatorId")
		ce := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrField, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ce, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["operatorId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["companyId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
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
		case "createdAt":
			*params.Order[0].Field = "createdAt"
		case "updatedAt":
			*params.Order[0].Field = "updatedAt"
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
		oh.Logger.ErrorContext(r.Context(), "[operatorInfo/ListOperatorInfo/ServerHTTP] [OptInfoService.ListOperators: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), optViews}
	oh.Response(r.Context(), oh.Logger, w, nil, dataBuf)
	return
}

func (oh *OperatorInfoHandlers) GetOperatorInfo(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := oh.HttpRequestParse(r, params)
	if err != nil {
		oh.Logger.ErrorContext(r.Context(), "[operatorInfo/GetOperatorInfo] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrOperator, service.ErrMalformed, service.ErrNull, err.Error())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}

	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrId, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}

	requestId := oh.GetTraceId(r)
	optView, ccErr := oh.OptInfoService.GetOperatorInfoByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[operatorInfo/GetOperatorInfo/ServerHTTP] [OptInfoService.GetOperatorInfoByID: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, optView)
	return
}

func (oh *OperatorInfoHandlers) isCreateAdmin(iRole int) bool {
	if (0xF0 & iRole) > 0 {
		return true
	}
	return false
}
func (oh *OperatorInfoHandlers) CreateOperator(w http.ResponseWriter, r *http.Request) {
	var params = new(model.CreateOptInfoParams)
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
	if params.Role == nil || *params.Role <= 0 {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrRole, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if oh.isCreateAdmin(*params.Role) {
		if !GAccessTokenH.isRootRequest(r) {
			ccErr := service.NewError(service.ErrOperator, service.ErrUnsupported, service.ErrNull, service.ErrNoAuthority)
			oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
			return
		}
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
	var params = new(model.ModifyOptInfoParams)
	err := oh.HttpRequestParse(r, params)
	if err != nil {
		oh.Logger.ErrorContext(r.Context(), "[operator/UpdateOperator] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrOperator, service.ErrMalformed, service.ErrNull, err.Error())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.OperatorID == nil || *params.OperatorID <= 0 {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrId, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	updateFields := make(map[string]interface{})
	if params.Name != nil {
		if *params.Name == "" {
			ccErr := service.NewError(service.ErrOperator, service.ErrNotAllowed, service.ErrEmpty, service.ErrNull)
			oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
			return
		}
		if utf8.RuneCountInString(*params.Name) > NameMaxLen || !utils.VerStrP(*params.Name) {
			ccErr := service.NewError(service.ErrOperator, service.ErrInvalid, service.ErrName, service.ErrNull)
			oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
			return
		}
		updateFields["name"] = *params.Name
	}
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
	ccErr := oh.OptInfoService.UpdateOperator(r.Context(), *params.OperatorID, updateFields)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[opreator/UpdateOperator/ServerHTTP] [OptInfoService.UpdateOperator: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, nil)
	return
}

func (oh *OperatorInfoHandlers) DeleteOperator(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteIDParams)
	err := oh.HttpRequestParse(r, params)
	if err != nil {
		oh.Logger.ErrorContext(r.Context(), "[opreator/DeleteOperator] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrOperator, service.ErrMalformed, service.ErrNull, err.Error())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrOperator, service.ErrMiss, service.ErrId, service.ErrNull)
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	requestId := oh.GetTraceId(r)
	ccErr := oh.OptInfoService.DeleteOperatorInfoByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		oh.Logger.WarnContext(r.Context(), "[opreator/DeleteOperator/ServerHTTP] [OptInfoService.DeleteOperatorInfoByID: %s]", ccErr.Detail())
		oh.Response(r.Context(), oh.Logger, w, ccErr, nil)
		return
	}
	oh.Response(r.Context(), oh.Logger, w, nil, nil)
	return
}
