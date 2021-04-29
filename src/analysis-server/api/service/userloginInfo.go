package service

import (
	"analysis-server/api/db"
	"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"context"
	"database/sql"
	"time"
)

type LoginInfoService struct {
	Logger     *log.Logger
	LogInfoDao *db.LoginInfoDao
	OptInfoDao *db.OperatorInfoDao
	Db         *sql.DB
}

func (ls *LoginInfoService) CreateLoginInfo(ctx context.Context, params *model.LoginInfoParams,
	requestId string) (*model.LoginInfoView, CcError) {
	ls.Logger.InfoContext(ctx, "CreateLoginInfo method start, "+"login name:%s", *params.Name)
	FuncName := "LoginInfoService/login/CreateLoginInfo"
	tx, err := ls.Db.Begin()
	if err != nil {
		ls.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	//defer RollbackLog(ctx, ls.Logger, FuncName, tx)
	//generate login information
	loginInfo := new(model.LoginInfo)
	loginInfo.Name = *params.Name
	loginInfo.Role = *params.Role
	loginInfo.ClientIp = *params.ClientIp
	loginInfo.BeginedAt = time.Now()

	if err := ls.LogInfoDao.Create(ctx, tx, loginInfo); err != nil {
		ls.Logger.ErrorContext(ctx, "[%s] [LogInfoDao.Create: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	loginView := ls.LoginInfoMdelToView(loginInfo)
	ls.Logger.InfoContext(ctx, "CreateLoginInfo method end,login name:%s", *params.Name)
	//update the operator information
	updateParams := make(map[string]interface{}, 2)
	updateParams["UpdatedAt"] = time.Now()
	updateParams["Status"] = utils.UserOnline
	err = ls.OptInfoDao.Update(ctx, tx, *params.Name, updateParams)
	if err != nil {
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ls.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	ls.Logger.InfoContext(ctx, "the OptInfoDao.Update end, login name:%s", *params.Name)
	return loginView, nil
}

func (ls *LoginInfoService) ListLoginInfo(ctx context.Context,
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
	optInfos, err := ls.LogInfoDao.List(ctx, ls.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		ls.Logger.ErrorContext(ctx, "[LoginInfoService/service/ListLoginInfo] [LogInfoDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return OptViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, optInfo := range optInfos {
		optInfoView := ls.LoginInfoMdelToView(optInfo)
		OptViewSlice = append(OptViewSlice, optInfoView)
	}
	optInfoCount := len(optInfos)
	return OptViewSlice, optInfoCount, nil
}

// convert accSubject to LoginInfoView ...
func (ls *LoginInfoService) LoginInfoMdelToView(loginInfo *model.LoginInfo) *model.LoginInfoView {
	loginView := new(model.LoginInfoView)
	loginView.Name = loginInfo.Name
	loginView.ClientIp = loginInfo.ClientIp
	loginView.Role = loginInfo.Role
	loginView.BeginedAt = loginInfo.BeginedAt
	loginView.EndedAt = loginInfo.EndedAt
	return loginView
}

// func (ls *LoginInfoService) GetLoginInfoByName(ctx context.Context, strUserName string,
// 	requestId string) (*model.LoginInfoView, CcError) {
// 	optInfo, err := ls.LogInfoDao.Get(ctx, ls.Db, strUserName)
// 	switch err {
// 	case nil:
// 	case sql.ErrNoRows:
// 		return nil, NewCcError(cons.CodeUserNameWrong, ErrOperator, ErrNotFound, ErrNull, "the user name is not exist")
// 	default:
// 		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
// 	}
// 	optInfoView := ls.LoginInfoMdelToView(optInfo)
// 	return optInfoView, nil
// }

func (ls *LoginInfoService) UpdateLoginInfo(ctx context.Context, strUserName string, params map[string]interface{}) CcError {
	FuncName := "LoginInfoService/UpdateLoginInfo"
	tx, err := ls.Db.Begin()
	if err != nil {
		ls.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	//defer RollbackLog(ctx, ls.Logger, FuncName, tx)
	_, err = ls.LogInfoDao.Get(ctx, tx, strUserName)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return NewCcError(cons.CodeUserNameWrong, ErrOperator, ErrNotFound, ErrNull, "the user name is not exist")
	default:
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	err = ls.LogInfoDao.Update(ctx, tx, strUserName, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	if err = tx.Commit(); err != nil {
		ls.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	return nil
}
