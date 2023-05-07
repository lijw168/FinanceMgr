package mgr

import (
	"financeMgr/src/analysis-server/model"
	"financeMgr/src/analysis-server/sdk/options"
	"financeMgr/src/analysis-server/sdk/util"
	// "errors"
	// "fmt"
)

type MenuInfo struct {
}

func (mi *MenuInfo) ListMenuInfo_json(params []byte) ([]byte, error) {
	action := "ListMenuInfo"
	return ListOpsResources_json(action, params)
}

func (mi *MenuInfo) ListMenuInfo(opts *options.ListOptions) (int64, []*model.MenuInfoView, error) {
	action := "ListMenuInfo"
	var ret []*model.MenuInfoView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}
