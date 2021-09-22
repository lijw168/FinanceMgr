package service

import (
	"analysis-server/api/db"
	//"analysis-server/api/utils"
	"analysis-server/model"
	//cons "common/constant"
	"common/log"
	"context"
	"database/sql"
)

type MenuInfoService struct {
	Logger  *log.Logger
	MenuDao *db.MenuInfoDao
	Db      *sql.DB
}

// convert accSubject to accSubjectView ...
func (ms *MenuInfoService) MenuInfoModelToView(menuInfo *model.MenuInfo) *model.MenuInfoView {
	menuView := new(model.MenuInfoView)
	menuView.MenuID = menuInfo.MenuID
	menuView.MenuName = menuInfo.MenuName
	menuView.MenuLevel = menuInfo.MenuLevel
	menuView.ParentMenuID = menuInfo.ParentMenuID
	return menuView
}

func (ms *MenuInfoService) ListMenuInfo(ctx context.Context,
	params *model.ListParams) ([]*model.MenuInfoView, int, CcError) {
	menuInfoViewSlice := make([]*model.MenuInfoView, 0)
	filterFields := make(map[string]interface{})
	limit, offset := -1, 0
	if params.Filter != nil {
		for _, f := range params.Filter {
			switch *f.Field {
			case "menuId", "menuName", "menuLevel", "parentMenuId":
				filterFields[*f.Field] = f.Value
			default:
				return menuInfoViewSlice, 0, NewError(ErrMenuInfo, ErrUnsupported, ErrField, *f.Field)
			}
		}
	}
	if params.DescLimit != nil {
		limit = *params.DescLimit
		if params.DescOffset != nil {
			offset = *params.DescOffset
		}
	}
	orderField := "menuSerialNum"
	//asc
	orderDirection := 0
	// if params.Order != nil {
	// 	orderField = *params.Order[0].Field
	// 	orderDirection = *params.Order[0].Direction
	// }
	menuInfos, err := ms.MenuDao.List(ctx, ms.Db, filterFields, limit, offset, orderField, orderDirection)
	if err != nil {
		ms.Logger.ErrorContext(ctx, "[MenuInfoService/service/ListMenuInfo] [MenuDao.List: %s, filterFields: %v]", err.Error(), filterFields)
		return menuInfoViewSlice, 0, NewError(ErrSystem, ErrError, ErrNull, err.Error())
	}

	for _, menuInfo := range menuInfos {
		menuInfoView := ms.MenuInfoModelToView(menuInfo)
		menuInfoViewSlice = append(menuInfoViewSlice, menuInfoView)
	}
	accSubInfoCount := len(menuInfos)
	return menuInfoViewSlice, accSubInfoCount, nil
}
