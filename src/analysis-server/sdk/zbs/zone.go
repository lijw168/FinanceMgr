package zbs

import (
	"errors"

	"analysis-server/model"
	"analysis-server/sdk/options"
	"analysis-server/sdk/util"
)

type Pool struct {
}

func (z *Pool) CreatePool(opts *options.CreatePoolOptions) (*model.PoolView, error) {
	action := "CreatePool"
	switch {
	case opts.Name == "":
		return nil, errors.New("name is required")
	}
	params := model.CreatePoolParams{
		Name:      &opts.Name,
		Type:      &opts.Type,
		MediaType: opts.MediaType,
		Rgc:       &opts.Rgc,
		ObjSize:   &opts.ObjSize,
	}
	result, err := util.DoRequest(action, params)
	if err != nil {
		return nil, err
	}
	view := &model.PoolView{}
	util.FormatView(result.Data, &view)
	return view, nil
}

func (s *Pool) DeletePool(opts *options.DeleteOptions) (*model.AtDeleteView, error) {
	if opts.Id == "" {
		return nil, errors.New("Id are required")
	}
	action := "DeletePool"
	para := &model.DeletePoolParams{
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

func (s *Pool) DescribePool(opts *options.BaseOptions) (*model.PoolView, error) {
	action := "DescribePool"
	bt, err := DescribeOpsResource(action, opts)
	if err != nil {
		return nil, err
	}
	view := &model.PoolView{}
	err = util.FormatView(bt, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (z *Pool) DescribePools(opts *options.ListOptions) (int64, []*model.PoolView, error) {
	action := "DescribePools"
	var ret []*model.PoolView
	desc, err := ListOpsResources(action, opts)
	if err != nil {
		return -1, nil, err
	}
	if err := util.FormatView(desc.Elements, &ret); err != nil {
		return -1, nil, err
	}
	return desc.Tc, ret, nil
}

func (s *Pool) UpdatePoolStatus(opts *options.UpdateStatusOptions) (*model.PoolUpdateView, error) {
	if opts.Id == "" {
		return nil, errors.New("Id is required")
	}
	action := "UpdatePoolStatus"
	para := &model.UpdatePoolStatusParams{
		Id:     &opts.Id,
		Status: &opts.Status,
	}
	dv, err := util.DoRequest(action, para)
	if err != nil {
		return nil, err
	}

	if dv.Data == nil {
		return nil, nil
	}

	view := &model.PoolUpdateView{}
	err = util.FormatView(dv.Data, view)
	if err != nil {
		return nil, err
	}
	return view, nil
}
