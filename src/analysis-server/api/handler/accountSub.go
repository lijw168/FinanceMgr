package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"financeMgr/src/analysis-server/api/service"
	"financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
	"financeMgr/src/common/log"
)

type AccountSubHandlers struct {
	CCHandler
	Logger        *log.Logger
	AccSubService *service.AccountSubService
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
	if isLackBaseParams([]string{"subjectId", "companyId"}, params.Filter) {
		ah.Logger.ErrorContext(r.Context(), "lack base param  operatorId")
		ce := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ce, nil)
		return
	}
	if params.Filter != nil {
		filterMap := map[string]utils.Attribute{}
		filterMap["subjectId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["companyId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["subjectName"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["commonId"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["subjectLevel"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["subjectDirection"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["subjectType"] = utils.Attribute{Type: utils.T_Int, Val: nil}
		filterMap["subjectStyle"] = utils.Attribute{Type: utils.T_String, Val: nil}
		filterMap["mnemonicCode"] = utils.Attribute{Type: utils.T_String, Val: nil}
		if !utils.ValiFilter(filterMap, params.Filter) {
			ce := service.NewError(service.ErrAccSub, service.ErrValue, service.ErrField, service.ErrNull)
			ah.Response(r.Context(), ah.Logger, w, ce, nil)
			return
		}
	}
	if (params.Order != nil) && (len(params.Order) > 0) {
		switch *params.Order[0].Field {
		case "subjectId":
			*params.Order[0].Field = "subjectId"
		case "commonId":
			*params.Order[0].Field = "commonId"
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
			ah.Response(r.Context(), ah.Logger, w, ce, nil)
			return
		}
		switch *params.Order[0].Direction {
		case cons.Order_Asc, cons.Order_Desc:
		default:
			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd,
				strconv.Itoa(*params.Order[0].Direction))
			ah.Response(r.Context(), ah.Logger, w, ce, nil)
			return
		}
	}
	accSubViews, count, ccErr := ah.AccSubService.ListAccSub(r.Context(), params)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[accSub/ListAccSub/ServerHTTP] [AccSubService.ListAccSub: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), accSubViews}
	ah.Response(r.Context(), ah.Logger, w, nil, dataBuf)
	return
}

// func (ah *AccountSubHandlers) ListYearBalance(w http.ResponseWriter, r *http.Request) {
// 	var params = new(model.ListParams)
// 	err := ah.HttpRequestParse(r, params)
// 	if err != nil {
// 		ah.Logger.ErrorContext(r.Context(), "[accSub/ListYearBalance] [HttpRequestParse: %v]", err)
// 		ccErr := service.NewError(service.ErrAccSub, service.ErrMalformed, service.ErrNull, err.Error())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	if isLackBaseParams([]string{"subjectId", "companyId"}, params.Filter) {
// 		ah.Logger.ErrorContext(r.Context(), "lack base param  operatorId")
// 		ce := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrId, service.ErrNull)
// 		ah.Response(r.Context(), ah.Logger, w, ce, nil)
// 		return
// 	}
// 	//由于目前获取期初余额时的过滤条件，可能只用到subjectId和companyId，所以也就不进行参数判断了。
// 	// if params.Filter != nil {
// 	// 	filterMap := map[string]utils.Attribute{}
// 	// 	filterMap["subjectId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
// 	// 	filterMap["companyId"] = utils.Attribute{Type: utils.T_Int, Val: nil}
// 	// 	filterMap["subjectName"] = utils.Attribute{Type: utils.T_String, Val: nil}
// 	// 	filterMap["commonId"] = utils.Attribute{Type: utils.T_String, Val: nil}
// 	// 	filterMap["subjectLevel"] = utils.Attribute{Type: utils.T_Int, Val: nil}
// 	// 	filterMap["subjectDirection"] = utils.Attribute{Type: utils.T_Int, Val: nil}
// 	// 	filterMap["subjectType"] = utils.Attribute{Type: utils.T_Int, Val: nil}
// 	// 	filterMap["subjectStyle"] = utils.Attribute{Type: utils.T_String, Val: nil}
// 	// 	filterMap["mnemonicCode"] = utils.Attribute{Type: utils.T_String, Val: nil}
// 	// 	if !utils.ValiFilter(filterMap, params.Filter) {
// 	// 		ce := service.NewError(service.ErrAccSub, service.ErrValue, service.ErrField, service.ErrNull)
// 	// 		ah.Response(r.Context(), ah.Logger, w, ce, nil)
// 	// 		return
// 	// 	}
// 	// }
// 	if (params.Order != nil) && (len(params.Order) > 0) {
// 		switch *params.Order[0].Field {
// 		case "subjectId":
// 			*params.Order[0].Field = "subjectId"
// 		default:
// 			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrField, *params.Order[0].Field)
// 			ah.Response(r.Context(), ah.Logger, w, ce, nil)
// 			return
// 		}
// 		switch *params.Order[0].Direction {
// 		case cons.Order_Asc, cons.Order_Desc:
// 		default:
// 			ce := service.NewError(service.ErrOrder, service.ErrInvalid, service.ErrOd, strconv.Itoa(*params.Order[0].Direction))
// 			ah.Response(r.Context(), ah.Logger, w, ce, nil)
// 			return
// 		}
// 	}
// 	yearBalViews, count, ccErr := ah.AccSubService.ListYearBalance(r.Context(), params)
// 	if ccErr != nil {
// 		ah.Logger.WarnContext(r.Context(), "[accSub/ListYearBalance/ServerHTTP] [AccSubService.ListYearBalance: %s]", ccErr.Detail())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	dataBuf := &DescData{(int64)(count), yearBalViews}
// 	ah.Response(r.Context(), ah.Logger, w, nil, dataBuf)
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
		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	requestId := ah.GetTraceId(r)

	accSubView, ccErr := ah.AccSubService.GetAccSubById(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[accSub/GetAccSub/ServerHTTP] [AccSubService.GetAccSubById: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, accSubView)
	return
}

// func (ah *AccountSubHandlers) GetYearBalance(w http.ResponseWriter, r *http.Request) {
// 	var params = new(model.DescribeIdParams)
// 	err := ah.HttpRequestParse(r, params)
// 	if err != nil {
// 		ah.Logger.ErrorContext(r.Context(), "[accSub/GetYearBalance] [HttpRequestParse: %v]", err)
// 		ccErr := service.NewError(service.ErrYearBalance, service.ErrMalformed, service.ErrNull, err.Error())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}

// 	if params.ID == nil || *params.ID <= 0 {
// 		ccErr := service.NewError(service.ErrYearBalance, service.ErrMiss, service.ErrId, service.ErrNull)
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	requestId := ah.GetTraceId(r)

// 	yearBal, ccErr := ah.AccSubService.GetYearBalanceById(r.Context(), *params.ID, requestId)
// 	if ccErr != nil {
// 		ah.Logger.WarnContext(r.Context(), "[accSub/GetYearBalance/ServerHTTP] [AccSubService.GetYearBalanceById: %s]", ccErr.Detail())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	ah.Response(r.Context(), ah.Logger, w, nil, yearBal)
// 	return
// }

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
	if params.CommonID == nil || *params.CommonID == "" {
		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrCommonId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if utf8.RuneCountInString(*params.CommonID) > 10 || !utils.VerCommonIdP(*params.CommonID) {
		ccErr := service.NewError(service.ErrAccSub, service.ErrInvalid, service.ErrCommonId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}

	if params.SubjectLevel == nil || *params.SubjectLevel > 4 {
		ccErr := service.NewCcError(cons.CodeInvalAccSubLevel, service.ErrAccSub, service.ErrInvalid, service.ErrSubLevel, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.SubjectDirection == nil || *params.SubjectDirection > 2 {
		ccErr := service.NewCcError(cons.CodeInvalAccSubDir, service.ErrAccSub, service.ErrInvalid, service.ErrSubdir, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.SubjectType == nil || *params.SubjectType > 5 {
		ccErr := service.NewCcError(cons.CodeInvalAccSubType, service.ErrAccSub, service.ErrInvalid, service.ErrType, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.CompanyID == nil || *(params.CompanyID) <= 0 {
		ccErr := service.NewError(service.ErrAccSub, service.ErrInvalid, service.ErrCompanyId, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	if params.SubjectStyle == nil || *params.SubjectStyle == "" {
		ccErr := service.NewError(service.ErrAccSub, service.ErrInvalid, service.ErrSubStyle, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	requestId := ah.GetTraceId(r)
	accSubView, ccErr := ah.AccSubService.CreateAccSub(r.Context(), params, requestId)
	ah.Logger.InfoContext(r.Context(), "AccSubService.CreateAccSub in CreateAccSub.")
	if ccErr != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/CreateAccSub/ServerHTTP] [AccSubService.CreateAccSub: %s]", ccErr.Detail())
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

	if params.SubjectID == nil || *params.SubjectID <= 0 {
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
		updateFields["subjectName"] = *params.SubjectName
	}
	if params.CommonID != nil {
		if utf8.RuneCountInString(*params.CommonID) > 10 || !utils.VerCommonIdP(*params.CommonID) {
			ccErr := service.NewError(service.ErrAccSub, service.ErrInvalid, service.ErrCommonId, service.ErrNull)
			ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
			return
		}
		updateFields["commonId"] = *params.CommonID
	}
	if params.SubjectLevel != nil {
		updateFields["subjectLevel"] = *params.SubjectLevel
	}
	if params.SubjectDirection != nil {
		updateFields["subjectDirection"] = *params.SubjectDirection
	}
	if params.SubjectType != nil {
		updateFields["subjectType"] = *params.SubjectType
	}
	if params.CompanyID != nil {
		updateFields["companyId"] = *params.CompanyID
	}
	if params.SubjectStyle != nil {
		updateFields["subjectStyle"] = *params.SubjectStyle
	}
	if params.MnemonicCode != nil {
		updateFields["mnemonicCode"] = *params.MnemonicCode
	}
	if len(updateFields) == 0 {
		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ccErr := ah.AccSubService.UpdateAccSubById(r.Context(), *params.SubjectID, updateFields)
	if ccErr != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/AccountSubHandlers/ServerHTTP] [AccSubService.UpdateAccSubById: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, nil)
	return
}

//因为期初余额字段，在任何时候都可以修改，所以修改该字段和修改其他字段进行分开。
// func (ah *AccountSubHandlers) UpdateYearBalance(w http.ResponseWriter, r *http.Request) {
// 	var params = new(model.OptYearBalanceParams)
// 	err := ah.HttpRequestParse(r, params)
// 	if err != nil {
// 		ah.Logger.ErrorContext(r.Context(), "[accSub/UpdateYearBalance] [HttpRequestParse: %v]", err)
// 		ccErr := service.NewError(service.ErrAccSub, service.ErrMalformed, service.ErrNull, err.Error())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}

// 	if params.SubjectID == nil || *params.SubjectID <= 0 {
// 		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrId, service.ErrNull)
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}

// 	if params.Balance == nil {
// 		ccErr := service.NewError(service.ErrAccSub, service.ErrMiss, service.ErrChangeContent, service.ErrNull)
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	ccErr := ah.AccSubService.UpdateYearBalanceById(r.Context(), *params.SubjectID, *params.Balance)
// 	if ccErr != nil {
// 		ah.Logger.ErrorContext(r.Context(), "[accSub/AccountSubHandlers/ServerHTTP] [AccSubService.UpdateAccSubById: %s]", ccErr.Detail())
// 		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
// 		return
// 	}
// 	ah.Response(r.Context(), ah.Logger, w, nil, nil)
// 	return
// }

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
	ccErr := ah.AccSubService.DeleteAccSubByID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/DeleteAccSub/ServerHTTP] [AccSubService.DeleteAccSubByID: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, nil)
	return
}

func (ah *AccountSubHandlers) QueryAccSubReference(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/QueryAccSubReference] [HttpRequestParse: %v]", err)
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
	iCount, ccErr := ah.AccSubService.QueryAccSubReferenceBySubID(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		errInfo := fmt.Sprintf("[accSub/QueryAccSubReference/ServerHTTP] [AccSubService.JudgeAccSubReferenceBySubID: %s]", ccErr.Detail())
		ah.Logger.ErrorContext(r.Context(), errInfo)
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	//dataBuf := &DescData{iCount, nil}
	ah.Response(r.Context(), ah.Logger, w, nil, iCount)
	return
}

func (ah *AccountSubHandlers) CopyAccSubTemplate(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/CopyAccSubTemplate] [HttpRequestParse: %v]", err)
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
	accSubViews, count, ccErr := ah.AccSubService.CopyAccSubTemplate(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ah.Logger.WarnContext(r.Context(), "[accSub/CopyAccSubTemplate/ServerHTTP] [AccSubService.CopyAccSubTemplate: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	dataBuf := &DescData{(int64)(count), accSubViews}
	ah.Response(r.Context(), ah.Logger, w, nil, dataBuf)
	return
}

func (ah *AccountSubHandlers) GenerateAccSubTemplate(w http.ResponseWriter, r *http.Request) {
	var params = new(model.DescribeIdParams)
	err := ah.HttpRequestParse(r, params)
	if err != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/GenerateAccSubTemplate] [HttpRequestParse: %v]", err)
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
	ccErr := ah.AccSubService.GenerateAccSubTemplate(r.Context(), *params.ID, requestId)
	if ccErr != nil {
		ah.Logger.ErrorContext(r.Context(), "[accSub/GenerateAccSubTemplate/ServerHTTP] [AccSubService.GenerateAccSubTemplate: %s]", ccErr.Detail())
		ah.Response(r.Context(), ah.Logger, w, ccErr, nil)
		return
	}
	ah.Response(r.Context(), ah.Logger, w, nil, nil)
	return
}
