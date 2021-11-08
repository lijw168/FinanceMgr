package service

import (
	"context"
	"database/sql"

	"analysis-server/api/db"
	//"analysis-server/api/utils"
	"analysis-server/model"
	cons "common/constant"
	"common/log"
	"fmt"
)

type ResouceInfoService struct {
	Logger     *log.Logger
	VInfoDao   *db.VoucherInfoDao
	CompanyDao *db.CompanyDao
	Db         *sql.DB
}

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
	comInfo, err := rs.CompanyDao.GetComGroupIdByOperatorId(ctx, tx, operatorId)
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
			resInfo, err := rs.getResourceData(ctx, comInfo, tx)
			if err != nil {
				return nil, err
			}
			resInfoSlice = append(resInfoSlice, resInfo)
		}
	} else if comInfo.CompanyGroupID == 0 {
		resInfo, err := rs.getResourceData(ctx, comInfo, tx)
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
func (rs *ResouceInfoService) getResourceData(ctx context.Context, pComView *model.CompanyInfo,
	tx *sql.Tx) (*model.ResourceInfoView, CcError) {
	resInfo := new(model.ResourceInfoView)
	resInfo.CompanyId = pComView.CompanyID
	resInfo.CompanyName = pComView.CompanyName
	//get voucherInfo
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	orderField := ""
	orderDirection := 0
	filterFields["companyId"] = pComView.CompanyID
	voucherInfos, err := rs.VInfoDao.SimpleList(ctx, tx, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		rs.Logger.ErrorContext(ctx, "[ResouceInfoService/service/GetResourceData] [VInfoDao.SimpleList: %s, filterFields: %v]",
			err.Error(), filterFields)
		return nil, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}
	yearSlice := make([]int, len(voucherInfos))
	for _, item := range voucherInfos {
		//先暂时这样写，等以后想到了解决方案，再继续实行该函数。
		//yearSlice = append(yearSlice, item.Year)
		fmt.Printf("year:%d", item.VoucherMonth)
		yearSlice = append(yearSlice, 2021)
	}
	resInfo.YearSlice = yearSlice
	return resInfo, nil
}
