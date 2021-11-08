package handler

import (
	"net/http"
	//"unicode/utf8"

	"analysis-server/api/service"
	//"analysis-server/api/utils"
	"analysis-server/model"
	//cons "common/constant"
	"common/log"
)

type ResourceInfoHandlers struct {
	CCHandler
	Logger     *log.Logger
	ResService *service.ResouceInfoService
}

func (rh *ResourceInfoHandlers) InitResourceInfo(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := rh.HttpRequestParse(r, params)
	if err != nil {
		rh.Logger.ErrorContext(r.Context(), "[resource/InitResourceInfo] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrResInfo, service.ErrMalformed, service.ErrNull, err.Error())
		rh.Response(r.Context(), rh.Logger, w, ccErr, nil)
		return
	}

	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrResInfo, service.ErrMiss, service.ErrId, service.ErrNull)
		rh.Response(r.Context(), rh.Logger, w, ccErr, nil)
		return
	}
	requestId := rh.GetTraceId(r)

	resViews, ccErr := rh.ResService.GetResouceByOptId(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		rh.Logger.WarnContext(r.Context(), "[resource/InitResourceInfo/ServerHTTP] [ResService.GetResouceByOptId: %s]", ccErr.Detail())
		rh.Response(r.Context(), rh.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{int64(len(resViews)), resViews}
	rh.Response(r.Context(), rh.Logger, w, nil, dataBuf)
	return
}
