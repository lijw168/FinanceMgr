package handler

import (
	"net/http"
	//"unicode/utf8"

	"financeMgr/src/analysis-server/api/service"
	"financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"

	//cons "financeMgr/src/commonconstant"
	"financeMgr/src/common/log"
)

type MenuInfoHandlers struct {
	CCHandler
	Logger      *log.Logger
	MenuService *service.MenuInfoService
}

func (ah *MenuInfoHandlers) ListMenuInfo(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/ListMenuInfo] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrMenuInfo, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	//暂时先不添加这个功能
	// if isLackBaseParams([]string{"subjectId", "companyId"}, params.Filter) {
	// 	ah.Logger.ErrorContext(r.Context(), "lack base param  operatorId")
	// 	ce := service.NewError(service.ErrMenuInfo, service.ErrMiss, service.ErrId, service.ErrNull)
	// 	ah.Response(r.Context(), ah.Logger, w, ce, nil)
	// 	return
	// }
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["MenuId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["menuName"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["menuLevel"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["parentMenuId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrMenuInfo, service.ErrValue, service.ErrField, service.ErrNull)
			ah.Response(r.Context(), ah.Logger, w, ce, nil)
			return
		}
	}
	accSubViews, count, ccErr := ah.MenuService.ListMenuInfo(r.Context(), params)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[accSub/ListMenuInfo/ServerHTTP] [MenuService.ListMenuInfo: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), accSubViews}
	ah.Response(r.Context(), ah.Logger, w, nil, dataBuf)
	return
}
