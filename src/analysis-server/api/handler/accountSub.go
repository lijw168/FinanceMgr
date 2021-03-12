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

type AccountSubHandlers struct {
	CCHandler
	Logger        *log.Logger
	AccSubService service.AccountSubService
}

func (ah *AccountSubHandlers) ListAccSub(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ListSubjectParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/ListAccSub] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrAccSub, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["subjectId"] = utils.Attribute{utils.T_Int, nil}
		filterMap["subjectName"] = utils.Attribute{utils.T_String, nil}
		filterMap["subjectLevel"] = utils.Attribute{utils.T_Int, nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrAccSub, service.ErrValue, service.ErrField, service.ErrNull)
			ah.Response(r.Context(), ah.Logger, w, ce, nil)
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
	// 		ah.Response(r.Context(), ah.Logger, w, ce, nil)
	// 		return
	// 	}
	// 	switch *params.Order[0].Direction {
	// 	case cons.Order_Asc, cons.Order_Desc:
	// 	default:
	// 		ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, string(*params.Order[0].Direction))
	// 		ah.Response(r.Context(), ah.Logger, w, ce, nil)
	// 		return
	// 	}
	// }

	accSubViews, count, ccErr := ah.AccSubService.ListAccSub(r.Context(), params)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[accSub/ListAccSub/ServeHTTP] [AccSubService.ListAccSub: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), accSubViews}
	ah.Response(r.Context(), ah.Logger, w, nil, dataBuf)
	return
}

// func (ah *AccountSubHandlers) GetAccSubByName(w http.ResponseWriter, r *http.Request) {
// 	var params = new(model.DescribeNameParams)
// 	err := ah.HttpRequestParse(r, params)
// 	if err != nil {
// 		ah.Logger.ErrorContext(r.Context(), "[accSub/GetAccSubByName] [HttpRequestParse: %v]", err)
// 		ccErr := service.NewError(service.ErrAccSub, service.ErrMalformed, service.ErrNull, err.Error())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}

// 	if params.Name == nil || *params.Name == "" {
// 		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrName, service.ErrNull)
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	requestId := ah.GetTraceId(r)

// 	accSubView, ccErr := ah.AccSubService.GetAccSubByName(r.Context(), params, requestId)
// 	if ccErr != nil {
// 		ah.Logger.WarnContext(r.Context(), "[accSub/GetAccSubByName/ServeHTTP] [AccSubService.GetaccSubByName: %s]", ccErr.Detail())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	ah.Response(r.Context(), ah.Logger, w, nil, accSubView)
// 	return
// }

func (ah *AccountSubHandlers) GetAccSub(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/GetAccSub] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrAccSub, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}

	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.Id, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	requestId := ah.GetTraceId(r)

	accSubView, ccErr := ah.AccSubService.GetAccSubById(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[accSub/GetAccSub/ServeHTTP] [AccSubService.GetAccSubById: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, accSubView)
	return
}

func (ah *AccountSubHandlers) CreateAccSub(w http.ResponseWriter, r *http.Request) {
	var params = new(model.CreateSubjectParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/CreateAccSub] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrAccSub, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.SubjectName == nil || *params.SubjectName == "" {
		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrName, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if utf8.RuneCountInString(*params.SubjectName) > NameMaxLen || !utils.VerStrP(*params.SubjectName) {
		ccErr := service.NewError(service.ErrAccSub, service.ErrInvalid, service.ErrName, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}

	if params.SubjectLevel == nil || *params.SubjectLevel > 4 {
		ccErr = service.NewCcError(cons.CodeInvalLevel, service.ErrAccSub, service.ErrInvalid, service.ErrSubLevel, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}

	requestId := ah.GetTraceId(r)
	accSubView, ccErr = ah.AccSubService.CreateAccSub(r.Context(), params, requestId)
	ah.Logger.InfoContext(r.Context(), "AccSubService.CreateAccSub in CreateAccSub.")
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[accSub/CreateAccSub/ServeHTTP] [AccSubService.CreateAccSub: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, accSubView)
	return
}

func (ah *AccountSubHandlers) UpdateAccSub(w http.ResponseWriter, r *http.Request) {
	var params = new(model.ModifySubjectParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/AccountSubHandlers] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrAccSub, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}

	if params.SubjectID == nil || *params.SubjectID == "" {
		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}

	updateFields := make(map[string]interface{})
	if params.SubjectName != nil {
		if utf8.RuneCountInString(*params.SubjectName) > NameMaxLen || !utils.VerStrP(*params.SubjectName) {
			ccErr := service.NewError(service.ErrAccSub, service.ErrInvalid, service.ErrSubjectName, service.ErrNull)
			ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
			return
		}
		updateFields["SubjectName"] = *params.SubjectName
	}
	if params.SubjectLevel != nil {
		updateFields["SubjectLevel"] = *params.SubjectLevel
	}
	if len(updateFields) == 0 {
		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ccErr := ah.AccSubService.UpdateAccSubById(r.Context(), *params.SubjectID, updateFields)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[accSub/AccountSubHandlers/ServeHTTP] [AccSubService.UpdateAccSubById: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, nil)
	return
}

func (ah *AccountSubHandlers) DeleteAccSub(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DeleteSubjectParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/DeleteAccSub] [HttpRequestParse: %v]", err)
		ccErr := service.NewError(service.ErrAccSub, service.ErrMalformed, service.ErrNull, err.Error())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.ID == nil || *params.ID <= 0 {
		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	requestId := ah.GetTraceId(r)
	ccErr = ah.AccSubService.DeleteAccSubByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[accSub/DeleteAccSub/ServeHTTP] [AccSubService.DeleteaccSubByName: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, nil)
	return
}
