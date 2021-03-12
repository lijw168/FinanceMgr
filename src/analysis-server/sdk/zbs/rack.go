package zbs

import (
	"errors"
	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type Rack struct {
}

func (s *Rack) CreateRack(opts *options.CreateRackOptions) (*model.RackView, error) {
	action := "CreateRack"
	switch {
	case opts.Tag == "":
		return nil, errors.New("tag is required")
	}
	params := model.CreateRackParams{
		Tag: &opts.Tag,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.RackView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (s *Rack) DeleteRack(opts *options.DeleteOptions) (*model.AtDeleteView, error) {
	if opts.Id == "" {
		return nil, errors.New("Id are required")
	}
	action := "DeleteRack"
	para := &model.DeleteRackParams{
		Id: &opts.Id,
	}
	dv, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}

	view := &model.AtDeleteView{}
	err = util.FormatView(dv.Data, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (s *Rack) DescribeRack(opts *options.BaseOptions) (*model.RackView, error) {
	action := "DescribeRack"
	bt, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.RackView{}
	err = util.FormatView(bt, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (z *Rack) DescribeRacks(opts *options.ListOptions) (int64, []*model.RackView, error) {
	action := "DescribeRacks"
	var ret []*model.RackView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}
