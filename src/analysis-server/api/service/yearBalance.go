package service

import (
	"analysis-server/api/db"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"context"
	"database/sql"
)

type YearBalanceService struct {
	Logger     *log.Logger
	YearBalDao *db.YearBalanceDao
	Db         *sql.DB
}

func (ys *YearBalanceService) CreateYearBalance(ctx context.Context, params *model.YearBalanceParams,
	requestId string) CcError {
	ys.Logger.InfoContext(ctx, "CreateYearBalance method start, "+"subjectName:%s", *params.SubjectID)
	if err := ys.YearBalDao.Create(ctx, ys.Db, params); err != nil {
		FuncName := "YearBalanceService/yearBalance/CreateYearBalance"
		ys.Logger.ErrorContext(ctx, "[%s] [YearBalDao.Create: %s]", FuncName, err.Error())
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	ys.Logger.InfoContext(ctx, "CreateYearBalance method end, "+"subjectName:%s", *params.SubjectID)
	return nil
}

func (ys *YearBalanceService) GetYearBalanceById(ctx context.Context, subjectID int,
	requestId string) (*model.YearBalanceView, CcError) {
	accSubView, err := ys.YearBalDao.GetByID(ctx, ys.Db, subjectID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeYearBalanceNotExist, ErrYearBalance, ErrNotFound, ErrNull, "the year balance is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	return accSubView, nil
}

func (ys *YearBalanceService) DeleteYearBalanceByID(ctx context.Context, subjectID int,
	requestId string) CcError {
	ys.Logger.InfoContext(ctx, "DeleteYearBalanceByID method begin, "+"subject:%d", subjectID)
	err := ys.YearBalDao.DeleteByID(ctx, ys.Db, subjectID)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	ys.Logger.InfoContext(ctx, "DeleteYearBalanceByID method end, "+"subject:%d", subjectID)
	return nil
}

func (ys *YearBalanceService) UpdateYearBalanceById(ctx context.Context, subjectID int,
	params map[string]interface{}) CcError {
	err := ys.YearBalDao.UpdateBySubID(ctx, ys.Db, subjectID, params)
	if err != nil {
		return NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	return nil
}
