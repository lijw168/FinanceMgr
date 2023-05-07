package service

import (
	"context"
	"database/sql"

	"financeMgr/src/analysis-server/api/db"
	//"financeMgr/src/analysis-server/api/utils"
	"financeMgr/src/analysis-server/model"
	cons "financeMgr/src/common/constant"
	"financeMgr/src/common/log"
)

type ResouceInfoService struct {
	Logger     *log.Logger
	VInfoDao   *db.VoucherInfoDao
	CompanyDao *db.CompanyDao
	Db         *sql.DB
}

//可以优化一下GetCompanyByOperatorId这个函数的返回值。
func (rs *ResouceInfoService) GetResouceByOptId(ctx context.Context, operatorId int,
	requestId string) ([]*model.ResourceInfoView, CcError) {
	//create
	rs.Logger.InfoContext(ctx, "GetResouceByOptId method start, "+"operator:%s", operatorId)
	FuncName := "ResouceInfoService/Resource/GetResouceByOptId"
	bIsRollBack := true
	// Begin transaction
	tx, err := rs.Db.Begin()
	if err != nil {
		rs.Logger.ErrorContext(ctx, "[%s] [DB.Begin: %s]", FuncName, err.Error())
		return nil, NewError(ErrSystem, ErrError, ErrNull, "tx begin error")
	}
	defer func() {
		if bIsRollBack {
			RollbackLog(ctx, rs.Logger, FuncName, tx)
		}
	}()
	//get company info
	comInfo, err := rs.CompanyDao.GetCompanyByOperatorId(ctx, tx, operatorId)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, NewCcError(cons.CodeComInfoNotExist, ErrCompany, ErrNotFound, ErrNull, "the company information is not exist")
	default:
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	//get resource info
	resInfoSlice := make([]*model.ResourceInfoView, 0)
	if comInfo.CompanyGroupID > 0 {
		//获取同一个公司组里的所有的公司信息
		filterFields := make(map[string]interface{})
		limit, offset := -1, 0
		filterFields["companyGroupId"] = comInfo.CompanyGroupID
		orderField := ""
		orderDirection := 0
		comInfos, err := rs.CompanyDao.List(ctx, tx, filterFields, limit, offset, orderField, orderDirection)
		if err != nil {
			rs.Logger.ErrorContext(ctx, "[ResouceInfoService/service/GetResouceByOptId] [CompanyDao.List: %s, filterFields: %v]",
				err.Error(), filterFields)
			return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
		}
		for _, comInfo := range comInfos {
			resInfo, err := rs.getResourceData(ctx, comInfo)
			if err != nil {
				return nil, err
			}
			resInfoSlice = append(resInfoSlice, resInfo)
		}
	} else if comInfo.CompanyGroupID == 0 {
		resInfo, err := rs.getResourceData(ctx, comInfo)
		if err != nil {
			return nil, err
		}
		resInfoSlice = append(resInfoSlice, resInfo)
	} else {
		panic("company group id is negative")
	}
	if err = tx.Commit(); err != nil {
		rs.Logger.ErrorContext(ctx, "[%s] [Commit Err: %v]", FuncName, err)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	bIsRollBack = false
	rs.Logger.InfoContext(ctx, "GetResouceByOptId method end")
	return resInfoSlice, nil
}

//这是中间计算的函数
func (rs *ResouceInfoService) getResourceData(ctx context.Context,
	pComView *model.CompanyInfo) (*model.ResourceInfoView, CcError) {
	resInfo := new(model.ResourceInfoView)
	resInfo.CompanyId = pComView.CompanyID
	resInfo.CompanyName = pComView.CompanyName
	iStartAccountYear := pComView.StartAccountPeriod / 100
	iLatestAccountYear := pComView.LatestAccountYear
	yearSlice := make([]int, 0, (iLatestAccountYear - iStartAccountYear + 1))
	for i := iStartAccountYear; i <= iLatestAccountYear; i++ {
		yearSlice = append(yearSlice, i)
	}
	resInfo.YearSlice = yearSlice
	return resInfo, nil
}
