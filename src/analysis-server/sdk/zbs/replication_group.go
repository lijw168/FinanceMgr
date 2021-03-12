package zbs

import (
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type Rg struct {
}

/*
func (z *Rg) CreateRg(opts *options.CreateRgOptions) (*model.RgView, error) {
	action := "CreateRg"
	switch {
	case opts.DeviceId == "":
		return nil, errors.New("DeviceId is required")
	}
	params := model.CreateRgParams{
		PoolId: &opts.DeviceId,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.RgView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (s *Rg) DeleteRg(opts *options.BaseOptions) error {
	action := "DeleteRg"
	return DeleteOpsResource(action, opts)
}
*/

func (r *Rg) DescribeRg(opts *options.BaseOptions) (*model.RgView, error) {
	action := "DescribeRg"
	bt, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.RgView{}
	err = util.FormatView(bt, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (r *Rg) DescribeRgs(opts *options.ListOptions) (int64, []*model.RgView, error) {
	action := "DescribeRgs"
	var ret []*model.RgView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}
