package service

import (
	"analysis-server/api/db"
	"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	comUtils "common/utils"
	"context"
	"database/sql"
	"time"
)

type AuthenService struct {
	Logger     *log.Logger
	LogInfoDao *db.LoginInfoDao
	OptInfoDao *db.OperatorInfoDao
	Db         *sql.DB
}

func (as *AuthenService) Login(ctx context.Context, params *model.LoginInfoParams,
	requestId string) (*model.LoginInfoView, CcError) {
	as.Logger.InfoContext(ctx, "Login method start, "+"login name:%s", *params.Name)
	FuncName := "AuthenService/login"
	tx, err := as.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer RollbackLog(ctx, as.Logger, FuncName, tx)
	//generate login information
	loginInfo := new(model.LoginInfo)
	loginInfo.Name = *params.Name
	loginInfo.Status = utils.UserOnline
	loginInfo.ClientIp = *params.ClientIp
	loginInfo.BeginedAt = time.Now()

	if err := as.LogInfoDao.Create(ctx, tx, loginInfo); err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [LogInfoDao.Create: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	loginView := as.LoginInfoMdelToView(loginInfo)
	as.Logger.InfoContext(ctx, "CreateLoginInfo method end,login name:%s", *params.Name)
	//update the operator information
	updateParams := make(map[string]interface{}, 2)
	updateParams["UpdatedAt"] = time.Now()
	updateParams["Status"] = utils.UserOnline
	err = as.OptInfoDao.Update(ctx, tx, *params.Name, updateParams)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	as.Logger.InfoContext(ctx, "the OptInfoDao.Update end, login name:%s", *params.Name)
	return loginView, nil
}

func (as *AuthenService) Logout(ctx context.Context, strUserName string) CcError {
	FuncName := "AuthenService/Logout"
	tx, err := as.Db.Begin()
	if err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer RollbackLog(ctx, as.Logger, FuncName, tx)
	_, err = as.LogInfoDao.Get(ctx, tx, strUserName)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeUserNameWrong, ErrOperator, ErrNotFound, ErrNull, "the logining user is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//generate login information
	filterFields := make(map[string]interface{})
	filterFields["Name"] = strUserName
	filterFields["Status"] = utils.UserOnline
	updateFields := make(map[string]interface{})
	updateFields["EndedAt"] = time.Now()
	updateFields["Status"] = utils.UserOffline
	err = as.LogInfoDao.Update(ctx, tx, filterFields, updateFields)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//update the operator information
	delete(updateFields, "EndedAt")
	updateFields["UpdatedAt"] = time.Now()
	//updateFields["Status"] = utils.UserOffline
	err = as.OptInfoDao.Update(ctx, tx, strUserName, updateFields)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		as.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	return nil
}

func (as *AuthenService) ListLoginInfo(ctx context.Context,
	params *model.ListParams) ([]*model.LoginInfoView, int, CcError) {
	OptViewSlice := make([]*model.LoginInfoView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "name", "clientIp", "beginedAt", "endedAt", "role":
				filterFields[*f.Field] = f.Value
			default:
				return OptViewSlice, 0, NewError(ErrOperator, ErrUnsupported, ErrField, *f.Field)
			}
		}
	}
	if params.DescLimit != nil {
		limit = *params.DescLimit
		if params.DescOffset != nil {
			offset = *params.DescOffset
		}
	}
	orderField := ""
	orderDirection := 0
	if params.Order != nil {
		orderField = *params.Order[0].Field
		orderDirection = *params.Order[0].Direction
	}
	optInfos, err := as.LogInfoDao.List(ctx, as.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		as.Logger.ErrorContext(ctx, "[AuthenService/service/ListLoginInfo] [LogInfoDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return OptViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, optInfo := range optInfos {
		optInfoView := as.LoginInfoMdelToView(optInfo)
		OptViewSlice = append(OptViewSlice, optInfoView)
	}
	optInfoCount := len(optInfos)
	return OptViewSlice, optInfoCount, nil
}

// convert accSubject to LoginInfoView ...
func (as *AuthenService) LoginInfoMdelToView(loginInfo *model.LoginInfo) *model.LoginInfoView {
	loginView := new(model.LoginInfoView)
	loginView.Name = loginInfo.Name
	loginView.ClientIp = loginInfo.ClientIp
	loginView.Status = loginInfo.Status
	loginView.BeginedAt = loginInfo.BeginedAt
	loginView.EndedAt = loginInfo.EndedAt
	loginView.AccessToken = comUtils.Uuid()
	return loginView
}

// func (as*AuthenService) GetLoginInfoByName(ctx context.Context, strUserName string,
// 	requestId string) (*model.LoginInfoView, CcError) {
// 	optInfo, err := as.LogInfoDao.Get(ctx, as.Db, strUserName)
// 	switch err {
// 	case nil:
// 	case sql.ErrNoRows:
// 		return nil, NewCcError(cons.CodeUserNameWrong, ErrOperator, ErrNotFound, ErrNull, "the user name is not exist")
// 	default:
// 		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	optInfoView := as.LoginInfoMdelToView(optInfo)
// 	return optInfoView, nil
// }
