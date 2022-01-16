package handler

import (
	"analysis-server/api/service"
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
	var params = new(model.DescribeIdParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/GetYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}

	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	requestId := yh.GetTraceId(r)

	yearBalView, ccErr := yh.YearBalService.GetYearBalanceById(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		yh.Logger.WarnContext(r.Context(), "[yearBalance/GetYearBalance/ServerHTTP] [YearBalService.GetYearBalanceById: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, yearBalView)
	return
}

func (yh *YearBalHandlers) CreateYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.YearBalanceParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/CreateYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.Summary == nil || *params.Summary == "" {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrYearBalSummary, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.SubjectDirection == nil || *params.SubjectDirection > 2 {
		ccErr := service.NewCcError(cons.CodeInvalAccSubDir, service.ErrYearBalance, service.ErrInvalid, service.ErrSubdir, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.SubjectID == nil || *(params.SubjectID) <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrInvalid, service.ErrId, service.ErrNull)
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

func (yh *YearBalHandlers) UpdateYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.YearBalanceParams)
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
	updateFields := make(map[string]interface{})
	if params.Summary != nil && *params.Summary != "" {
		updateFields["summary"] = *params.Summary
	}
	if params.SubjectDirection != nil {
		updateFields["subjectDirection"] = *params.SubjectDirection
	}
	if params.Balance != nil {
		updateFields["balance"] = *params.Balance
	}
	if len(updateFields) == 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	ccErr := yh.YearBalService.UpdateYearBalanceById(r.Context(), *params.SubjectID, updateFields)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/YearBalHandlers/ServerHTTP] [YearBalService.UpdateYearBalanceById: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
	return
}

func (yh *YearBalHandlers) DeleteYearBalance(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteIDParams)
	err := yh.HttpRequestParse(r, params)
	if err != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/DeleteYearBalance] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	requestId := yh.GetTraceId(r)
	ccErr := yh.YearBalService.DeleteYearBalanceByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		yh.Logger.ErrorContext(r.Context(), "[yearBalance/DeleteYearBalance/ServerHTTP] [YearBalService.DeleteYearBalanceByID: %s]", ccErr.Detail())
		yh.Response(r.Context(), yh.Logger, w, ccErr, nil)
		return
	}
	yh.Response(r.Context(), yh.Logger, w, nil, nil)
	return
}
